package main

import (
	"context"
	"flag"
	"fmt"
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
		"localhost:12111",
		"Endpoint to which client connects",
	)
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

/* Remove?
func subscribe(s pb.SubscriptionClient, ch string, rate int) (string, error) {
	//
	return "test", nil
}
*/

// Look at SubscriptionClient API in subscribe.pb.go
type SubscriptionClient interface {
	Subscribe(ctx context.Context, opts ...grpc.CallOption) (Subscription_SubscribeClient, error)
}

type subscriptionClient struct {
	cc *grpc.ClientConn
}

func NewSubscriptionClient(cc *grpc.ClientConn) pb.SubscriptionClient {
	//return &subscriptionClient{cc}
	return nil // Remove
}

func (c *subscriptionClient) Subscribe(ctx context.Context, opts ...grpc.CallOption) (Subscription_SubscribeClient, error) {
	/*
		stream, err := c.cc.NewStream(ctx, &_Subscription_serviceDesc.Streams[0], "/proto.Subscription/Subscribe", opts...)
		if err != nil {
			return nil, err
		}
		x := &subscriptionSubscribeClient{stream}
		return x, nil
	*/
	return nil, nil // Remove
}

type Subscription_SubscribeClient interface {
	Send(*pb.SubscribeMessage) error
	Recv() (*pb.NotificationMessage, error)
	grpc.ClientStream
}

type subscriptionSubscribeClient struct {
	grpc.ClientStream
}

func (x *subscriptionSubscribeClient) Send(m *pb.SubscribeMessage) error {
	return x.ClientStream.SendMsg(m)
}

func (x *subscriptionSubscribeClient) Recv() (*pb.NotificationMessage, error) {
	m := new(pb.NotificationMessage)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Main func inspired by lab3
func main() {
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	// reader := bufio.NewReader(os.Stdin)               // Why do we need a reader?
	// fmt.Printf("reader: %q", reader)                  // And why is this needed?
	conn, err := grpc.Dial(*endpoint, grpc.WithInsecure()) // WithInsecure: Disable transport security connection
	if err != nil {
		log.Fatalf("Error with creating connection to gRPC server: %v", err)
	}
	fmt.Printf("Connection to gRPC server created: %v", conn)
	defer conn.Close() // Closing connection when the surrounding function return
}
