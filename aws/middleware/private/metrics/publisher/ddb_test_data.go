package publisher

import (
	"github.com/aws/aws-sdk-go-v2/aws/middleware/private/metrics"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"time"
)

var emptyRequestItem = map[string]types.AttributeValue{
	"ServiceId": &types.AttributeValueMemberS{
		Value: "",
	},
	"OperationName": &types.AttributeValueMemberS{
		Value: "",
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
}

var completeRequestMetricDataDDBPublisher = []*metrics.MetricData{{
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
}

var completeRequestItem = map[string]types.AttributeValue{
	"ServiceId": &types.AttributeValueMemberS{
		Value: "sid",
	},
	"OperationName": &types.AttributeValueMemberS{
		Value: "operationname",
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
	"Throughput": &types.AttributeValueMemberN{
		Value: "70000000.00",
	},
}
