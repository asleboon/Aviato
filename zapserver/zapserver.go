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
	//"github.com/uis-dat320-fall18/assignments/lab6/zlog" REMOVE
)

var conn *net.UDPConn
var err error

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
		go dumpAll()
	case "c1":
		go recordAll()
		go showViewers("NRK1")
	case "c2":
		go recordAll()
		go showViewers("TV2 Norge")
	case "d":
		//TODO write code for task d ???
	case "e":
		//TODO write code for task e
		// go recordAll()
		// go top10Viewers()
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
	conn, err = net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("NewUDPServer: Error creating UDP connection")
	}
}

func readFromUDP() (string, error) {
	buf := make([]byte, 1024)          // Make a buffer used to store bytes read from UDP
	n, _, err := conn.ReadFromUDP(buf) // n = Number of bytes read
	str := string(buf[:n])
	return str, err
}

// dumpAll reads new STB events and prints to console
func dumpAll() {
	for {
		eventStr, err := readFromUDP()
		if err == nil { // ReadFromUDP error check
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			fmt.Printf("Dumped response: %v\n", eventStr)
		}
	}
}

// recordAll processes and stores new viewers in Zaplogger
func recordAll() {
	for {
		eventStr, err := readFromUDP()
		if err == nil { // ReadFromUDP error check
			fmt.Printf("ReadFromUDP: error: %v\n", err)
		} else {
			chZap, _, err := chzap.NewSTBEvent(eventStr) // We don't care about statuschange

			if err == nil { // NewSTBEvent error check
				fmt.Printf("Error: %v\n", err)
			} else {
				if chZap != nil {
					ztore.LogZap(*chZap) // Make a copy of pointer value
					fmt.Printf("Stored zap: %v\n", eventStr)
				}
			}
		}
	}
}

// showViewers compute number of viewers on channel and prints to console every second
func showViewers(chName string) {
	tickChan := time.NewTicker(time.Second)
	defer tickChan.Stop()

	for range tickChan.C { // Runs code inside loop ~ every second
		views := ztore.Viewers(chName)
		fmt.Printf("No. of viewers on %s is now %d\n", chName, views)
	}
}

// top10Viewers compute
func top10Viewers() {
	ztore.Channels()
}
