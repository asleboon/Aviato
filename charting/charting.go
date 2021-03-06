package charting

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/wcharczuk/go-chart"
)

// DrawChart renders a graph with viewer from a channel
func DrawChart(channelName string, channelViews []float64, viewTime []time.Time) {
	graph := chart.Chart{
		XAxis: chart.XAxis{
			Name:           "Time",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),              // Displays the x-axis
			ValueFormatter: chart.TimeMinuteValueFormatter, // Add desired time format
		},
		YAxis: chart.YAxis{
			Name:      "Viewers",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(), // Displays the x-axis
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

	// Add legend
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	filePath := fmt.Sprintf("../charting/" + channelName + "Viewers" + ".png")
	filePath = strings.Replace(filePath, " ", "", -1) // Remove possible whitespace from channelnames
	fo, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}

// DrawMulChart renders a graph with viewer from two channels
func DrawMulChart(channelOne string, viewsOne []float64, timesOne []time.Time, channelTwo string, viewsTwo []float64, timesTwo []time.Time) {
	graphOne := chart.TimeSeries{
		Name: channelOne,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
		},
		XValues: timesOne,
		YValues: viewsOne,
	}

	graphTwo := chart.TimeSeries{
		Name: channelTwo,
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
			Name:           "Time",
			NameStyle:      chart.StyleShow(),
			Style:          chart.StyleShow(),              // Displays the x-axis
			ValueFormatter: chart.TimeMinuteValueFormatter, // Add desired time format
		},
		YAxis: chart.YAxis{
			Name:      "Viewers",
			NameStyle: chart.StyleShow(),
			Style:     chart.StyleShow(), // Displays the x-axis
		},
		Series: []chart.Series{
			graphOne,
			graphTwo,
		},
	}

	// Add legend
	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	buffer := bytes.NewBuffer([]byte{})
	graph.Render(chart.PNG, buffer)

	// Write file to charting folder
	filePath := fmt.Sprintf("../charting/" + channelOne + "And" + channelTwo + "Viewers" + ".png")
	filePath = strings.Replace(filePath, " ", "", -1) // Remove possible whitespace from channelnames
	fo, err := os.Create(filePath)
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}
