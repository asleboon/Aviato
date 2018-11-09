package zlog

import (
	"strconv"
	"strings"
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
	viewer map[string]*val // Key: IP address
	lock   sync.Mutex
}

type val struct {
	channel string // Previous channel watched
	volume  int    // Previous volume value
	mute    int    // Previous mute value
}

// Global variables
var prevVol *viewerVolume

// NewViewersZapLogger initializes a new map for storing views per channel.
// Viewers adhere Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	dm := DurationMuted{duration: make(map[string]channelMute, 0)}
	prevVol = &viewerVolume{viewer: make(map[string]*val, 0)}
	return &dm
}

// LogStatus handles a status chnage
func (dm *DurationMuted) LogStatus(s chzap.StatusChange) {
	prev, exists := prevVol.viewer[s.IP]
	if exists {
		prevVolume := prev.volume // Remove?
		prevMute := prev.mute     // Remove?
	}

	statusType, statusValue := strings.Split(s.Status, ":")[0], strings.Split(s.Status, ":")[1]
	statusValue = strings.TrimSpace(statusValue)
	statusValueInt, err := strconv.Atoi(strings.TrimSpace(statusValue))

	switch statusType {

	case "Mute_Status":
		// Case 1, s.Status == "Mute_Status: 1"
		// And not in viewersIP slice --> New muted viewer on channel
		if statusValueInt == 1 {
			prevVol.viewer[s.IP].mute = 1

			// Case 2, s.Status == "Mute_Status: 0"
			// 3. If statuschange is Mute_Status = 0 and prev volume > 0 --> One less muted viewer on channel
		} else if statusValueInt == 0 {
			prevVol.viewer[s.IP].mute = 0
		}

	case "HDMI_Status":
		// Case 3, s.Status == "HDMI_Status: 1"
		// And (prev volume = 0 or mute_status = 1) --> New muted viewer on channel
		if statusValueInt == 1 {

			// Case 4, s.Status == "HDMI_Status: 0"
			// And (prev volume = 0 or mute_status = 1) --> One less muted viewer on channel
		} else if statusValueInt == 0 {

		}

	case "Volume":
		// Case 5, s.Status == "Volume: >0"
		// And prev volume = 0 and prev mutestatus = 0 --> One less muted viewer on channel
		if statusValueInt > 0 && statusValueInt <= 100 {

			// Case 6, s.Status == "Volume: 0"
			// And prev Mute_Status = 0 --> One more muted viewer on channel
		} else if statusValueInt == 0 {

		}
	}
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
