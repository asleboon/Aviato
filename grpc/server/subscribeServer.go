package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"time"

	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"github.com/uis-dat320-fall18/Aviato/zlog"
	"google.golang.org/grpc"
)

// SubscribeServer exported? Or not exported?
type SubscribeServer struct {
	// kvMap map[string]int // store channels and number of viewers?
	logger zlog.ZapLogger
	lock   sync.Mutex
}

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
			// TODO: Send top 10 list
			fmt.Printf(in.String()) // Only for debug, remove afterwards
		}
	}
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	grpcServer := grpc.NewServer()

	// TODO: Finish. Remove .Output() ?
	// Start zapserver and top 10 calculation
	output, err := exec.Command("go", "run", "-lab f", "../../zapserver").Output()
	if err == nil {
		fmt.Printf("%v", output)
	}
	server := &SubscribeServer{logger: zlog.NewViewersZapLogger()}
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
