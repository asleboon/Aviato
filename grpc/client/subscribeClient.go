package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

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
	refreshRate = flag.String(
		"refreshRate",
		"1",
		"Refresh rate at which the client will get a top 10 channel response from the server",
	)
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func dumpTop10() {	// input stream w/ appropriate input type
}

func main() {
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	conn, err := grpc.Dial(*endpoint, grpc.WithInsecure()) // WithInsecure: Disable transport security connection
	if err != nil {
		log.Fatalf("Error with creating connection to gRPC server: %v", err)
	}
	fmt.Printf("Connection to gRPC server created\n")
	defer conn.Close()

	client := pb.NewSubscriptionClient(conn)

	stream, err := client.Subscribe(context.Background())
	rate, _ := strconv.ParseUint((*refreshRate), 10, 0)
	msg := &pb.SubscribeMessage{RefreshRate: rate}
	stream.Send(msg)

	waitchan := make(chan struct{})	// Wait channel
	// go dumpTop10(stream)
	// TODO: Refactor, create named func
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				fmt.Printf("End of file received. Client quitting...")
				return
			} else if err != nil {
				fmt.Printf("Error: %v", err)
				return
			}
			log.Printf("Top 10")
			fmt.Printf("%v", in.Top10)
		}
	}()
	stream.CloseSend()	// Need this?
	<-waitchan
}
