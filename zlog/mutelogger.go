package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

/*Extend the gRPC publish/subscribe server to include the top-10 most muted channels
(in terms of average muted duration per viewer) until that time point (refresh point),
along with the time of the day during which each of these channels had the highest number of viewers with muted TV.*/

/*
Need to store previous time from mute to unmuted
Need to include what channel it was muted on
*/

type DurationMuted struct {
	duration map[string]time.Duration // Key: channel name, value: duration muted
	lock     sync.Mutex
}

type MutedViewers struct {
	views map[string]int // Key: Channelname, value: mutedViewers
	lock  sync.Mutex
}

// NewViewersZapLogger initializes a new map for storing views per channel.
// Viewers adhere Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	vs := MutedViewers{views: make(map[string]int, 0)}
	return &vs
}

// also need to know what channel it was on
func (ms *MutedViewers) LogStatus(s chzap.StatusChange) {
	if s.Status == "Mute_Status: 0" {
		count, exists := (*ms).mutes[s.IP]
		// loop throught zap events to find out what channel he was whatching?
		// What if someone has a muted tv and then changes channel?
		//
	}
}

// LogZap updates count for the two channels in the zap
func (vs *MutedViewers) LogZap(z chzap.ChZap) {
	// (*vs).lock.Lock()
	// defer (*vs).lock.Unlock()

	// // Log views
	// count, exists := (*vs).views[z.ToChan]
	// if exists {
	// 	(*vs).views[z.ToChan] = count + 1
	// } else {
	// 	(*vs).views[z.ToChan] = 1
	// }

	// count, exists = (*vs).views[z.FromChan]
	// if exists {
	// 	(*vs).views[z.FromChan] = count - 1
	// } else {
	// 	(*vs).views[z.FromChan] = -1
	// }
}

// Entries returns the length of views map (# of channnels)
func (vs *MutedViewers) Entries() int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Entries")
	return len((*vs).views)
}

// Viewers return number of viewers for a channel
func (vs *MutedViewers) Viewers(channelName string) int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}
	return 0 // Not found in views map = 0 zaps
}

// Channels creates a list of channels in the viewers.
func (vs *MutedViewers) Channels() []string {
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
func (vs *MutedViewers) ChannelsViewers() []*ChannelViewers {
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
