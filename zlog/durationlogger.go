package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

// Do we need to implement locks?
// TODO: Discuss data structure
// TODO: Implement in grpc server:
// Run duration logger and add extra field in Subscribe msg

// We should use pointers if the map is accessed concurrently
// Don't use it unless it is necessesary.
// https://bit.ly/2Qyj5Zr
// DurationChan stores total viewtime per channel
type DurationChan struct {
	duration map[string]time.Time // Key: channel name, value: total duration(viewtime)
	lock     sync.Mutex
}

// prevZap stores previous channel
type prevZapIP map[string]*prevZap // Key: IP address, value: channel name and start time
type prevZap struct {
	channel string
	start   time.Time
}

type globalStats struct {
	duration time.Time // Total duration(viewtime)
	zaps     int       // Total number of zaps
}

// NewDurationZapLogger duration logger data structure
// Adheres Zaplogger interface.
func NewDurationZapLogger() ZapLogger {
	du := DurationChan{duration: make(map[string]time.Time, 0)}
	//prevZap{}
	return &du
}

// LogZap updates count for the two channels in the zap
func (du *DurationChan) LogZap(z chzap.ChZap) {
	//(*du).lock.Lock()
	//defer (*vs).lock.Unlock()

	// Log views
	/*count, exists := (*vs).views[z.ToChan]
	if exists {
		(*vs).views[z.ToChan] = count + 1
	} else {
		(*vs).views[z.ToChan] = 1
	}

	count, exists = (*vs).views[z.FromChan]
	if exists {
		(*vs).views[z.FromChan] = count - 1
	} else {
		(*vs).views[z.FromChan] = -1
	}
	*/
}

// Entries returns the length of views map (# of channnels)
func (du *DurationChan) Entries() int {
	//(*vs).lock.Lock()
	//defer (*vs).lock.Unlock()
	//return len((*vs).views)
	return 0
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
