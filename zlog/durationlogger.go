package zlog

import (
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Do we need to implement locks?
// TODO: Discuss data structure
// TODO: Implement in grpc server:
// Run duration logger and add extra field in Subscribe msg

// lastZap pointer or not pointer?
type lastZapChan map[string]*lastZap // Key: IP address, value: channel name and start time
type lastZap struct {
	channel string
	start   time.Time
}

// Pointer or not pointer?
type totDurChan map[string]*channelStats // Key: channel name, value: total duration(viewtime) and viewers(current)
type channelStats struct {
	duration time.Time
	viewers  int
}

type globalStats struct {
	duration time.Time // Total duration(viewtime)
	zaps     int       // Total number of zaps
}

// NewDurationZapLogger duration logger data structure
// Adheres Zaplogger interface.
func NewDurationZapLogger() ZapLogger {
	// vs := Viewers{views: make(map[string]int, 0)} REMOVE
	// return &vs REMOVE
	return nil
}

// LogZap updates count for the two channels in the zap
func (vs *Viewers) LogZap(z chzap.ChZap) {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()

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

// Entries returns the length of views map (# of channnels)
func (vs *Viewers) Entries() int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	return len((*vs).views)
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
