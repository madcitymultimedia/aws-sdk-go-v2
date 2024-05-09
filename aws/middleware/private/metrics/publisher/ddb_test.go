package publisher

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
	"time"
)

type testDDBPublisherClient interface {
	PutItem(context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	AddTestData(*testing.T, map[string]types.AttributeValue)
}

type TestOutputClientWithError struct{}

func (TestOutputClientWithError) PutItem(ctx context.Context, input *dynamodb.PutItemInput, options ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	return nil, fmt.Errorf("PutItem error")
}

func (TestOutputClientWithError) AddTestData(*testing.T, map[string]types.AttributeValue) {}

type TestOutputClient struct {
	t    *testing.T
	item map[string]types.AttributeValue
}

func (c *TestOutputClient) PutItem(ctx context.Context, input *dynamodb.PutItemInput, options ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	if !reflect.DeepEqual(c.item, input.Item) {
		c.t.Errorf("Unexpected result, should be '%q' but was '%q'", c.item, input.Item)
	}
	return nil, nil
}

func (c *TestOutputClient) AddTestData(t *testing.T, item map[string]types.AttributeValue) {
	c.t = t
	c.item = item
}

func TestDDBOutputMetrics(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Data          []*metrics.MetricData
		DDBClient     testDDBPublisherClient
		ExpectedError error
		ExpectItem    map[string]types.AttributeValue
	}{
		"emptyRequestMetricData": {
			Data:          []*metrics.MetricData{{}},
			ExpectedError: nil,
			DDBClient:     &TestOutputClient{},
			ExpectItem: map[string]types.AttributeValue{
				"ServiceOperation": &types.AttributeValueMemberS{
					Value: ".",
				},
				"ApiCallDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"ApiCallSuccessful": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"MarshallingDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"EndpointResolutionDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"RetryCount": &types.AttributeValueMemberN{
					Value: "-1.00",
				},
				"InThroughput": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"OutThroughput": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"MaxConcurrency": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"AvailableConcurrency": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"ConcurrencyAcquireDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"PendingConcurrencyAcquires": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"SigningDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"UnmarshallingDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"TimeToFirstByte": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"ServiceCallDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
				"BackoffDelayDuration": &types.AttributeValueMemberN{
					Value: "0.00",
				},
			},
		},
		"putItemError": {
			Data:          []*metrics.MetricData{{}},
			DDBClient:     TestOutputClientWithError{},
			ExpectedError: fmt.Errorf("PutItem error"),
		},
		"completeRequestMetricData": {
			Data: []*metrics.MetricData{{
				RequestStartTime:         time.Unix(1234, 0),
				RequestEndTime:           time.Unix(1434, 0),
				SerializeStartTime:       time.Unix(1234, 0),
				SerializeEndTime:         time.Unix(1434, 0),
				ResolveEndpointStartTime: time.Unix(1234, 0),
				ResolveEndpointEndTime:   time.Unix(1434, 0),
				Success:                  1,
				ClientRequestID:          "crid",
				ServiceID:                "sid",
				OperationName:            "operationname",
				PartitionID:              "partitionid",
				Region:                   "region",
				RequestContentLength:     100,
				Stream:                   metrics.StreamMetrics{},
				Attempts: []metrics.AttemptMetrics{{
					ServiceCallStart:          time.Unix(1234, 0),
					ServiceCallEnd:            time.Unix(1434, 0),
					FirstByteTime:             time.Unix(1234, 0),
					ConnRequestedTime:         time.Unix(1234, 0),
					ConnObtainedTime:          time.Unix(1434, 0),
					CredentialFetchStartTime:  time.Unix(1234, 0),
					CredentialFetchEndTime:    time.Unix(1434, 0),
					SignStartTime:             time.Unix(1234, 0),
					SignEndTime:               time.Unix(1434, 0),
					DeserializeStartTime:      time.Unix(1234, 0),
					DeserializeEndTime:        time.Unix(1434, 0),
					RetryDelay:                200,
					ResponseContentLength:     100,
					StatusCode:                200,
					RequestID:                 "reqid",
					ExtendedRequestID:         "exreqid",
					HTTPClient:                "Default",
					MaxConcurrency:            5,
					PendingConnectionAcquires: 1,
					AvailableConcurrency:      2,
					ActiveRequests:            3,
					ReusedConnection:          false,
				},
					{
						ServiceCallStart:          time.Unix(1234, 0),
						ServiceCallEnd:            time.Unix(1434, 0),
						FirstByteTime:             time.Unix(1234, 0),
						ConnRequestedTime:         time.Unix(1234, 0),
						ConnObtainedTime:          time.Unix(1434, 0),
						CredentialFetchStartTime:  time.Unix(1234, 0),
						CredentialFetchEndTime:    time.Unix(1434, 0),
						SignStartTime:             time.Unix(1234, 0),
						SignEndTime:               time.Unix(1434, 0),
						DeserializeStartTime:      time.Unix(1234, 0),
						DeserializeEndTime:        time.Unix(1434, 0),
						RetryDelay:                100,
						ResponseContentLength:     100,
						StatusCode:                200,
						RequestID:                 "reqid",
						ExtendedRequestID:         "exreqid",
						HTTPClient:                "Default",
						MaxConcurrency:            5,
						PendingConnectionAcquires: 1,
						AvailableConcurrency:      2,
						ActiveRequests:            3,
						ReusedConnection:          false,
					}},
			},
				{
					RequestStartTime:         time.Unix(1234, 0),
					RequestEndTime:           time.Unix(1634, 0),
					SerializeStartTime:       time.Unix(1234, 0),
					SerializeEndTime:         time.Unix(1634, 0),
					ResolveEndpointStartTime: time.Unix(1234, 0),
					ResolveEndpointEndTime:   time.Unix(1634, 0),
					Success:                  1,
					ClientRequestID:          "crid",
					ServiceID:                "sid",
					OperationName:            "operationname",
					PartitionID:              "partitionid",
					Region:                   "region",
					RequestContentLength:     100,
					Stream:                   metrics.StreamMetrics{},
					Attempts: []metrics.AttemptMetrics{{
						ServiceCallStart:          time.Unix(1234, 0),
						ServiceCallEnd:            time.Unix(1634, 0),
						FirstByteTime:             time.Unix(1234, 0),
						ConnRequestedTime:         time.Unix(1234, 0),
						ConnObtainedTime:          time.Unix(1634, 0),
						CredentialFetchStartTime:  time.Unix(1234, 0),
						CredentialFetchEndTime:    time.Unix(1634, 0),
						SignStartTime:             time.Unix(1234, 0),
						SignEndTime:               time.Unix(1634, 0),
						DeserializeStartTime:      time.Unix(1234, 0),
						DeserializeEndTime:        time.Unix(1634, 0),
						RetryDelay:                100,
						ResponseContentLength:     100,
						StatusCode:                200,
						RequestID:                 "reqid",
						ExtendedRequestID:         "exreqid",
						HTTPClient:                "Default",
						MaxConcurrency:            10,
						PendingConnectionAcquires: 2,
						AvailableConcurrency:      4,
						ActiveRequests:            3,
						ReusedConnection:          false,
					},
						{
							ServiceCallStart:          time.Unix(1234, 0),
							ServiceCallEnd:            time.Unix(1634, 0),
							FirstByteTime:             time.Unix(1634, 0),
							ConnRequestedTime:         time.Unix(1234, 0),
							ConnObtainedTime:          time.Unix(1634, 0),
							CredentialFetchStartTime:  time.Unix(1234, 0),
							CredentialFetchEndTime:    time.Unix(1634, 0),
							SignStartTime:             time.Unix(1234, 0),
							SignEndTime:               time.Unix(1634, 0),
							DeserializeStartTime:      time.Unix(1234, 0),
							DeserializeEndTime:        time.Unix(1634, 0),
							RetryDelay:                100,
							ResponseContentLength:     100,
							StatusCode:                200,
							RequestID:                 "reqid",
							ExtendedRequestID:         "exreqid",
							HTTPClient:                "Default",
							MaxConcurrency:            10,
							PendingConnectionAcquires: 2,
							AvailableConcurrency:      4,
							ActiveRequests:            3,
							ReusedConnection:          false,
						}},
				},
			},
			DDBClient: &TestOutputClient{},
			ExpectItem: map[string]types.AttributeValue{
				"ServiceOperation": &types.AttributeValueMemberS{
					Value: "sid.operationname",
				},
				"ApiCallDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"ApiCallSuccessful": &types.AttributeValueMemberN{
					Value: "1.00",
				},
				"MarshallingDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"EndpointResolutionDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"RetryCount": &types.AttributeValueMemberN{
					Value: "1.00",
				},
				"InThroughput": &types.AttributeValueMemberN{
					Value: "0.38",
				},
				"OutThroughput": &types.AttributeValueMemberN{
					Value: "0.38",
				},
				"MaxConcurrency": &types.AttributeValueMemberN{
					Value: "7.50",
				},
				"AvailableConcurrency": &types.AttributeValueMemberN{
					Value: "3.00",
				},
				"ConcurrencyAcquireDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"PendingConcurrencyAcquires": &types.AttributeValueMemberN{
					Value: "1.50",
				},
				"SigningDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"UnmarshallingDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"TimeToFirstByte": &types.AttributeValueMemberN{
					Value: "100000000000.00",
				},
				"ServiceCallDuration": &types.AttributeValueMemberN{
					Value: "300000000000.00",
				},
				"BackoffDelayDuration": &types.AttributeValueMemberN{
					Value: "125.00",
				},
			},
			ExpectedError: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			publisher := NewDDBPublisher()

			for _, data := range c.Data {
				data.ComputeRequestMetrics()
				publisher.PostRequestMetrics(data)
			}

			c.DDBClient.AddTestData(t, c.ExpectItem)

			err := publisher.Output(c.DDBClient, "")

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}
		})
	}
}
