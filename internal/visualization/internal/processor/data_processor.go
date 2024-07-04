
package processor

import (
	"context"
	"fmt"
	"sort"

	pb "visualization/api/proto"
	client "visualization/data"
	"datasource/managers/query"
)

// DataProcessor handles the processing of data for visualizations.
type DataProcessor struct {
	dataSourceClient *client.DataSourceClient
}

// NewDataProcessor creates a new DataProcessor instance.
func NewDataProcessor(dataSourceAddress string) (*DataProcessor, error) {
	dsClient, err := client.NewDataSourceClient(dataSourceAddress)
	if err != nil {
		return nil, fmt.Errorf("failed to create DataSource client: %w", err)
	}
	return &DataProcessor{
		dataSourceClient: dsClient,
	}, nil
}

// ProcessData processes the data based on the visualization request.
func (dp *DataProcessor) ProcessData(ctx context.Context, req *pb.CreateVisualizationRequest) (map[string]interface{}, error) {
	// Create a query.Query object
	q := query.Query{
		Type:       query.Select,
		Collection: req.DataSourceId,
		Fields:     append(req.Dimensions, req.Measures...),
		// Conditions: req.Filters,
		// Limit:      int(req.Limit),
		// Offset:     int(req.Offset),
	}

	// Execute the query
	results, err := dp.dataSourceClient.ExecuteQuery(ctx, req.DataSourceId, q)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	// Process the results
	processedData, err := dp.processQueryResults(results, req.Dimensions, req.Measures)
	if err != nil {
		return nil, fmt.Errorf("failed to process query results: %w", err)
	}

	return processedData, nil
}

// processQueryResults processes the query results into a format suitable for visualization
func (dp *DataProcessor) processQueryResults(results []map[string]interface{}, dimensions, measures []string) (map[string]interface{}, error) {
	output := make(map[string]interface{})
	dimensionValues := make(map[string]map[string]bool)
	aggregatedData := make(map[string]map[string]float64)

	for _, result := range results {
		key := dp.getDimensionKey(result, dimensions)
		if _, exists := aggregatedData[key]; !exists {
			aggregatedData[key] = make(map[string]float64)
			for _, measure := range measures {
				aggregatedData[key][measure] = 0
			}
		}

		for _, dim := range dimensions {
			if dimensionValues[dim] == nil {
				dimensionValues[dim] = make(map[string]bool)
			}
			dimValue := fmt.Sprintf("%v", result[dim])
			dimensionValues[dim][dimValue] = true
		}

		for _, measure := range measures {
			if value, ok := result[measure].(float64); ok {
				aggregatedData[key][measure] += value
			} else if value, ok := result[measure].(int64); ok {
				aggregatedData[key][measure] += float64(value)
			}
		}
	}

	// Convert dimension values to sorted slices
	dimensionSlices := make(map[string][]string)
	for dim, valueMap := range dimensionValues {
		values := make([]string, 0, len(valueMap))
		for value := range valueMap {
			values = append(values, value)
		}
		sort.Strings(values)
		dimensionSlices[dim] = values
	}

	output["dimensions"] = dimensionSlices
	output["measures"] = measures
	output["data"] = aggregatedData

	return output, nil
}

// getDimensionKey creates a unique key based on dimension values
func (dp *DataProcessor) getDimensionKey(item map[string]interface{}, dimensions []string) string {
	var key string
	for _, dim := range dimensions {
		key += fmt.Sprintf("%v|", item[dim])
	}
	return key
}

// TimeSeriesProcessor handles time series data processing
