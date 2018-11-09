package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// TODO: Move mutelogger to viewerslogger and rename to advancedlogger??
// TODO: Implement handling in gRPC server and client

type DurationMuted struct {
	duration map[string]channelMute // Key: channel name
	lock     sync.Mutex
}

type channelMute struct {
	viewersIP     []string      // Slice of muted IP address watching the channel
	duration      time.Duration // Total mute duration
	viewers       int           // Current number of viewers
	maxMuted      time.Time     // Time with highest number of muted views
	numberOfMuted int           // Current number of muted viewers
}

type viewerVolume struct {
	duration map[string]prevVal // Key: IP address
	lock     sync.Mutex
}

type prevVal struct {
	volume int // Previous volume value
	mute   int // Previous mute value
}

// NewViewersZapLogger initializes a new map for storing views per channel.
// Viewers adhere Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	dm := DurationMuted{duration: make(map[string]channelMute, 0)}
	return &dm
}

// LogStatus handles a status chnage
func (dm *DurationMuted) LogStatus(s chzap.StatusChange) {
	// 1. If statuschange is Mute_Status: 1 and prev volume > 0 --> New muted viewer on channel
	if s.Status == "Mute_Status: 1" {
		// count, exists := (*ms).mutes[s.IP]

	}

	// 2. If statuschange is Volume = 0 and prev Mute_Status = 0 --> One more muted viewer on channel

	// 3. If statuschange is Mute_Status = 0 and prev volume > 0 --> One less muted viewer on channel

	// 4. If statuschange is HDMI_Status = 0 and (prev volume = 0 or mute_status = 1) --> One less muted viewer on channel

	// 5. If statuschange is HDMI_Status = 1 and (prev volume = 0 or mute_status = 1) --> New muted viewer on channel

	// 6. If statuschange is Volume > 0 and prev volume = 0 and prev mutestatus = 0 --> One less muted viewer on channel

}

// LogZap moves current viewer between channels in DurationMuted if muted
func (dm *DurationMuted) LogZap(z chzap.ChZap) {
	// TODO: Implement
	// If IP in muted list for fromChan, move mute to toChan
	// (*vs).lock.Lock()
	// defer (*vs).lock.Unlock()
}

// func (du *DurationMuted) LogMuted(s chzap.StatusChange) {
// 	prev.lock.Lock()
// 	defer prev.lock.Unlock()
// 	pZap, exists := prev.prevZap[s.IP]
// 	if s.Status == "Volume: 0" {
// 		// This will not work uless
// 		(*du).lock.Lock()
// 		defer (*du).lock.Unlock()
// 		(*du).duration[pZap.ToChan] += newDur // Add duration for channel
// 		delete(prev.prevZap, s.IP)            // Remove prevZap froom this IP
// 	}
// }

// --------------------------------------------------------------------

// Entries returns the length of views map (# of channnels)
func (dm *DurationMuted) Entries() int {
	(*dm).lock.Lock()
	defer (*dm).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Entries")
	return len((*vs).views)
}

// Viewers return number of viewers for a channel
func (dm *DurationMuted) Viewers(channelName string) int {
	// (*dm).lock.Lock()
	// defer (*dm).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "Viewers")

	// count, exists := (*vs).views[channelName]
	// if exists {
	// 	return count
	// }
	return 0 // Not found in views map = 0 zaps
}

// Channels creates a list of channels in the viewers.
func (dm *DurationMuted) Channels() []string {
	// (*dm).lock.Lock()
	// defer (*dm).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "Channels")

	// channels := make([]string, 0)
	// for channel := range (*vs).views {
	// 	channels = append(channels, channel)
	// }
	// return channels
	return nil
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (dm *DurationMuted) ChannelsViewers() []*ChannelViewers {
	// (*dm).lock.Lock()
	// defer (*dm).lock.Unlock()
	// defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	// res := make([]*ChannelViewers, 0)
	// for channel, viewers := range (*vs).views {
	// 	channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
	// 	res = append(res, &channelViewer)
	// }
	// return res
	return nil
}
