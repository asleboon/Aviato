package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Logger type contains datastructure for logging viewers, duration and mute stats from a set-top box.
type Logger struct {
	view     map[string]int           // Key: channel name, value: number of viewers (viewerslogger)
	duration map[string]time.Duration // Key: channel name, value: total viewtime (durationlogger)
	prevZap  map[string]*chzap.ChZap  // Key: IP address, value: previous zap
	mute     map[string]*chanMute     // Key: channel name, value: mute stats (mutelogger)
	prevMute map[string]*muteStat     // Key: IP address, value: previous mute
	lock     sync.Mutex
}

type chanMute struct {
	duration      time.Duration // Total mute duration
	maxMutedTime  time.Time     // Date and time with highest number of muted views
	maxMutedNum   int           // Time with highest number of muted views
	numberOfMuted int           // Current number of muted viewers
}

type muteStat struct {
	channel   string    // Previous channel watched
	volume    string    // Previous volume value
	mute      string    // Previous mute value
	muteStart time.Time // Time when mute was started
}

// NewAdvancedZapLogger creates a new logger. Adhere Zaplogger interface
func NewAdvancedZapLogger() ZapLogger {
	lg := Logger{
		view:     make(map[string]int, 0),
		duration: make(map[string]time.Duration, 0),
		prevZap:  make(map[string]*chzap.ChZap, 0),
		mute:     make(map[string]*chanMute, 0),
		prevMute: make(map[string]*muteStat, 0),
	}
	return &lg
}

// LogZap updates loggers when a new zap event is received
func (lg *Logger) LogZap(z chzap.ChZap) {
	(*lg).lock.Lock()
	defer (*lg).lock.Unlock()

	// Log views
	count, exists := (*lg).view[z.ToChan]
	if exists {
		(*lg).view[z.ToChan] = count + 1
	} else {
		(*lg).view[z.ToChan] = 1
	}

	count, exists = (*lg).view[z.FromChan]
	if exists {
		(*lg).view[z.FromChan] = count - 1
	} else {
		(*lg).view[z.FromChan] = -1
	}
}

// LogStatus updates loggers when a new status event is received
func (lg *Logger) LogStatus(s chzap.StatusChange) {

}

// Entries returns the length of the views map (# of channels)
func (lg *Logger) Entries() int {
	(*lg).lock.Lock()
	defer (*lg).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Entries")
	return len((*lg).view)
}

// Viewers return number of viewers for a channel
func (lg *Logger) Viewers(channelName string) int {
	(*lg).lock.Lock()
	defer (*lg).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*lg).view[channelName]
	if exists {
		return count
	}
	return 0 // Not found in views map = 0 zaps
}

// Channels creates a list of channels in the viewers.
func (lg *Logger) Channels() []string {
	(*lg).lock.Lock()
	defer (*lg).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Channels")

	channels := make([]string, 0)
	for channel := range (*lg).view {
		channels = append(channels, channel)
	}
	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (lg *Logger) ChannelsViewers() []*ChannelViewers {
	(*lg).lock.Lock()
	defer (*lg).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	res := make([]*ChannelViewers, 0)
	for channel, viewers := range (*lg).view {
		channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &channelViewer)
	}
	return res
}
