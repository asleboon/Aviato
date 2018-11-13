package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"github.com/uis-dat320-fall18/Aviato/zlog"
	"google.golang.org/grpc"
)

type SubscribeServer struct {
	logger zlog.AdvZapLogger
}

var conn *net.UDPConn
var err error

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	endpoint = flag.String(
		"endpoint",
		"localhost:1994", // Changed port from std to 1994 to avoid problems during testing.
		"Endpoint on which server runs. Preferable",
	)
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func parseFlags() {
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}
}

// TODO: Split gRPC server and zapserver part into separate files
func startZapServer() {
	log.Println("Starting ZapServer...")
	// Build UDP address
	addr, _ := net.ResolveUDPAddr("udp", "224.0.1.130:10000")

	// Create connection
	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("NewUDPServer: Error creating UDP connection")
	}
}

func readFromUDP() (string, error) {
	buf := make([]byte, 256)           // UDP packages usually ~50-70 bytes
	n, _, err := conn.ReadFromUDP(buf) // n = Number of bytes read
	str := string(buf[:n])
	return str, err
}

// recordAll processes and stores new viewers in Zaplogger
func (s *SubscribeServer) recordAll() {
	for {
		eventStr, err := readFromUDP()

		if err != nil { // ReadFromUDP error check
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			chZap, stChange, err := chzap.NewSTBEvent(eventStr)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			} else if chZap != nil {
				s.logger.LogZap(*chZap) // Pass a copy of pointer value
			} else if stChange != nil {
				s.logger.LogStatus(*stChange) // Pass a copy of pointer value
			}
		}
	}
}

func (s *SubscribeServer) top10Viewers() string {
	channels := s.logger.ChannelsViewers() // Map of all channels with number of viewers

	// Sort channels by views, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Viewers > channels[j].Viewers
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}
		top10Str += fmt.Sprintf("%v. %v, viewers: %v", count+1, v.Channel, v.Viewers)
	}
	top10Str += "\n\n" // Easy way to create space between top 10 prints
	return top10Str
}

func (s *SubscribeServer) top10Duration() string {
	channels := s.logger.ChannelsDuration() // Map of all channels with total duration

	// Sort channels by total duration, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Duration > channels[j].Duration
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}
		top10Str += fmt.Sprintf("%v. %v, total duration: %v", count+1, v.Channel, v.Duration)
	}
	top10Str += "\n\n"
	return top10Str
}

func (s *SubscribeServer) top10Mute() string {
	channels := s.logger.ChannelsMute() // Map of all channels with avg. muted duration per viewer
	// Sort channels by avg mute per viewer, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].AvgMute > channels[j].AvgMute
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}
		top10Str += fmt.Sprintf("%v. %v, average muted duration per viewer: %v\n", count+1, v.Channel, v.AvgMute)
		top10Str += fmt.Sprintf("Time with highest number of muted viewers: %v", v.MaxMuteTime)
	}
	top10Str += "\n\n"
	return top10Str
}

func (s *SubscribeServer) sma(smaChannel string, smaLength uint64) string {
	views := s.logger.Viewers(smaChannel)
	timeNow := time.Now()
	s.logger.sma[timeNow] = views

	sumViewers := 0
	count := 0
	for k, v := range s.logger.sma {
		if timeNow.Sub(k) < (time.Duration(smaLength) * time.Second) {
			sumViewers += v
			count++
		}
	}
	if count == 0 {
		return fmt.Sprintf("Simple moving average for %s: %d\n", smaChannel, 0)
	}
	return fmt.Sprintf("Simple moving average for %s: %d\n", smaChannel, sumViewers/count)
}

// Subscribe handles a client subscription request
func (s *SubscribeServer) Subscribe(stream pb.Subscription_SubscribeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF { // Do we need this?
			return nil
		} else if err != nil {
			return err
		}

		tickChan := time.NewTicker(time.Second * time.Duration(in.RefreshRate))
		defer tickChan.Stop()
		for range tickChan.C { // Runs code inside loop ~ at specified refresh rate
			resString := ""
			if in.StatisticsType == "viewership" {
				resString = s.top10Viewers()
			} else if in.StatisticsType == "duration" {
				resString = s.top10Duration()
			} else if in.StatisticsType == "mute" {
				resString = s.top10Mute()
			} else if in.StatisticsType == "SMA" {
				fmt.Printf("in.SmaChannel")
				fmt.Printf(in.SmaChannel)
				resString = s.sma(in.SmaChannel, in.SmaLength) // Have to compile proto file again

			}

			err := stream.Send(&pb.NotificationMessage{Top10: resString})
			if err != nil {
				return err
			}
		}
	}
}

func main() {
	parseFlags()
	grpcServer := grpc.NewServer()
	startZapServer()

	server := &SubscribeServer{logger: zlog.NewAdvancedZapLogger()}
	go server.recordAll() // Record all zaps and store in logger

	pb.RegisterSubscriptionServer(grpcServer, server)

	listener, err := net.Listen("tcp", *endpoint)
	if err != nil {
		log.Fatalf("net.listen error: %v\n", err)
	}

	fmt.Printf("Preparing to serve incoming requests...\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("Error with gRPC serve. Quitting...")
	}
}
