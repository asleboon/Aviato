// Zap Collection Server
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/uis-dat320-fall18/Aviato/zlog"
)

var (
	maddr      = flag.String("mcast", "224.0.1.130:10000", "multicast ip:port")
	labnum     = flag.String("lab", "c1", "which lab exercise to run")
	showHelp   = flag.Bool("h", false, "show this help message and exit")
	memprofile = flag.String("memprofile", "", "write memory profile to this file")
	cpuprofile = flag.String("cpuprofile", "", "write cpu profile to this file")
)

var ztore zlog.ZapLogger
var ztoreGraph *zlog.Chartlogger

func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

func parseFlags() {
	flag.Usage = Usage
	flag.Parse()
	if *showHelp {
		flag.Usage()
		os.Exit(0)
	}
}

func main() {
	parseFlags()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Kill, os.Interrupt)

	server, _ := NewUDPServer(*maddr)
	server.runLab()

	// Here we wait for CTRL-C or some other kill signal
	<-signalChan
	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		fmt.Println("Saved memory profile")
		fmt.Println("Analyze with: go tool pprof $GOPATH/bin/zapserver", *memprofile)
	}
}
