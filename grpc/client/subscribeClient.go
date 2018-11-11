package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"google.golang.org/grpc"
)

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	endpoint = flag.String(
		"endpoint",
		"localhost:1994", // Changed port from std to 1994 to avoid problems during testing.
		"Endpoint to which client connects",
	)
	refreshRate = flag.Uint64(
		"refreshRate",
		1,
		"Refresh rate at which the client will get a top 10 channel response from the server. Default: 1 second.",
	)
	statisticsType = flag.String(
		"statisticsType",
		"viewership",
		"Statistics type for which this client want to subscribe for. Options: viewership (default) , muted or duration.",
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

// dumpTop10 receives stream values and prints
func dumpTop10(stream pb.Subscription_SubscribeClient, sType string) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			fmt.Printf("End of file received. Client quitting...")
			return
		} else if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		log.Printf("Top 10: " + sType)
		fmt.Printf("%v", in.Top10)
	}
}

func main() {
	parseFlags()

	conn, err := grpc.Dial(*endpoint, grpc.WithInsecure()) // WithInsecure: Disable transport security connection
	if err != nil {
		log.Fatalf("Error with creating connection to gRPC server: %v", err)
	}
	fmt.Printf("\nConnection to gRPC server created\n\n")
	defer conn.Close()

	client := pb.NewSubscriptionClient(conn)
	stream, err := client.Subscribe(context.Background()) // context.Background(): Non-nil, empty Context
	if err != nil {
		log.Fatalf("Client failed to subscribe: %v", err)
	}

	err = stream.Send(&pb.SubscribeMessage{RefreshRate: *refreshRate, StatisticsType: *statisticsType}) // Send subscribe msg to gRPC server
	stream.CloseSend()                                                                                  // Client will not send more messages on the stream

	waitchan := make(chan struct{}) // Wait channel so main does not return
	go dumpTop10(stream, *statisticsType)
	<-waitchan
}
