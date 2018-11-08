package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

// TODO: Implement locks
// TODO: Implement in grpc server:
// Run duration logger and add extra field in Subscribe msg

// Use pointers and locks when data is access concurrently
// https://bit.ly/2Qyj5Zr

// DurationChan stores total viewtime per channel
type DurationChan struct {
	duration map[string]time.Time // Key: channel name, value: total duration(viewtime)
	lock     sync.Mutex
}

// prevZap stores previous channel
type prevZapIP struct {
	prevZap map[string]chzap.ChZap // Key: IP address, value: prev zap
	lock    sync.Mutex
}

type zap struct {
	channel string
	start   time.Time
}

type globalStats struct {
	duration time.Duration // Total duration(viewtime)
	zaps     int           // Total number of zaps
	lock     sync.Mutex
}

// Global variables
var global *globalStats
var prev *prevZapIP

// NewDurationZapLogger duration logger data structure
// DurationChan adheres Zaplogger interface.
func NewDurationZapLogger() ZapLogger {
	du := DurationChan{duration: make(map[string]time.Time, 0)}
	prev = &prevZapIP{prevZap: make(map[string]chzap.ChZap, 0)}
	global = &globalStats{}
	return &du
}

// LogZap updates duration counter
func (du *DurationChan) LogZap(z chzap.ChZap) {
	(*du).lock.Lock()
	defer (*du).lock.Unlock()
	// TODO: Implement
}

// Log status removes previous zap from IP address if TV is turned off
func (du *DurationChan) LogStatus(s chzap.StatusChange) {
	(*du).lock.Lock()
	defer (*du).lock.Unlock()
	// TODO: Implement
}

// Entries returns the length of views map (# of channnels)
func (du *DurationChan) Entries() int {
	(*du).lock.Lock()
	defer (*du).lock.Unlock()
	return len((*du).duration)
}

// Viewers return number of viewers for a channel
func (du *DurationChan) Viewers(channelName string) int {
	//(*vs).lock.Lock()
	//defer (*vs).lock.Unlock()
	/*defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}*/
	return 0
}

// Channels creates a list of channels in the viewers.
func (du *DurationChan) Channels() []string {
	//(*vs).lock.Lock()
	//defer (*vs).lock.Unlock()
	/*defer util.TimeElapsed(time.Now(), "Channels")

	channels := make([]string, 0)
	for channel := range (*vs).views {
		channels = append(channels, channel)
	}
	return channels*/
	return nil
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (du *DurationChan) ChannelsViewers() []*ChannelViewers {
	// (*vs).lock.Lock()
	// defer (*vs).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	// res := make([]*ChannelViewers, 0)
	// for channel, viewers := range (*vs).views {
	// 	channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
	// 	res = append(res, &channelViewer)
	// }
	// return res
	return nil
}
