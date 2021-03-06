package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"

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
		"rate",
		1,
		"Refresh rate at which the client will get a response from the server in seconds. Default: 1",
	)
	statisticsType = flag.String(
		"type",
		"viewership",
		"Statistics type which this client want to subscribe to. Options: viewership(default), mute, duration or sma",
	)
	smaChannel = flag.String(
		"smaChan",
		"NRK1",
		"Simple moving average will be calculated for this Channel",
	)
	smaLength = flag.Uint64(
		"smaLen",
		10,
		"Interval for the simple moving average",
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
		if in.Statistics != "" {
			log.Printf("Statistics type: " + sType)
			fmt.Printf("%v", in.Statistics)
		} else {
			log.Print("Something is wrong with your command")
			os.Exit(1)
		}
	}
}

func main() {
	parseFlags()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)

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

	err = stream.Send(&pb.SubscribeMessage{RefreshRate: *refreshRate, StatisticsType: *statisticsType, SmaChannel: *smaChannel, SmaLength: *smaLength}) // Send subscribe msg to gRPC server
	stream.CloseSend()                                                                                                                                  // Client will not send more messages on the stream

	go dumpTop10(stream, *statisticsType)

	// Here we wait for CTRL-C or some other kill signal
	s := <-signalChan
	fmt.Println("Client stopping on", s, "signal")
}
