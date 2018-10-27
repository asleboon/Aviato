// +build !solution

// Leave an empty line above this comment.
//
// Zap Collection Server
package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/zlog"
	//"github.com/uis-dat320-fall18/assignments/lab6/zlog"
)

var conn *net.UDPConn

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
		go dumpAll()
	case "c1":
		//TODO write code for recording and showing # of viewers on NRK1
		go recordAll()
		go showViewers("NRK1")
	case "c2":
		//TODO write code for task c2
		// go recordAll()
		// go showViewers("TV2 Norge")
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
	// Build UDP address
	addr, _ := net.ResolveUDPAddr("udp", "224.0.1.130:10000")

	// Create connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("NewUDPServer: Error creating UDP connection")
	}
}

// func runServer(conn *net.UDPConn) {
// 	for {
// 		buf := make([]byte, 1024)        // Make a buffer used to store bytes read from UDP
// 		n, _, _ := conn.ReadFromUDP(buf) // n = Number of bytes read
// 		txt := string(buf[:n])
// 		fmt.Printf("New response: %v\n", txt)
// 	}
// }

// dumpAll reads new STB events and prints to console
func dumpAll() {
	for {
		buf := make([]byte, 1024)        // Make a buffer used to store bytes read from UDP
		n, _, _ := conn.ReadFromUDP(buf) // n = Number of bytes read
		eventStr := string(buf[:n])
		fmt.Printf("Dumped response: %v\n", eventStr)
	}
}

// recordAll processes and stores new viewers in Zaplogger
func recordAll() {
	for {
		buf := make([]byte, 1024)        // Make a buffer used to store bytes read from UDP
		n, _, _ := conn.ReadFromUDP(buf) // n = Number of bytes read
		eventStr := string(buf[:n])
		fmt.Printf("Recorded response: %v\n", eventStr)
		event := chzap.ChZap.NewSTBevent(eventStr)
		ztore.LogZap(event)
	}
}

// showViewers compute number of viewers on channel and prints to console every second
func showViewers(chName string) {
	// Ticker sends the time every second. Also adjust intervals for slow recievers
	tickChan := time.NewTicker(time.Second)
	defer tickChan.Stop()

	for tick := range tickChan.C {
		views := ztore.Viewers(chName)
		fmt.Printf("%v", views)
	}
}
