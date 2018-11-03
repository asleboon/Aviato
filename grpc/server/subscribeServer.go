package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"google.golang.org/grpc"
)

// SubscribeServer exported? Or not exported?
type SubscribeServer struct {
	kv    map[string]int // store channels and number of viewers?
	mutex sync.Mutex
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
		"Endpoint on which server runs or to which client connects",
	)
)

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

/* Remove?
// subscribe.proto: rpc Subscribe(stream SubscribeMessage) returns (stream NotificationMessage)
func (ss *SubscribeServer) Subscribe(ctx context.Context, req *pb.SubscribeMessage) (*pb.NotificationMessage, error) {
	return nil, nil
}
*/

// Look at SubscriptionServer API in subscribe.pb.go
type SubscriptionServer interface {
	Subscribe(Subscription_SubscribeServer) error
}

func RegisterSubscriptionServer(s *grpc.Server, srv SubscriptionServer) {
	// s.RegisterService(&_Subscription_serviceDesc, srv)
}

func _Subscription_Subscribe_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(SubscriptionServer).Subscribe(&subscriptionSubscribeServer{stream})
}

type Subscription_SubscribeServer interface {
	// Send(*NotificationMessage) error
	// Recv() (*SubscribeMessage, error)
	grpc.ServerStream
}

type subscriptionSubscribeServer struct {
	grpc.ServerStream
}

func (x *subscriptionSubscribeServer) Send(m *pb.NotificationMessage) error {
	return x.ServerStream.SendMsg(m)
}

func (x *subscriptionSubscribeServer) Recv() (*pb.SubscribeMessage, error) {
	m := new(pb.SubscribeMessage)
	if err := x.ServerStream.RecvMsg(m); err != nil {
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

	listener, err := net.Listen("tcp", *endpoint)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Listener started on %v\n", *endpoint)
	}

	server := new(SubscribeServer)
	server.kv = make(map[string]int)
	grpcServer := grpc.NewServer()

	// Must create proto first. cmd: protoc --go_out=plugins=grpc:. subscribe.proto
	//pb.RegisterKeyValueServiceServer(grpcServer, server)

	fmt.Printf("Preparing to serve incoming requests.\n")
	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
