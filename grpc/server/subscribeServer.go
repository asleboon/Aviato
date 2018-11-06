package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"github.com/uis-dat320-fall18/Aviato/zlog"
	"google.golang.org/grpc"
)

// TODO: Top 10 calc
// TODO: Reformat code

// SubscribeServer exported? Or not exported?
type SubscribeServer struct {
	// kvMap map[string]int // store channels and number of viewers?
	logger zlog.ZapLogger
	lock   sync.Mutex
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
		"localhost:12111",
		"Endpoint on which server runs",
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

func startServer() {
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
		eventStr, err := readFromUDP() // Something wrong with readFromUDP

		if err != nil { // ReadFromUDP error check
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			chZap, _, err := chzap.NewSTBEvent(eventStr) // We don't care about statuschange
			if err != nil {                              // NewSTBEvent error check
				fmt.Printf("Error: %v\n", err)
			} else {
				if chZap != nil {
					s.logger.LogZap(*chZap) // Make a copy of pointer value
				}
			}
		}
	}
}

func (s *SubscribeServer) Subscribe(stream pb.Subscription_SubscribeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		tickChan := time.NewTicker(time.Second * time.Duration(in.RefreshRate))
		defer tickChan.Stop()
		for range tickChan.C { // Runs code inside loop ~ at specified refresh rate
			// Create a top 10 map
			channels := s.logger.ChannelsViewers()
			fmt.Printf("%v", channels) // Only for debug, remove afterwards

			// Sort the channelviews, descending
			sort.Slice(channels, func(i, j int) bool {
				return channels[i].Viewers > channels[j].Viewers
			})

			if len(channels) > 10 { // Only want top 10
				channels = channels[:10]
			}
			//msg := make([]string, 0)
			msg := ""
			// Create a string slice with top 10 ??
			counter := 1
			for k, v := range channels {
				fmt.Printf("%v", k)
				str := string(counter) + ". "
				str += v.Channel
				str += ". Viewers: "
				str += string(v.Viewers)
				msg += str + "\n"
				//msg = append(msg, str)
				counter++
			}
			// Send top 10 to subscriber
			stream.Send(&pb.NotificationMessage{Notification: msg})
		}
	}
}

func main() {
	parseFlags()

	grpcServer := grpc.NewServer()
	startServer() // Start zapserver

	// Create new server with viewerslogger
	server := &SubscribeServer{logger: zlog.NewViewersZapLogger()}
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
