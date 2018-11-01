package main

import (
	"context"
	"sync"

	//"google.golang.org/grpc"
	pb "github.com/uis-dat320-fall18/Aviato/proto"
)

type SubscribeServer struct {
	kv    map[string]int // store channels and number of viewers?
	mutex sync.Mutex
}

// subscribe.proto: rpc Subscribe(stream SubscribeMessage) returns (stream NotificationMessage)
func (ss *SubscribeServer) Subscribe(ctx context.Context, req *pb.SubscribeMessage) (*pb.NotificationMessage, error) {
	return nil, nil
}

func main() {
}
