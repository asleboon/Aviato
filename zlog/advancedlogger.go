package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Logger type contains datastructure for logging viewers, duration and mute stats from a set-top box.
type Logger struct {
	viewers  map[string]int           // Key: channel name, value: number of viewers (viewerslogger)
	duration map[string]time.Duration // Key: channel name, value: total viewtime (durationlogger)
	prevZap  map[string]chzap.ChZap   // Key: IP address, value: previous zap
	mute     map[string]chanMute      // Key: channel name, value: mute stats (mutelogger)
	prevMute map[string]muteStat      // Key: IP address, value: previous mute
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

	logZapViewers(z, lg)  // go?
	logZapDuration(z, lg) // go?
	logZapMute(z, lg)     // go?
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
		prevVol.viewer[z.IP] = &val{}
		prev = lg.prevMute[z.IP]
	}

	fromChannelStats, channelExists := lg.mute[z.FromChan]
	if channelExists {
		lg.viewers[z.ToChan]--
		if prev.mute == "1" || prev.volume == "0" {
			fromChannelStats.numberOfMuted--
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
			toChanStats.numberOfMuted++
			if prev.muteStart.IsZero() {
				prev.muteStart = z.Time
			}
		}
	} else {
		lg.mute[z.ToChan] = chanMute{}
	}
}

// LogStatus updates loggers when a new status event is received
func (lg *Logger) LogStatus(s chzap.StatusChange) {
	lg.lock.Lock() // or (*lg)?
	defer lg.lock.Unlock()

	logStatusMute(s, lg)
	logStatusDuration(s, lg)
}

func logStatusDuration(s chzap.StatusChange, lg *Logger) {
	// TODO: Implement
}

func logStatusMute(s chzap.StatusChange, lg *Logger) {
	// TODO: Implement
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

	count, exists := lg.viewers[channelName]
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
	for channel := range lg.viewers {
		channels = append(channels, channel)
	}
	return channels
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (lg *Logger) ChannelsViewers() []*ChannelViewers {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	res := make([]*ChannelViewers, 0)
	for channel, viewers := range (*lg).viewers {
		channelViewer := ChannelViewers{Channel: channel, Viewers: viewers}
		res = append(res, &channelViewer)
	}
	return res
}
