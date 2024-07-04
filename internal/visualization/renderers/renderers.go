package renderers

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"sort"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

// Renderer is the interface that all chart renderers must implement
type Renderer interface {
	Render(data map[string]interface{}) (string, error)
}

// BaseRenderer provides common functionality for renderers
type BaseRenderer struct{}

func (br *BaseRenderer) renderToBase64(renderFunc func(w io.Writer) error) (string, error) {
	buffer := new(bytes.Buffer)
	if err := renderFunc(buffer); err != nil {
		return "", fmt.Errorf("failed to render chart: %w", err)
	}
	return base64.StdEncoding.EncodeToString(buffer.Bytes()), nil
}

// BarChartRenderer renders bar charts
type BarChartRenderer struct {
	BaseRenderer
}

func (bcr *BarChartRenderer) Render(data map[string]interface{}) (string, error) {
	dimensions, ok := data["dimensions"].(map[string][]string)
	if !ok {
		return "", fmt.Errorf("invalid dimensions data")
	}
	measures, ok := data["measures"].([]string)
	if !ok {
		return "", fmt.Errorf("invalid measures data")
	}
	chartData, ok := data["data"].(map[string]map[string]float64)
	if !ok {
		return "", fmt.Errorf("invalid chart data")
	}

	// Assume we're using the first dimension for x-axis
	xAxisData := dimensions[measures[0]]

	bar := charts.NewBar()
	bar.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Bar Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
	)

	bar.SetXAxis(xAxisData)

	for _, measure := range measures {
		series := make([]opts.BarData, 0, len(xAxisData))
		for _, x := range xAxisData {
			series = append(series, opts.BarData{Value: chartData[x][measure]})
		}
		bar.AddSeries(measure, series)
	}

	return bcr.renderToBase64(bar.Render)
}

// LineChartRenderer renders line charts
type LineChartRenderer struct {
	BaseRenderer
}

func (lcr *LineChartRenderer) Render(data map[string]interface{}) (string, error) {
	dimensions, ok := data["dimensions"].(map[string][]string)
	if !ok {
		return "", fmt.Errorf("invalid dimensions data")
	}
	measures, ok := data["measures"].([]string)
	if !ok {
		return "", fmt.Errorf("invalid measures data")
	}
	chartData, ok := data["data"].(map[string]map[string]float64)
	if !ok {
		return "", fmt.Errorf("invalid chart data")
	}

	// Assume we're using the first dimension for x-axis
	xAxisData := dimensions[measures[0]]

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Line Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
	)

	line.SetXAxis(xAxisData)

	for _, measure := range measures {
		series := make([]opts.LineData, 0, len(xAxisData))
		for _, x := range xAxisData {
			series = append(series, opts.LineData{Value: chartData[x][measure]})
		}
		line.AddSeries(measure, series)
	}

	return lcr.renderToBase64(line.Render)
}

// PieChartRenderer renders pie charts
type PieChartRenderer struct {
	BaseRenderer
}

func (pcr *PieChartRenderer) Render(data map[string]interface{}) (string, error) {
	dimensions, ok := data["dimensions"].(map[string][]string)
	if !ok {
		return "", fmt.Errorf("invalid dimensions data")
	}
	measures, ok := data["measures"].([]string)
	if !ok || len(measures) == 0 {
		return "", fmt.Errorf("invalid measures data")
	}
	chartData, ok := data["data"].(map[string]map[string]float64)
	if !ok {
		return "", fmt.Errorf("invalid chart data")
	}

	// Use the first dimension and measure for the pie chart
	dimension := measures[0]
	measure := measures[0]

	pie := charts.NewPie()
	pie.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Pie Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
	)

	pieData := make([]opts.PieData, 0, len(dimensions[dimension]))
	for _, label := range dimensions[dimension] {
		pieData = append(pieData, opts.PieData{Name: label, Value: chartData[label][measure]})
	}
	pie.AddSeries("Category", pieData)

	return pcr.renderToBase64(pie.Render)
}

// ScatterPlotRenderer renders scatter plots
type ScatterPlotRenderer struct {
	BaseRenderer
}

func (spr *ScatterPlotRenderer) Render(data map[string]interface{}) (string, error) {
	dimensions, ok := data["dimensions"].(map[string][]string)
	if !ok {
		return "", fmt.Errorf("invalid dimensions data")
	}
	if len(dimensions) < 2 {
		return "", fmt.Errorf("need at least two dimensions for scatter plot")
	}
	measures, ok := data["measures"].([]string)
	if !ok || len(measures) < 2 {
		return "", fmt.Errorf("invalid measures data, need at least two measures")
	}
	chartData, ok := data["data"].(map[string]map[string]float64)
	if !ok {
		return "", fmt.Errorf("invalid chart data")
	}

	// Use the first two measures for x and y axes
	xMeasure, yMeasure := measures[0], measures[1]

	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Scatter Plot"}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
		charts.WithXAxisOpts(opts.XAxis{Name: xMeasure}),
		charts.WithYAxisOpts(opts.YAxis{Name: yMeasure}),
	)

	scatterData := make([]opts.ScatterData, 0, len(chartData))
	for _, dataPoint := range chartData {
		scatterData = append(scatterData, opts.ScatterData{
			Value:        []interface{}{dataPoint[xMeasure], dataPoint[yMeasure]},
			Symbol:       "circle",
			SymbolSize:   10,
			// ItemStyle:    &opts.ItemStyle{Color: "blue"},
			SymbolRotate: 0,
		})
	}
	scatter.AddSeries("Data Points", scatterData)

	return spr.renderToBase64(scatter.Render)
}

// TimeSeriesRenderer renders time series charts
type TimeSeriesRenderer struct {
	BaseRenderer
}

func (tsr *TimeSeriesRenderer) Render(data map[string]interface{}) (string, error) {
	timeDimension, ok := data["timeDimension"].(string)
	if !ok {
		return "", fmt.Errorf("invalid time dimension")
	}
	measures, ok := data["measures"].([]string)
	if !ok {
		return "", fmt.Errorf("invalid measures data")
	}
	chartData, ok := data["data"].(map[string]map[string]float64)
	if !ok {
		return "", fmt.Errorf("invalid chart data")
	}

	// Sort time keys
	timeKeys := make([]string, 0, len(chartData))
	for k := range chartData {
		timeKeys = append(timeKeys, k)
	}
	sort.Strings(timeKeys)

	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Time Series Chart"}),
		charts.WithTooltipOpts(opts.Tooltip{}),
		charts.WithLegendOpts(opts.Legend{}),
		charts.WithXAxisOpts(opts.XAxis{Name: timeDimension}),
	)

	line.SetXAxis(timeKeys)

	for _, measure := range measures {
		series := make([]opts.LineData, 0, len(timeKeys))
		for _, time := range timeKeys {
			series = append(series, opts.LineData{Value: chartData[time][measure]})
		}
		line.AddSeries(measure, series)
	}

	return tsr.renderToBase64(line.Render)
}

// RendererFactory creates renderers based on the visualization type
func RendererFactory(visualizationType string) (Renderer, error) {
	switch visualizationType {
	case "bar":
		return &BarChartRenderer{}, nil
	case "line":
		return &LineChartRenderer{}, nil
	case "pie":
		return &PieChartRenderer{}, nil
	case "scatter":
		return &ScatterPlotRenderer{}, nil
	case "timeseries":
		return &TimeSeriesRenderer{}, nil
	default:
		return nil, fmt.Errorf("unsupported visualization type: %s", visualizationType)
	}
}