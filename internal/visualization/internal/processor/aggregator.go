package processor

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"
	pb "visualization/api/proto"

	"datasource/managers/query"
)

// Aggregator handles data aggregation.
type Aggregator struct{}

// NewAggregator creates a new Aggregator instance.
func NewAggregator() *Aggregator {
	return &Aggregator{}
}

// Aggregate performs data aggregation based on dimensions and measures.
func (a *Aggregator) Aggregate(data []map[string]interface{}, dimensions, measures []string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	dimensionValues := make(map[string][]string)
	aggregatedData := make(map[string]map[string]float64)

	// Collect unique dimension values and initialize aggregated data
	for _, item := range data {
		key := a.getDimensionKey(item, dimensions)
		if _, exists := aggregatedData[key]; !exists {
			aggregatedData[key] = make(map[string]float64)
			for _, measure := range measures {
				aggregatedData[key][measure] = 0
			}
		}

		for _, dim := range dimensions {
			value := fmt.Sprintf("%v", item[dim])
			if !contains(dimensionValues[dim], value) {
				dimensionValues[dim] = append(dimensionValues[dim], value)
			}
		}

		for _, measure := range measures {
			if value, ok := item[measure].(float64); ok {
				aggregatedData[key][measure] += value
			}
		}
	}

	// Sort dimension values
	for dim, values := range dimensionValues {
		sort.Strings(values)
		dimensionValues[dim] = values
	}

	result["dimensions"] = dimensionValues
	result["measures"] = measures
	result["data"] = aggregatedData

	return result, nil
}

// getDimensionKey creates a unique key based on dimension values.
func (a *Aggregator) getDimensionKey(item map[string]interface{}, dimensions []string) string {
	var keyParts []string
	for _, dim := range dimensions {
		keyParts = append(keyParts, fmt.Sprintf("%v", item[dim]))
	}
	return strings.Join(keyParts, "|")
}

// contains checks if a string slice contains a specific value.
func contains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

type TimeSeriesProcessor struct {
	DataProcessor
}

// ProcessTimeSeriesData processes time series data
func (tp *TimeSeriesProcessor) ProcessTimeSeriesData(ctx context.Context, req *pb.CreateVisualizationRequest, timeDimension string, interval string) (map[string]interface{}, error) {
	results, err := tp.DataProcessor.dataSourceClient.ExecuteQuery(ctx, req.DataSourceId, query.Query{
		Type:       query.Select,
		Collection: req.DataSourceId,
		Fields:     append(req.Dimensions, req.Measures...),
		// Conditions: req.Filters,
		// Limit:      int(req.Limit),
		// Offset:     int(req.Offset),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	return tp.processTimeSeriesResults(results, timeDimension, req.Measures, interval)
}

func (tp *TimeSeriesProcessor) processTimeSeriesResults(results []map[string]interface{}, timeDimension string, measures []string, interval string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	timeSeriesData := make(map[string]map[string]float64)

	for _, result := range results {
		timeValue, ok := result[timeDimension].(string)
		if !ok {
			return nil, fmt.Errorf("invalid time value for dimension %s", timeDimension)
		}

		t, err := time.Parse(time.RFC3339, timeValue)
		if err != nil {
			return nil, fmt.Errorf("failed to parse time: %w", err)
		}

		bucketKey := tp.getTimeBucketKey(t, interval)
		if _, exists := timeSeriesData[bucketKey]; !exists {
			timeSeriesData[bucketKey] = make(map[string]float64)
			for _, measure := range measures {
				timeSeriesData[bucketKey][measure] = 0
			}
		}

		for _, measure := range measures {
			if value, ok := result[measure].(float64); ok {
				timeSeriesData[bucketKey][measure] += value
			} else if value, ok := result[measure].(int64); ok {
				timeSeriesData[bucketKey][measure] += float64(value)
			}
		}
	}

	output["timeDimension"] = timeDimension
	output["interval"] = interval
	output["measures"] = measures
	output["data"] = timeSeriesData

	return output, nil
}

func (tp *TimeSeriesProcessor) getTimeBucketKey(t time.Time, interval string) string {
	switch interval {
	case "hour":
		return t.Format("2006-01-02 15:00")
	case "day":
		return t.Format("2006-01-02")
	case "week":
		year, week := t.ISOWeek()
		return fmt.Sprintf("%d-W%02d", year, week)
	case "month":
		return t.Format("2006-01")
	case "year":
		return t.Format("2006")
	default:
		return t.Format(time.RFC3339)
	}
}

// Close closes the DataSourceClient connection
func (dp *DataProcessor) Close() error {
	return dp.dataSourceClient.Close()
}