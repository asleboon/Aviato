package zlog

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// TODO: Move mutelogger to viewerslogger and rename to advancedlogger??
// TODO: Implement handling in gRPC server and client
// TODO: Save HDMI-status? Don't think it is needed

type DurationMuted struct {
	duration map[string]*channelMute // Key: channel name
	lock     sync.Mutex
}

type channelMute struct {
	// viewersIP  []string   	// Slice of muted IP address watching the channel. Don't think this is needed
	duration      time.Duration // Total mute duration
	viewers       int           // Current number of viewers
	maxMutedTime  time.Time     // Date and time with highest number of muted views
	maxMutedNum   int           // Time with highest number of muted views
	numberOfMuted int           // Current number of muted viewers
}

type viewerVolume struct {
	viewer map[string]*val // Key: IP address
	lock   sync.Mutex
}

type val struct {
	channel   string    // Previous channel watched
	volume    int       // Previous volume value
	mute      int       // Previous mute value
	muteStart time.Time // Time when mute was started
}

// Global variables
var prevVol *viewerVolume

// NewViewersZapLogger initializes a new map for storing views per channel.
// Viewers adhere Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	dm := DurationMuted{duration: make(map[string]*channelMute, 0)}
	prevVol = &viewerVolume{viewer: make(map[string]*val, 0)}
	return &dm
}

// LogStatus handles a status chnage
func (dm *DurationMuted) LogStatus(s chzap.StatusChange) {
	//
	prev, ipExists := prevVol.viewer[s.IP]
	channelStats, channelExists := dm.duration[prev.channel]

	statusType, statusValue := strings.Split(s.Status, ":")[0], strings.Split(s.Status, ":")[1]
	statusValue = strings.TrimSpace(statusValue)
	statusValueInt, err := strconv.Atoi(strings.TrimSpace(statusValue))

	switch statusType {
	case "Mute_Status":
		if statusValueInt == 1 { // s.Status == "Mute_Status: 1"
			if ipExists {
				prev.mute = 1
				// If previous channel is known - Update channel stats
				if prev.channel != "" {
					// Set mute start time of not already set on this channel
					if prev.muteStart.IsZero() {
						prev.muteStart = s.Time
					}
					if channelExists {
						channelStats.numberOfMuted++
						if channelStats.numberOfMuted > channelStats.maxMutedNum {
							channelStats.maxMutedTime = s.Time
							channelStats.maxMutedNum = channelStats.numberOfMuted
						}
					} else { // Should never happen
						fmt.Printf("Failure: Channel assigned to IP address, but does not exist in DurationMuted.")
					}
				}
			} else { // IP address not previously encountered
				prevVol.viewer[s.IP] = &val{mute: 1}
			}
		} else if statusValueInt == 0 { // s.Status == "Mute_Status: 0"
			if ipExists {
				prev.mute = 0
				// If previous channel is known - Update channel stats
				if prev.channel != "" {
					if channelExists {
						if prev.volume > 0 { // Is this check needed?
							channelStats.numberOfMuted--
							channelStats.duration += prev.muteStart.Sub(s.Time)
						}
					} else { // Should never happen
						fmt.Printf("Failure: Channel assigned to IP address, but does not exist in DurationMuted.")
					}
					prev.muteStart = time.Time{} // Reset mute start time
				}
			} else { // IP address not previously encountered
				prevVol.viewer[s.IP] = &val{mute: 0}
			}
		}
	case "HDMI_Status":
		if statusValueInt == 1 { // s.Status == "HDMI_Status: 1"
			// If channel exists and (prev.volume == 0 || prev.mute == 1)
			// --> New muted viewer on channel
		} else if statusValueInt == 0 { // s.Status == "HDMI_Status: 0"
			// If channel exists and (prev.volume == 0 || prev.mute == 1)
			// --> One less muted viewer on channel
		}

	case "Volume":
		if statusValueInt > 0 && statusValueInt <= 100 { // s.Status == "Volume: >0"
			// If channel exists and (prev.volume == 0 || prev.mute == 0)
			// --> One less muted viewer on channel
		} else if statusValueInt == 0 { // s.Status == "Volume: 0"
			// If channel exists and (prev.mute == 0)
			// --> New muted viewer on channel
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
