// +build !solution

// Leave an empty line above this comment.
//
// Zap Collection Server
package main

import (
	"github.com/uis-dat320/glabs/lab7/zlog"
	"log"
)

// REMARK: This function should return (i.e. it should not block)
func runLab() {
	switch *labnum {
	case "a", "c1", "c2", "d", "e":
		ztore = zlog.NewSimpleZapLogger()
	case "f":
		//TODO activate with new ZapLogger data structure (task f)
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
	//TODO write this method (5p)
}
