package zlog

import (
	"fmt"
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Logger type contains datastructure for logging viewers, duration and mute stats from a set-top box.
// Remark: HDMI_Status and Volume not considered in mute and durationlogger
type Logger struct {
	viewers  map[string]int           // Key: channel name, value: current number of viewers (viewerslogger)
	duration map[string]time.Duration // Key: channel name, value: total viewtime (durationlogger)
	prevZap  map[string]chzap.ChZap   // Key: IP address, value: previous zap (used for durationlogger)
	prevMute map[string]*muteStat     // Key: IP address, value: previous mute (used for mutelogger)
	mute     map[string]*chanMute     // Key: channel name, value: mute stats (mutelogger)
	sma      map[string][]*smaStats   // Key: channel name, value: views
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
	//channel   string    // Previous channel watched. Do not need this
	mute      string    // Previous mute value
	muteStart time.Time // Time when mute was started
}

// NewAdvancedZapLogger creates a new logger. Adhere Zaplogger interface
func NewAdvancedZapLogger() AdvZapLogger {
	lg := Logger{
		viewers:  make(map[string]int, 0),
		duration: make(map[string]time.Duration, 0),
		prevZap:  make(map[string]chzap.ChZap, 0),
		prevMute: make(map[string]*muteStat, 0),
		mute:     make(map[string]*chanMute, 0),
		sma:      make(map[string][]*smaStats, 0),
	}
	return &lg
}

// LogZap updates loggers when a new zap event is received
func (lg *Logger) LogZap(z chzap.ChZap) {
	lg.lock.Lock() // or (*lg)?
	defer lg.lock.Unlock()
	logZapViewers(z, lg)  // Update viewers data structure
	logZapDuration(z, lg) // Update durationdata structure
	//logZapMute(z, lg)     // Update mute data structure
}

func logZapViewers(z chzap.ChZap, lg *Logger) {
	// ToChan handling
	count, exists := lg.viewers[z.ToChan]
	if exists {
		lg.viewers[z.ToChan] = count + 1
	} else {
		lg.viewers[z.ToChan] = 1
	}

	// FromChan handling
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
	if ipExists == true { // If no prev mute values exist for this IP, do nothing

		// From channel handling
		fromChannelStats, channelExists := lg.mute[z.FromChan]
		if !channelExists { // Initialize chanMute strcut for this channel
			lg.mute[z.FromChan] = &chanMute{muteViewers: make(map[string]bool, 0)}
			fromChannelStats = lg.mute[z.FromChan]
		}
		if prev.mute == "1" {
			fromChannelStats.numberOfMute--
			fromChannelStats.duration += z.Duration(prev.muteStart)
		}

		// To channel handling
		toChannelStats, channelExists := lg.mute[z.ToChan]
		if !channelExists {
			lg.mute[z.ToChan] = &chanMute{muteViewers: make(map[string]bool, 0)}
			fromChannelStats = lg.mute[z.ToChan]
		}
		if prev.mute == "1" {
			// Increment number of mutes on channel and set maxMuteNum and maxMuteTime if true
			toChannelStats.numberOfMute++
			if toChannelStats.numberOfMute > toChannelStats.maxMuteNum {
				toChannelStats.maxMuteTime = time.Now()
				toChannelStats.maxMuteNum = toChannelStats.numberOfMute
			}
			// Update prev mute
			prev.mute = "1"
			prev.muteStart = z.Time

			// Add IP address to map of IP addresses that have viewed this channel muted
			toChannelStats.muteViewers[z.IP] = true
		}
	}
}

// LogStatus updates loggers when a new status event is received
func (lg *Logger) LogStatus(s chzap.StatusChange) {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	logStatusMute(s, lg) // Update mute data structure
}

func logStatusMute(s chzap.StatusChange, lg *Logger) {
	pZap, pZapExists := lg.prevZap[s.IP]
	if pZapExists { // No updates if no zap events on previous zap registred on this IP address
		channelStats, channelExists := lg.mute[pZap.ToChan]
		prev, ipExists := lg.prevMute[s.IP]

		if s.Status == "Mute_Status: 1" || s.Status == "Mute_Status: 0" {
			if !ipExists { // Create new muteStat struct for IP if not present
				lg.prevMute[s.IP] = &muteStat{}
				prev = lg.prevMute[s.IP]
			}
			if !channelExists { // Create new chanMute struct for IP if not present
				lg.mute[pZap.ToChan] = &chanMute{muteViewers: make(map[string]bool, 0)}
				channelStats = lg.mute[pZap.ToChan]
			}
		}

		if s.Status == "Mute_Status: 1" {
			if channelExists {
				lg.mute[pZap.ToChan].numberOfMute++
				if lg.mute[pZap.ToChan].numberOfMute > channelStats.maxMuteNum {
					lg.mute[pZap.ToChan].maxMuteTime = time.Now()
					lg.mute[pZap.ToChan].maxMuteNum = channelStats.numberOfMute
				}
				lg.mute[pZap.ToChan].muteViewers[s.IP] = true
			}
			// Update prev mute values
			prev.mute = "1"
			prev.muteStart = s.Time

		} else if s.Status == "Mute_Status: 0" {
			if channelExists {
				channelStats.numberOfMute--
				if !prev.muteStart.IsZero() {
					channelStats.duration += s.Time.Sub(prev.muteStart)
				}
			}
			// Update prev mute value
			prev.mute = "0"
		}
	}
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

func (lg *Logger) ChannelsSMA(channelName string) *map[string][]*smaStats {
	count := lg.viewers[channelName]
	fmt.Printf("count: %q", count)
	output := &smaStats{Views: count, TimeAdded: time.Now()}
	lg.sma[channelName] = append(lg.sma[channelName], output)
	return &lg.sma
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (lg *Logger) ChannelsViewers() []*AdvChannelViewers {
	lg.lock.Lock()
	defer lg.lock.Unlock()
	defer util.TimeElapsed(time.Now(), "ChannelsViewers")

	res := make([]*AdvChannelViewers, 0)
	for channel, viewers := range lg.viewers {
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
	for channel, duration := range lg.duration {
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
	for channel, mute := range lg.mute {
		avgMute := 0.0
		if len(mute.muteViewers) > 0 {
			avgMute = mute.duration.Seconds() / float64(len(mute.muteViewers))
			// TODO: Check that avgMute is desired format
		}
		if avgMute > 0 { // Don't want to include channels without a valid average mute in result
			fmt.Printf("%v", time.Duration(avgMute)*time.Second)
			advChannelMute := AdvChannelMute{Channel: channel, AvgMute: time.Duration(avgMute), MaxMuteTime: mute.maxMuteTime}
			res = append(res, &advChannelMute)
		}
	}
	return res
}
