package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

type Chartlogger struct {
	views map[string][]*viewTime // Key: Channelname, value: viewtime
	lock  sync.Mutex
}
type viewTime struct {
	times time.Time
	views float64
}

func NewChartLogger() *Chartlogger {
	cl := Chartlogger{views: make(map[string][]*viewTime, 0)}
	return &cl
}

// LogZap updates count for the two channels in the zap
func (cl *Chartlogger) LogZap(z chzap.ChZap) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	// Log views
	vtSlice := cl.views[z.ToChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &viewTime{times: z.Time, views: 1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &viewTime{times: z.Time, views: prevVt.views + 1})
	}

	vtSlice = cl.views[z.FromChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &viewTime{times: z.Time, views: -1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &viewTime{times: z.Time, views: prevVt.views - 1})
	}
}

func GetChartVal(channelName string, cl *Chartlogger) []*viewTime {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.views[channelName]
}

// Entries returns the length of views map (# of channnels)
func (cl *Chartlogger) Entries() int {
	return 0
}

// Viewers return number of viewers for a channel
func (cl *Chartlogger) Viewers(channelName string) int {
	return 0
}

// Channels creates a list of channels in the viewers.
func (cl *Chartlogger) Channels() []string {
	return nil
}

// ChannelsViewers creates a ChannelViewers slice (# of viewers per channel)
func (cl *Chartlogger) ChannelsViewers() []*ChannelViewers {
	return nil
}
