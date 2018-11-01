package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"google.golang.org/grpc"
)

const addr = "localhost:xxxxx"

func subscribe(s pb.SubscriptionClient, ch string, rate int) (string, error) {
	//
	return "test", nil
}

func main() {
	// copied from lab3
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("reader: %q", reader)
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not subscribe %v", err)
	}
	fmt.Printf("conn: %q", conn)
}
