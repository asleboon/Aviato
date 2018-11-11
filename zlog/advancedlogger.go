package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Logger type contains datastructure for logging viewers, duration and mute stats from a set-top box.
type Logger struct {
	viewers  map[string]int           // Key: channel name, value: current number of viewers (viewerslogger)
	duration map[string]time.Duration // Key: channel name, value: total viewtime (durationlogger)
	prevZap  map[string]chzap.ChZap   // Key: IP address, value: previous zap
	mute     map[string]chanMute      // Key: channel name, value: mute stats (mutelogger)
	prevMute map[string]muteStat      // Key: IP address, value: previous mute
	lock     sync.Mutex
}

type chanMute struct {
	duration     time.Duration   // Total mute duration
	maxMuteTime  time.Time       // Date and time with highest number of muted views
	maxMuteNum   int             // Time with highest number of muted views
	numberOfMute int             // Current number of muted viewers
	muteViewers  map[string]bool // Map of IP addresses that have viewed this channel muted
}

type muteStat struct {
	channel   string    // Previous channel watched
	volume    string    // Previous volume value
	mute      string    // Previous mute value
	muteStart time.Time // Time when mute was started
}

// NewAdvancedZapLogger creates a new logger. Adhere Zaplogger interface
func NewAdvancedZapLogger() AdvZapLogger {
	lg := Logger{
		viewers:  make(map[string]int, 0),
		duration: make(map[string]time.Duration, 0),
		prevZap:  make(map[string]chzap.ChZap, 0),
		mute:     make(map[string]chanMute, 0),
		prevMute: make(map[string]muteStat, 0),
	}
	return &lg
}

// LogZap updates loggers when a new zap event is received
func (lg *Logger) LogZap(z chzap.ChZap) {
	lg.lock.Lock() // or (*lg)?
	defer lg.lock.Unlock()
	// Can we run these as go routines? (Remember locks!)
	logZapViewers(z, lg)  // Update viewers data structure
	logZapDuration(z, lg) // Update durationdata structure
	logZapMute(z, lg)     // Update mute data structure
}

func logZapViewers(z chzap.ChZap, lg *Logger) {
	count, exists := (*lg).viewers[z.ToChan]
	if exists {
		lg.viewers[z.ToChan] = count + 1
	} else {
		lg.viewers[z.ToChan] = 1
	}

	count, exists = lg.viewers[z.FromChan]
	if exists {
		lg.viewers[z.FromChan] = count - 1
	} else {
		lg.viewers[z.FromChan] = -1
	}
}

func logZapDuration(z chzap.ChZap, lg *Logger) {
	pZap, exists := lg.prevZap[z.IP]

	if exists {
		newDur := z.Duration(pZap.Time)    // Duration between previous and this zap on IP
		lg.duration[pZap.ToChan] += newDur // Add duration for channel
	}
	lg.prevZap[z.IP] = z // Update prevZap to include new zap event for IP
}

func logZapMute(z chzap.ChZap, lg *Logger) {
	prev, ipExists := lg.prevMute[z.IP]
	if !ipExists {
		lg.prevMute[z.IP] = muteStat{}
		prev = lg.prevMute[z.IP]
	}

	fromChannelStats, channelExists := lg.mute[z.FromChan]
	if channelExists {
		lg.viewers[z.ToChan]--
		if prev.mute == "1" || prev.volume == "0" {
			fromChannelStats.numberOfMute--
			fromChannelStats.duration += z.Duration(prev.muteStart)
		}
	} else {
		lg.mute[z.ToChan] = chanMute{}
	}

	toChanStats, channelExists := lg.mute[z.ToChan]
	if channelExists {
		lg.viewers[z.ToChan]++
		prev.channel = z.ToChan
		if prev.mute == "1" || prev.volume == "0" {
			toChanStats.numberOfMute++
			if prev.muteStart.IsZero() {
				prev.muteStart = z.Time
			}
		}
		_, ipMuteExists := lg.mute.muteViewers[z.IP]
		if !ipMuteExists {
			lg.mute.muteViewers[z.IP] = true
		}
	} else {
		lg.mute[z.ToChan] = chanMute{}
	}
}

// LogStatus updates loggers when a new status event is received
func (lg *Logger) LogStatus(s chzap.StatusChange) {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	// Can we run these as go routines? (Remember locks!)
	logStatusDuration(s, lg) // Update duration data structure
	logStatusMute(s, lg)     // Update mute data structure
}

func logStatusDuration(s chzap.StatusChange, lg *Logger) {
	// TODO: Implement. Partly implemented in durationlogger.go
}

func logStatusMute(s chzap.StatusChange, lg *Logger) {
	// TODO: Implement. Partly implemented in mutelogger.go
}

// Entries returns the length of the views map (# of channels)
func (lg *Logger) Entries() int {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Entries")
	return len(lg.viewers)
}

// Viewers return number of viewers for a channel
func (lg *Logger) Viewers(channelName string) int {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*lg).viewers[channelName]
	if exists {
		return count
	}
	return 0 // Not found in views map = 0 zaps
}

// Channels creates a list of channels in the viewers.
func (lg *Logger) Channels() []string {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Channels")

	channels := make([]string, 0)
	for channel := range (*lg).viewers {
		channels = append(channels, channel)
	}
	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (lg *Logger) ChannelsViewers() []*AdvChannelViewers {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	res := make([]*AdvChannelViewers, 0)
	for channel, viewers := range (*lg).viewers {
		advChannelViewer := AdvChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &advChannelViewer)
	}
	return res
}

// ChannelsDuration creates a ChannelDuration slice (total duration per channel)
func (lg *Logger) ChannelsDuration() []*AdvChannelDuration {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsDuration")

	res := make([]*AdvChannelDuration, 0)
	for channel, duration := range (*lg).duration {
		advChannelDuration := AdvChannelDuration{Channel: channel, Duration: duration}
		res = append(res, &advChannelDuration)
	}
	return res
}

// ChannelsMute creates a ChannelMute slice (avg. muted duration per viewer per channel)
func (lg *Logger) ChannelsMute() []*AdvChannelMute {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsMute")

	// Create slice with avg. mute duration per viewer and time of the day with highest number of muted viewers
	res := make([]*AdvChannelMute, 0)
	for channel, mute := range (*lg).mute {
		avgMute := 0
		if len(mute.muteViewers) > 0 {
			avgMute = int(mute.duration.Seconds()) / len(mute.muteViewers)
		}
		advChannelMute := AdvChannelMute{Channel: channel, AvgMute: avgMute, MaxMuteTime: mute.maxMuteTime}
		res = append(res, &advChannelMute)
	}
	return res
}
