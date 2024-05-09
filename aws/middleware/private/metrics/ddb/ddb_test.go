package ddb

import (
	"github.com/aws/aws-sdk-go-v2/internal/sdk"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"reflect"
	"testing"
	"time"
)

func TestCreateNewEntry(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Namespace     string
		ExpectedEntry Entry
	}{
		"success": {
			ExpectedEntry: Entry{
				metrics:    []metric{},
				dimensions: [][]string{{}},
				fields:     map[string]interface{}{},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			actualEntry := NewEntry()
			if !reflect.DeepEqual(actualEntry, c.ExpectedEntry) {
				t.Errorf("Entry contained unexpected values")
			}
		})
	}
}

func TestBuild(t *testing.T) {

	sdk.NowTime = func() time.Time {
		return time.Unix(1234, 0)
	}

	cases := map[string]struct {
		Configure      func(entry *Entry)
		ExpectedError  error
		ExpectedResult map[string]types.AttributeValue
	}{
		"normalEntry": {
			Configure: func(entry *Entry) {
				entry.AddMetric("testMetric1", 1)
				entry.AddDimension("testDimension1", "dim1")
			},
			ExpectedError: nil,
			ExpectedResult: map[string]types.AttributeValue{
				"testMetric1": &types.AttributeValueMemberN{
					Value: "1",
				},
				"testDimension1": &types.AttributeValueMemberS{
					Value: "dim1",
				},
			},
		},
	}

	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			entry := NewEntry()

			c.Configure(&entry)

			result, err := entry.Build()

			if !reflect.DeepEqual(err, c.ExpectedError) {
				t.Errorf("Unexpected error, should be '%s' but was '%s'", c.ExpectedError, err)
			}

			if !reflect.DeepEqual(result, c.ExpectedResult) {
				t.Errorf("Unexpected result, should be '%s' but was '%s'", c.ExpectedResult, result)
			}
		})
	}
}
