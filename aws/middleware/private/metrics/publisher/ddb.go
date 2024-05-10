package publisher

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics/ddb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DDBPublisher is a MetricPublisher implementation that publishes metrics to DDB
type DDBPublisher struct {
	additionalDimensions map[string]string
	requestDataEntries   []*metrics.MetricData
	streamDataEntries    []*metrics.MetricData
	totalAttempts        int // recorded in advance to calculate attempt metrics average
}

func NewDDBPublisher() *DDBPublisher {
	return &DDBPublisher{
		additionalDimensions: map[string]string{},
		requestDataEntries:   []*metrics.MetricData{},
		streamDataEntries:    []*metrics.MetricData{},
	}
}

// SetAdditionalDimension adds a dimension's to additionalDimensions
func (p *DDBPublisher) SetAdditionalDimension(key string, value string) {
	p.additionalDimensions[key] = value
}

// RemoveAdditionalDimension removes a dimension from additionalDimensions
func (p *DDBPublisher) RemoveAdditionalDimension(key string) {
	delete(p.additionalDimensions, key)
}

// PostRequestMetrics stores the request metrics in DDBPublisher
func (p *DDBPublisher) PostRequestMetrics(data *metrics.MetricData) error {
	p.requestDataEntries = append(p.requestDataEntries, data)
	p.totalAttempts = p.totalAttempts + len(data.Attempts)
	return nil
}

// PostStreamMetrics caches stream metrics to DDBPublisher
func (p *DDBPublisher) PostStreamMetrics(data *metrics.MetricData) error {
	if data.Stream.Throughput > 0 {
		p.streamDataEntries = append(p.streamDataEntries, data)
	}
	return nil
}

type ddbClient interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

// Output calculate average metrics and send that to ddb
func (p *DDBPublisher) Output(svc ddbClient, tableName string) error {
	if len(p.requestDataEntries) == 0 {
		return fmt.Errorf("no metric data is cached in publisher")
	}

	entry := ddb.NewEntry()
	p.populateWithAdditionalDimensions(&entry)

	// this 2 dimensions might not need since they can be identified through the input tableName
	// e.g. tableName = S3.PutObject
	entry.AddDimension(metrics.ServiceIDKey, p.requestDataEntries[0].ServiceID)
	entry.AddDimension(metrics.OperationNameKey, p.requestDataEntries[0].OperationName)

	p.averageRequestMetrics(&entry)

	p.averageStreamMetrics(&entry)

	item, err := entry.Build()
	if err != nil {
		return fmt.Errorf("error generating ddb item for metric: %s", err.Error())
	}

	_, err = svc.PutItem(context.Background(), &dynamodb.PutItemInput{
		TableName: aws.String(tableName),
		Item:      item,
	})
	if err != nil {
		return err
	}

	// cleans our stored metricDataEntries for next round test
	p.requestDataEntries = []*metrics.MetricData{}
	p.streamDataEntries = []*metrics.MetricData{}
	return nil
}

func (p DDBPublisher) averageRequestMetrics(entry *ddb.Entry) {
	// average metrics per request
	var apiCallDuration float64
	var apiCallSuccessful float64
	var marshallingDuration float64
	var endpointResolutionDuration float64

	var retryCount float64

	var inThroughput float64
	var outThroughput float64

	// average metrics per attempt
	var maxConcurrency float64
	var availableConcurrency float64
	var concurrencyAcquireDuration float64
	var pendingConcurrencyAcquires float64
	var signingDuration float64
	var unmarshallingDuration float64
	var timeToFirstByte float64
	var serviceCallDuration float64
	var backoffDelayDuration float64

	size := float64(len(p.requestDataEntries))
	for _, data := range p.requestDataEntries {
		apiCallDuration += float64(data.APICallDuration.Nanoseconds()) / size
		marshallingDuration += float64(data.MarshallingDuration.Nanoseconds()) / size
		endpointResolutionDuration += float64(data.EndpointResolutionDuration.Nanoseconds()) / size
		retryCount += float64(data.RetryCount)

		if data.InThroughput != 0 {
			inThroughput += data.InThroughput / size
		}
		if data.OutThroughput != 0 {
			outThroughput += data.OutThroughput / size
		}

		if data.Success > 0 {
			apiCallSuccessful++
		}

		attempts := float64(p.totalAttempts)
		for _, attempt := range data.Attempts {
			maxConcurrency += float64(attempt.MaxConcurrency) / attempts
			availableConcurrency += float64(attempt.AvailableConcurrency) / attempts
			concurrencyAcquireDuration += float64(attempt.ConcurrencyAcquireDuration.Nanoseconds()) / attempts
			pendingConcurrencyAcquires += float64(attempt.PendingConnectionAcquires) / attempts
			signingDuration += float64(attempt.SigningDuration.Nanoseconds()) / attempts
			unmarshallingDuration += float64(attempt.UnMarshallingDuration.Nanoseconds()) / attempts
			timeToFirstByte += float64(attempt.TimeToFirstByte.Nanoseconds()) / attempts
			serviceCallDuration += float64(attempt.ServiceCallDuration.Nanoseconds()) / attempts
			backoffDelayDuration += float64(attempt.RetryDelay.Nanoseconds()) / attempts
		}
	}

	retryCount /= size
	apiCallSuccessful /= size

	entry.AddMetric(metrics.APICallDurationKey, apiCallDuration)
	entry.AddMetric(metrics.APICallSuccessfulKey, apiCallSuccessful)
	entry.AddMetric(metrics.MarshallingDurationKey, marshallingDuration)
	entry.AddMetric(metrics.EndpointResolutionDurationKey, endpointResolutionDuration)

	entry.AddMetric(metrics.RetryCountKey, retryCount)
	entry.AddMetric(metrics.InThroughputKey, inThroughput)
	entry.AddMetric(metrics.OutThroughputKey, outThroughput)

	entry.AddMetric(metrics.MaxConcurrencyKey, maxConcurrency)
	entry.AddMetric(metrics.AvailableConcurrencyKey, availableConcurrency)
	entry.AddMetric(metrics.ConcurrencyAcquireDurationKey, concurrencyAcquireDuration)
	entry.AddMetric(metrics.PendingConcurrencyAcquiresKey, pendingConcurrencyAcquires)
	entry.AddMetric(metrics.SigningDurationKey, signingDuration)
	entry.AddMetric(metrics.UnmarshallingDurationKey, unmarshallingDuration)
	entry.AddMetric(metrics.TimeToFirstByteKey, timeToFirstByte)
	entry.AddMetric(metrics.ServiceCallDurationKey, serviceCallDuration)
	entry.AddMetric(metrics.BackoffDelayDurationKey, backoffDelayDuration)
}

func (p DDBPublisher) averageStreamMetrics(entry *ddb.Entry) {
	if size := float64(len(p.streamDataEntries)); size > 0 {
		var throughput float64
		for _, data := range p.streamDataEntries {
			throughput += data.Stream.Throughput / size
		}
		entry.AddMetric(metrics.StreamThroughputKey, throughput)
	}
}

func (p *DDBPublisher) populateWithAdditionalDimensions(entry *ddb.Entry) {
	for k := range p.additionalDimensions {
		entry.AddDimension(k, p.additionalDimensions[k])
	}
}
