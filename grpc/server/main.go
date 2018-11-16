package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"runtime/pprof"

	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"github.com/uis-dat320-fall18/Aviato/zlog"
	"google.golang.org/grpc"
)

type UDPServer struct {
	conn *net.UDPConn
}

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
	memprofile = flag.String(
		"memprofile",
		"",
		"write memory profile to this file",
	)
	cpuprofile = flag.String(
		"cpuprofile",
		"",
		"write cpu profile to this file",
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

func main() {
	parseFlags()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)

	grpcServer := grpc.NewServer()
	udpServer, err := NewUDPServer("224.0.1.130:10000") // have to specify this specifically unless we use another flag
	subscribeServer := &SubscribeServer{logger: zlog.NewAdvancedZapLogger()}

	go subscribeServer.recordAll() // Record all zaps and store in logger

	pb.RegisterSubscriptionServer(grpcServer, subscribeServer)

	listener, err := net.Listen("tcp", *endpoint)
	if err != nil {
		log.Fatalf("net.listen error: %v\n", err)
	}

	fmt.Printf("Preparing to serve incoming requests...\n")
	go grpcServer.Serve(listener)

	// Here we wait for CTRL-C or some other kill signal
	s := <-signalChan
	fmt.Println("Server stopping on", s, "signal")
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Println("Saved memory profile")
		fmt.Println("Analyze with: go tool pprof $GOPATH/bin/zapserver", *memprofile)
	}
}
