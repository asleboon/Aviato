package charting

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

// DrawChart renders a graph with viewer from a channel
func DrawChart(channelName string, channelViews []float64, viewTime []time.Time) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Time",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(), //enables / displays the x-axis
			ValueFormatter: chart.TimeHourValueFormatter,
		},
		YAxis: chart.YAxis{
			Name:      "Viewers",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(), //enables / displays the y-axis
		},
		Series: []chart.Series{
			chart.TimeSeries{
				Name: channelName,
				Style: chart.Style{
					Show:        true,
					StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
					FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
				},
				XValues: viewTime,
				YValues: channelViews,
			},
		},
	}

	// add legend
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	// Write to file.
	imgName := fmt.Sprintf(channelName + ".png")
	fo, err := os.Create(imgName)
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}

// DrawMulChart renders a graph with viewer from two channels
func DrawMulChart(channelOne string, viewsOne []float64, timesOne []time.Time, channelTwo string, viewsTwo []float64, timesTwo []time.Time) {
	chartName := fmt.Sprintf(channelOne + " viewers")
	graphOne := chart.TimeSeries{
		Name: chartName,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
		},
		XValues: timesOne,
		YValues: viewsOne,
	}

	chartName = fmt.Sprintf(channelTwo + " viewers")
	graphTwo := chart.TimeSeries{
		Name: chartName,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
		},
		XValues: timesTwo,
		YValues: viewsTwo,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:      "Time",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(), //enables / displays the x-axis
		},
		YAxis: chart.YAxis{
			Name:      "Viewers",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(), //enables / displays the y-axis
			// Range: &chart.ContinuousRange{
			// 	Max: 20000.0,
			// 	Min: 0.0,
			// },
		},
		Series: []chart.Series{
			graphOne,
			graphTwo,
		},
	}

	// add legend
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	// Write to file.
	fileName := fmt.Sprintf(channelOne + "AND" + channelTwo + ".png")
	fo, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}
