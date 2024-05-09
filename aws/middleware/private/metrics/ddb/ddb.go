package ddb

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Entry struct {
	metrics    []metric
	dimensions [][]string
	fields     map[string]interface{}
}

type metric struct {
	Name string
}

// NewEntry creates a new Entry with the specified namespace
func NewEntry() Entry {
	return Entry{
		metrics:    []metric{},
		dimensions: [][]string{{}},
		fields:     map[string]interface{}{},
	}
}

// Build constructs the DDBItem from Entry
func (e *Entry) Build() (map[string]types.AttributeValue, error) {
	item := make(map[string]types.AttributeValue, 0)
	for k, v := range e.fields {
		switch t := v.(type) {
		case float64:
			item[k] = &types.AttributeValueMemberN{
				Value: fmt.Sprintf("%.2f", t),
			}
		case string:
			item[k] = &types.AttributeValueMemberS{
				Value: t,
			}
		default:
		}
	}

	return item, nil
}

func (e *Entry) AddDimension(key string, value string) {
	e.dimensions[0] = append(e.dimensions[0], key)
	e.fields[key] = value
}

func (e *Entry) AddMetric(key string, value float64) {
	e.metrics = append(e.metrics, metric{key})
	e.fields[key] = value
}
