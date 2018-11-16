package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	"github.com/uis-dat320-fall18/Aviato/util"
)

// Viewers holds a lock and a map with key: Channel, Value: viewers
type Viewers struct {
	views map[string]int // Key: Channelname, value: Viewers
	lock  sync.Mutex
}

var chartTime []time.Time
var chartViews []float64

// NewViewersZapLogger initializes a new map for storing views per channel.
// Viewers adhere Zaplogger interface.
func NewViewersZapLogger() ZapLogger {
	vs := Viewers{views: make(map[string]int, 0)}
	return &vs
}

// LogZap updates count for the two channels in the zap
func (vs *Viewers) LogZap(z chzap.ChZap) {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()

	// Log views
	count, exists := (*vs).views[z.ToChan]
	if exists {
		(*vs).views[z.ToChan] = count + 1
	} else {
		(*vs).views[z.ToChan] = 1
	}

	count, exists = (*vs).views[z.FromChan]
	if exists {
		(*vs).views[z.FromChan] = count - 1
	} else {
		(*vs).views[z.FromChan] = -1
	}
}
func (vs *Viewers) LogStatus(z chzap.StatusChange) {
}

// Entries returns the length of views map (# of channnels)
func (vs *Viewers) Entries() int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Entries")
	return len((*vs).views)
}

// Viewers return number of viewers for a channel
func (vs *Viewers) Viewers(channelName string) int {
	(*vs).lock.Lock()
	defer (*vs).lock.Unlock()
	defer util.TimeElapsed(time.Now(), "Viewers")

	count, exists := (*vs).views[channelName]
	if exists {
		return count
	}
	chartViews = append(chartViews, float64(count))
	chartTime = append(chartTime, time.Now())
	return 0 // Not found in views map = 0 zaps
}

// Channels creates a list of channels in the viewers.
func (vs *Viewers) Channels() []string {
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
func (vs *Viewers) ChannelsViewers() []*ChannelViewers {
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

func (vs *Viewers) StupidChart() ([]float64, []time.Time) {
	return chartViews, chartTime
}
