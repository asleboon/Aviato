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
	duration map[string]time.Duration // Key: channel name, value: total duration(viewtime)
	lock     sync.Mutex
}

// prevZap stores previous channel
type prevZapIP struct {
	prevZap map[string]chzap.ChZap // Key: IP address, value: prev zap
	lock    sync.Mutex
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
	du := DurationChan{duration: make(map[string]time.Duration, 0)}
	prev = &prevZapIP{prevZap: make(map[string]chzap.ChZap, 0)}
	global = &globalStats{}
	return &du
}

// LogZap updates duration counter
func (du *DurationChan) LogZap(z chzap.ChZap) {
	prev.lock.Lock()
	defer prev.lock.Unlock()
	pZap, exists := prev.prevZap[z.IP]

	if exists {
		newDur := z.Duration(pZap.Time) // Duration between previous and this zap on IP

		(*du).lock.Lock()
		defer (*du).lock.Unlock()
		(*du).duration[pZap.ToChan] += newDur // Add duration for channel
	}
	prev.prevZap[z.IP] = z
}

// LogStatus stores duration and removes previous zap from IP address if TV is turned off
func (du *DurationChan) LogStatus(s chzap.StatusChange, z chzap.ChZap) {
	prev.lock.Lock()
	defer prev.lock.Unlock()
	pZap, exists := prev.prevZap[z.IP]
	if s.Status == "HDMI_Status: 0" {
		if exists {
			newDur := z.Duration(pZap.Time)

			(*du).lock.Lock()
			defer (*du).lock.Unlock()
			(*du).duration[pZap.ToChan] += newDur

			// remove previous zap
			prev.prevZap[z.IP] = chzap.ChZap{}
		}
	}
}

// Entries returns the length of views map (# of channels)
func (du *DurationChan) Entries() int {
	(*du).lock.Lock()
	defer (*du).lock.Unlock()
	return len((*du).duration)
}

// Viewers return number of viewers for a channel
func (du *DurationChan) Viewers(channelName string) int {
	//(*du).lock.Lock()
	//defer (*du).lock.Unlock()
	/*defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}*/
	return 0
}

// Channels creates a list of channels in the prevChannels map.
// Doesn´t really make sense here
func (du *DurationChan) Channels() []string {
	// (*du).lock.Lock()
	// defer (*du).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "Channels")

	// channels := make([]string, 0)
	// for channel := range (*du).duration {
	// 	channels = append(channels, channel)
	// }
	return nil //channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (du *DurationChan) ChannelsViewers() []*ChannelViewers {
	// (*du).lock.Lock()
	// defer (*du).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	// res := make([]*ChannelViewers, 0)
	// for channel, duration := range (*du).duration {

	// 	channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
	// 	res = append(res, &channelViewer)
	// }
	return nil // return res
}
