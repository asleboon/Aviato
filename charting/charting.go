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
	imgName := fmt.Sprintf(channelName + ".png")
	fo, err := os.Create(imgName)
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(fo)
	fw.Write(buffer.Bytes())
}

/*
// DrawMulChart renders a graph with viewer from two channels
func DrawMulChart(channelOne string, viewsOne []float64, timesOne []time.Time, channeTwo string, viewsTwo []float64, timesTwo []time.Time) {
	chartName := fmt.Sprintf(channelOne + " viewers")
	graphOne := chart.TimeSeries{
		Name: chartName,
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(1).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(1).WithAlpha(64),
		},
		XValues: viewTime,
		YValues: channelViews1,
	}

	chartName = fmt.Sprintf(channelTwo + " viewers")
	graphTwo := chart.TimeSeries{
		Name: "Tv2 Viewers",
		Style: chart.Style{
			Show:        true,
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
		},
		XValues: viewTime,
		YValues: channelViews2,
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
			Range: &chart.ContinuousRange{
				Max: 20000.0,
				Min: 0.0,
			},
		},
		Series: []chart.Series{
			nrkGraph,
			tv2Graph,
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
*/
