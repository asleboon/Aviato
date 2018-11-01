package zlog

import (
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Viewers contains a map with Key: Channelname, Value: Viewers
// Remark: Full zap info not stored here. Use simplelogger for that.
type Viewers map[string]int

// NewViewersZapLogger initializes a new map for storing views per channel. Adheres Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	vs := make(Viewers, 0)
	return &vs
}

// LogZap updates count for the two channels in the zap
func (vs *Viewers) LogZap(z chzap.ChZap) {
	count, exists := (*vs)[z.ToChan]
	if exists {
		(*vs)[z.ToChan] = count + 1
	} else {
		(*vs)[z.ToChan] = 1
	}

	count, exists = (*vs)[z.FromChan]
	if exists {
		(*vs)[z.FromChan] = count - 1
	} else {
		(*vs)[z.FromChan] = -1
	}
}

// Entries returns the length og the Viewers map (# of channnels)
func (vs *Viewers) Entries() int {
	return len(*vs)
}

// Viewers return number of viewers for a channel
func (vs *Viewers) Viewers(channelName string) int {
	count, exists := (*vs)[channelName]
	if exists {
		return count
	}
	return 0
}

// Channels creates a list of channels in the viewers.
func (vs *Viewers) Channels() []string {
	channels := make([]string, 0)
	for channel := range *vs {
		channels = append(channels, channel)
	}
	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (vs *Viewers) ChannelsViewers() []*ChannelViewers {
	defer util.TimeElapsed(time.Now(), "simple.ChannelsViewers")

	res := make([]*ChannelViewers, 0)
	for channel, viewers := range *vs {
		channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &channelViewer)
	}
	return res
}
