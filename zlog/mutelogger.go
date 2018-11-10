package zlog

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

// TODO: Move mutelogger to viewerslogger and rename to advancedlogger??
// TODO: Implement handling in gRPC server and client

type DurationMuted struct {
	duration map[string]*channelMute // Key: channel name
	lock     sync.Mutex
}

type channelMute struct {
	duration      time.Duration // Total mute duration
	viewers       int           // Current number of viewers
	maxMutedTime  time.Time     // Date and time with highest number of muted views
	maxMutedNum   int           // Time with highest number of muted views
	numberOfMuted int           // Current number of muted viewers
}

type viewerStats struct {
	viewer map[string]*val // Key: IP address
	lock   sync.Mutex
}

type val struct {
	channel   string    // Previous channel watched
	volume    string    // Previous volume value
	mute      string    // Previous mute value
	muteStart time.Time // Time when mute was started
}

// Global variables
var prevVol *viewerStats

func NewMuteZapLogger() ZapLogger {
	dm := DurationMuted{duration: make(map[string]*channelMute, 0)}
	prevVol = &viewerStats{viewer: make(map[string]*val, 0)}
	return &dm
}

// LogStatus handles a status chnage
func (dm *DurationMuted) LogStatus(s chzap.StatusChange) {
	prevVol.lock.Lock()
	defer prevVol.lock.Unlock()
	prev, ipExists := prevVol.viewer[s.IP]
	if !ipExists { // Create new entry if IP never been logged before
		prevVol.viewer[s.IP] = &val{}
		prev = prevVol.viewer[s.IP]
	}

	dm.lock.Lock()
	defer dm.lock.Unlock()
	channelStats, channelExists := dm.duration[prev.channel]

	statusType, statusValue := strings.Split(s.Status, ":")[0], strings.Split(s.Status, ":")[1]
	statusValue = strings.TrimSpace(statusValue)

	switch statusType {
	case "Mute_Status":
		if statusValue == "1" { // TV muted
			prev.mute = "1"
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
		} else if statusValue == "0" { // TV unmuted
			prev.mute = "0"
			// If previous channel is known - Update channel stats
			if prev.channel != "" {
				if channelExists {
					if prev.volume != "0" || prev.volume != "" { // Is this check needed?
						channelStats.numberOfMuted--
						if !prev.muteStart.IsZero() { // If this is not true, muteStart never set
							channelStats.duration += prev.muteStart.Sub(s.Time)
						}
					}
				} else { // Should never happen
					fmt.Printf("Failure: Channel assigned to IP address, but does not exist in DurationMuted.")
				}
				prev.muteStart = time.Time{} // Reset mute start time
			}
		}
	case "HDMI_Status":
		if statusValue == "1" { // TV connected
			// If channel exists and (prev.volume == 0 || prev.mute == 1)
			// --> New muted viewer on channel
		} else if statusValue == "0" { // TV disconnected
			// If channel exists and (prev.volume == 0 || prev.mute == 1)
			// --> One less muted viewer on channel
		}

	case "Volume":
		if statusValue == "0" { // Volume adjusted to 0
			// If channel exists and (prev.mute == 0)
			// --> New muted viewer on channel
		} else { // Volume adjusted to 1-100
			// If channel exists and (prev.volume == 0 || prev.mute == 0)
			// --> One less muted viewer on channel
		}
	}
}

// LogZap updates muted statistics if TV is muted
func (dm *DurationMuted) LogZap(z chzap.ChZap) {
	prevVol.lock.Lock()
	defer prevVol.lock.Unlock()
	prev, ipExists := prevVol.viewer[z.IP]
	if !ipExists {
		prevVol.viewer[z.IP] = &val{}
		prev = prevVol.viewer[z.IP]
	}

	dm.lock.Lock()
	defer dm.lock.Unlock()
	fromChannelStats, channelExists := dm.duration[z.FromChan]
	if channelExists {
		fromChannelStats.viewers--
		if prev.mute == "1" || prev.volume == "0" {
			fromChannelStats.numberOfMuted--
			fromChannelStats.duration += z.Duration(prev.muteStart)
		}
	} else {
		dm.duration[z.ToChan] = &channelMute{viewers: -1}
	}

	toChannelStats, channelExists := dm.duration[z.ToChan]
	if channelExists {
		toChannelStats.viewers++
		prev.channel = z.ToChan
		if prev.mute == "1" || prev.volume == "0" {
			toChannelStats.numberOfMuted++
			if prev.muteStart.IsZero() {
				prev.muteStart = z.Time
			}
		}
	} else {
		dm.duration[z.ToChan] = &channelMute{viewers: 1}
	}
}

// --------------------------------------------------------------------

// Entries returns the length of views map (# of channnels)
func (dm *DurationMuted) Entries() int {
	//(*dm).lock.Lock()
	//defer (*dm).lock.Unlock()
	//defer util.TimeElapsed(time.Now(), "Entries")
	//return len((*vs).views)
	return 0
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
