package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

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
		"localhost:12111",
		"Endpoint to which client connects",
	)
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

/*
func subTop10(client pb.NewSubscriptionClient, rate uint32) {
	stream, err := client.Subscribe(context.Background())
	waitc := make(chan struct{})
	msg := &pb.SubscribeMessage{RefreshRate: rate}
	for {
		fmt.Println("sleeping")
		time.Sleep(2 * time.Second)
		fmt.Println("sending msg...")
		stream.Send(msg)
	}
	<-waitc
	stream.CloseSend()
}*/

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

	// go subTop10()

	stream, err := client.Subscribe(context.Background())
	waitc := make(chan struct{})
	msg := &pb.SubscribeMessage{RefreshRate: 2}
	go func() {
		for {
			fmt.Println("sleeping")
			time.Sleep(2 * time.Second)
			fmt.Println("sending msg...")
			stream.Send(msg)
		}
	}()
	<-waitc
	stream.CloseSend()

	// reader := bufio.NewReader(os.Stdin)               // Why do we need a reader?
	// fmt.Printf("reader: %q", reader)                  // And why is this needed?

}
