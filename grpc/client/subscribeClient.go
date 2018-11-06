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

// TODO: Refactor code when everything is working
// TODO: Implement flags for specifying refreshrate

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	endpoint = flag.String(
		"endpoint",
		"localhost:12111",
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
	fmt.Printf("Connection to gRPC server created: %v", conn)
	defer conn.Close() // Closing connection when the surrounding function return

	client := pb.NewSubscriptionClient(conn)

	stream, err := client.Subscribe(context.Background())
	rate, _ := strconv.ParseUint((*refreshRate), 10, 0)
	msg := &pb.SubscribeMessage{RefreshRate: rate}
	stream.Send(msg)

	for {
		fmt.Println("Waiting for top 10 from server...\n")
		top10, err := stream.Recv()
		// TODO: Notification message handling here!
		if err == io.EOF {
			fmt.Printf("End of file received. Client quitting...")
			return
		} else if err != nil {
			fmt.Printf("Error: %v", err)
			return
		}
		fmt.Printf("%v", top10)
	}
	// stream.CloseSend() Unreachable code

}
