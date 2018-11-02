package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Viewers contains a map with Key: Channelname, Value: Viewers
// Remark: Full zap info not stored here. Use simplelogger for that.
// TODO: Implement mutex lock on map to prevent errors when running go routines concurrently accessing the map.
//type Viewers map[string]int

type Viewers struct {
	zaps  []chzap.ChZap //??
	views map[string]int
	lock  sync.Mutex
}

// data
// lock

// NewViewersZapLogger initializes a new map for storing views per channel. Adheres Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	vs := Viewers{views: make(map[string]int, 0)}
	return &vs
}

// LogZap updates count for the two channels in the zap
func (vs *Viewers) LogZap(z chzap.ChZap) {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()

	// Log zap
	(*vs).zaps = append((*vs).zaps, z)

	// Log views
	count, exists := (*vs).views[z.ToChan]
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
}

// Entries returns the length og the Viewers map (# of channnels)
func (vs *Viewers) Entries() int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	return len((*vs).zaps)
}

// Viewers return number of viewers for a channel
func (vs *Viewers) Viewers(channelName string) int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}
	return 0
}

// Channels creates a list of channels in the viewers.
func (vs *Viewers) Channels() []string {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Channels")

	channels := make([]string, 0)
	for channel := range (*vs).views {
		channels = append(channels, channel)
	}
	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (vs *Viewers) ChannelsViewers() []*ChannelViewers {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	res := make([]*ChannelViewers, 0)
	for channel, viewers := range (*vs).views {
		channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &channelViewer)
	}
	return res
}
