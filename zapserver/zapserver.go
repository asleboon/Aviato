// +build !solution

// Leave an empty line above this comment.
//
// Zap Collection Server
package main

import (
	"fmt"
	"log"
	"net"

	"github.com/uis-dat320-fall18/assignments/lab6/zlog"
)

// REMARK: This function should return (i.e. it should not block)
func runLab() {
	switch *labnum {
	case "a", "c1", "c2", "d", "e":
		ztore = zlog.NewSimpleZapLogger()
	case "f":
		// TODO activate with new ZapLogger data structure (task f)
		// ztore = zlog.NewViewersZapLogger()
	}
	switch *labnum {
	case "a":
		//TODO write code for dumping zap events to console
		// go dumpAll()
	case "c1":
		//TODO write code for recording and showing # of viewers on NRK1
		// go recordAll()
		// go showViewers("NRK1")
	case "c2":
		//TODO write code for task c2
	case "d":
		//TODO write code for task d
	case "e":
		//TODO write code for task e
	case "f":
		//TODO write code for task f
	}
}

// REMARK: This function should return (i.e. it should not block)
func startServer() {
	log.Println("Starting ZapServer...")
	//TODO write this method (5p)C
	// Build UDP address
	addr, _ := net.ResolveUDPAddr("udp", "224.0.1.130:10000")

	// Create connection
	UDPConn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("NewUDPServer: Error creating UDP connection")
	}

	// Serve UDPServer
	for {
		buf := make([]byte, 1024)           // Make a buffer used to store bytes read from UDP
		n, _, _ := UDPConn.ReadFromUDP(buf) // n = Number of bytes read, addr = UDP connection address, err = Error
		txt := string(buf[:n])
		fmt.Printf("New response: %v", txt)
	}
}
