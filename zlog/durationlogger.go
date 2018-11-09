package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

// TODO: Move durationlogger to viewerslogger and rename to advancedlogger??
// TODO: Implement or remove global stats

// Use pointers and locks when data is access concurrently
// https://bit.ly/2Qyj5Zr

// DurationChan stores total viewtime per channel
type DurationViewtime struct {
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
	du := DurationViewtime{duration: make(map[string]time.Duration, 0)}
	prev = &prevZapIP{prevZap: make(map[string]chzap.ChZap, 0)}
	global = &globalStats{}
	return &du
}

// LogZap updates duration counter and prevZap based on new zap event
func (dw *DurationViewtime) LogZap(z chzap.ChZap) {
	prev.lock.Lock()
	defer prev.lock.Unlock()
	pZap, exists := prev.prevZap[z.IP]

	if exists {
		newDur := z.Duration(pZap.Time) // Duration between previous and this zap on IP

		(*dw).lock.Lock()
		defer (*dw).lock.Unlock()
		(*dw).duration[pZap.ToChan] += newDur // Add duration for channel
	}

	prev.prevZap[z.IP] = z // Update prevZap to include new zap event for IP
}

// LogStatus stores duration and removes previous zap from IP address if TV is turned off
func (dw *DurationViewtime) LogStatus(s chzap.StatusChange) {
	if s.Status == "HDMI_Status: 0" {
		prev.lock.Lock()
		defer prev.lock.Unlock()
		pZap, exists := prev.prevZap[s.IP]
		if exists {
			newDur := pZap.Duration(s.Time)

			(*dw).lock.Lock()
			defer (*dw).lock.Unlock()
			(*dw).duration[pZap.ToChan] += newDur // Add duration for channel
			delete(prev.prevZap, s.IP)            // Remove prevZap froom this IP
		}
	}
	if s.Status == "HDMI_Status: 1" {
		// TODO: Add to viewer based on prev zap
		// Need two maps: Valid prevZap and allPrevZap
	}
}

// ----------------------------------------------------------------------------------

// Entries returns the length of views map (# of channels)
func (dw *DurationViewtime) Entries() int {
	(*dw).lock.Lock()
	defer (*dw).lock.Unlock()
	return len((*dw).duration)
}

// Viewers return number of viewers for a channel
func (dw *DurationViewtime) Viewers(channelName string) int {
	//(*dw).lock.Lock()
	//defer (*dw).lock.Unlock()
	/*defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}*/
	return 0
}

// Channels creates a list of channels in the prevChannels map.
// Doesn´t really make sense here
func (dw *DurationViewtime) Channels() []string {
	// (*dw).lock.Lock()
	// defer (*dw).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "Channels")

	// channels := make([]string, 0)
	// for channel := range (*dw).duration {
	// 	channels = append(channels, channel)
	// }
	return nil //channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (dw *DurationViewtime) ChannelsViewers() []*ChannelViewers {
	// (*dw).lock.Lock()
	// defer (*dw).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	// res := make([]*ChannelViewers, 0)
	// for channel, duration := range (*dw).duration {

	// 	channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
	// 	res = append(res, &channelViewer)
	// }
	return nil // return res
}
