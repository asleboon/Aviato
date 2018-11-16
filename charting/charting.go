package charting

import (
	"bufio"
	"bytes"
	"os"
	"time"

	"github.com/wcharczuk/go-chart"
)

func DrawChart(channelName string, channelViews []float64, viewTime []time.Time) {
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
	fo, err := os.Create("image.png")
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}
