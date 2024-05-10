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
		RequestData   []*metrics.MetricData
		StreamData    []*metrics.MetricData
		DDBClient     testDDBPublisherClient
		ExpectedError error
		ExpectItem    map[string]types.AttributeValue
	}{
		"emptyRequestMetricsData": {
			RequestData:   []*metrics.MetricData{{}},
			StreamData:    []*metrics.MetricData{},
			ExpectedError: nil,
			DDBClient:     &TestOutputClient{},
			ExpectItem:    emptyRequestItem,
		},
		"putItemError": {
			RequestData:   []*metrics.MetricData{{}},
			DDBClient:     TestOutputClientWithError{},
			ExpectedError: fmt.Errorf("PutItem error"),
		},
		"completeRequestMetricData": {
			RequestData: completeRequestMetricDataDDBPublisher,
			StreamData: []*metrics.MetricData{
				{
					Stream: metrics.StreamMetrics{
						Throughput: 80000000,
					},
				},
				{
					Stream: metrics.StreamMetrics{
						Throughput: 60000000,
					},
				},
			},
			DDBClient:     &TestOutputClient{},
			ExpectItem:    completeRequestItem,
			ExpectedError: nil,
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			publisher := NewDDBPublisher()

			for _, data := range c.RequestData {
				data.ComputeRequestMetrics()
				publisher.PostRequestMetrics(data)
			}
			for _, data := range c.StreamData {
				publisher.PostStreamMetrics(data)
			}

			c.DDBClient.AddTestData(t, c.ExpectItem)

			err := publisher.Output(c.DDBClient, "")

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}
		})
	}
}
