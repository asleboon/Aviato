// +build !solution

// Leave an empty line above this comment.
//
// Zap Collection Server
package main

import (
	"fmt"
	"log"
	"net"
	"sort"
	"time"

	"github.com/uis-dat320-fall18/Aviato/charting"
	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/zlog"
)

type UDPServer struct {
	conn *net.UDPConn
}

func (server *UDPServer) runLab() {
	switch *labnum {
	case "a", "c1", "c2", "d", "e":
		ztore = zlog.NewSimpleZapLogger()
	case "f":
		ztore = zlog.NewViewersZapLogger()
	case "g":
		ztoreGraph = zlog.NewChartLogger()
	}
	switch *labnum {
	case "a":
		go server.dumpAll()
	case "c1":
		go server.recordAll()
		go showViewers("NRK1")
	case "c2":
		go server.recordAll()
		go showViewers("TV2 Norge")
	case "d":
		// See answer in separate document.
	case "e":
		go server.recordAll()
		go top10Viewers()
	case "f":
		go server.recordAll()
		go top10Viewers()
	case "g":
		go server.recordAll()
		go drawChart("NRK1", "TV2 Norge")
	}
}

func NewUDPServer(addr string) (*UDPServer, error) {
	log.Println("Starting ZapServer...")
	// Build UDP address
	udpAddr, err := net.ResolveUDPAddr("udp", addr)

	// Create connection
	connUDP, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("NewUDPServer: Error creating UDP connection")
	}
	return &UDPServer{conn: connUDP}, nil
}

func (server *UDPServer) readFromUDP() (string, error) {
	buf := make([]byte, 256)                  // UDP packages usually ~50-70 bytes
	n, _, err := server.conn.ReadFromUDP(buf) // n = Number of bytes read
	str := string(buf[:n])
	return str, err
}

// dumpAll reads new STB events and prints to console
func (server *UDPServer) dumpAll() {
	for {
		eventStr, err := server.readFromUDP()
		if err != nil { // ReadFromUDP error check
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			fmt.Printf("Dumped response: %v\n", eventStr)
		}
	}
}

// recordAll processes and stores new viewers in Zaplogger
func (server *UDPServer) recordAll() {
	for {
		eventStr, err := server.readFromUDP()

		if err != nil {
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			chZap, _, err := chzap.NewSTBEvent(eventStr) // We don't care about statuschange here
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else {
				if chZap != nil {
					if ztore != nil {
						ztore.LogZap(*chZap) // (*chZap) = copy of pointer value
					}
					if ztoreGraph != nil {
						ztoreGraph.LogZap(*chZap) // Logger for logging data needed to create viewers graph
					}
				}
			}
		}
	}
}

// showViewers compute number of viewers on channel and prints every second
func showViewers(chName string) {
	tickChan := time.NewTicker(time.Second)
	defer tickChan.Stop()

	for range tickChan.C { // Runs code inside loop ~ every second
		views := ztore.Viewers(chName)
		fmt.Printf("No. of viewers on %s is now %d\n", chName, views)
	}
}

// top10Viewers prints top 10 channel views list every second
func top10Viewers() {
	tickChan := time.NewTicker(time.Second)
	defer tickChan.Stop()

	for range tickChan.C { // Runs code inside loop ~ every second
		channels := calculateTop10()

		fmt.Println("Top 10 channels with most viewers:")
		for i, c := range channels {
			fmt.Printf("%d. %v\n", i+1, c)
		}
		fmt.Println()
	}
}

// calculatTop10 computes top 10 views list
func calculateTop10() []*zlog.ChannelViewers {
	channels := ztore.ChannelsViewers()

	// Sort the channelviews, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Viewers > channels[j].Viewers
	})

	if len(channels) > 10 { // Only want top 10
		channels = channels[:10]
	}

	return channels
}

// calculatTop10 computes top 10 views list
func calculateTop10Muted() []*zlog.ChannelViewers {
	channels := ztore.ChannelsViewers()

	// Sort the channelviews, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Viewers > channels[j].Viewers
	})

	if len(channels) > 10 { // Only want top 10
		channels = channels[:10]
	}
	return channels
}

// drawChart creates viewers charts for two channels
func drawChart(channelOne string, channelTwo string) {
	tickChan := time.NewTicker(time.Hour * 24)
	defer tickChan.Stop()
	for range tickChan.C { // Runs code inside loop every 24hrs
		viewsOne, timesOne, viewsTwo, timesTwo := []float64{}, []time.Time{}, []float64{}, []time.Time{}
		log.Printf("Drawing chart for '%v' and '%v'...\n", channelOne, channelTwo)
		data := ztoreGraph.GetChartVal(channelOne)
		for _, value := range data {
			timesOne = append(timesOne, value.Times)
			viewsOne = append(viewsOne, value.Views)
		}
		data = ztoreGraph.GetChartVal(channelTwo)
		for _, value := range data {
			timesTwo = append(timesTwo, value.Times)
			viewsTwo = append(viewsTwo, value.Views)
		}
		// Draw one chart for each channel
		charting.DrawChart(channelOne, viewsOne, timesOne)
		charting.DrawChart(channelTwo, viewsTwo, timesTwo)

		// Draw a combined chart for both channels
		charting.DrawMulChart(channelOne, viewsOne, timesOne, channelTwo, viewsTwo, timesTwo)
	}

}
