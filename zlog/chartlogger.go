package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

type Chartlogger struct {
	v    map[string][]*ViewTime // Key: Channelname, value: viewtime
	lock sync.Mutex
}
type ViewTime struct {
	Times time.Time
	Views float64
}

func NewChartLogger() *Chartlogger {
	cl := Chartlogger{v: make(map[string][]*ViewTime, 0)}
	return &cl
}

// LogZap updates count for the two channels in the zap
func (cl *Chartlogger) LogZap(z chzap.ChZap) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	// Log views
	_, exists := cl.v[z.ToChan]
	if !exists {
		cl.v[z.ToChan] = []*ViewTime{&ViewTime{Times: z.Time, Views: 1}}
	} else {
		prevVt := cl.v[z.ToChan][len(cl.v[z.ToChan])-1]
		cl.v[z.ToChan] = append(cl.v[z.ToChan], &ViewTime{Times: z.Time, Views: prevVt.Views + 1})
	}

	_, exists = cl.v[z.FromChan]
	if !exists {
		cl.v[z.FromChan] = []*ViewTime{&ViewTime{Times: z.Time, Views: -1}}
	} else {
		prevVt := cl.v[z.FromChan][len(cl.v[z.FromChan])-1]
		cl.v[z.FromChan] = append(cl.v[z.FromChan], &ViewTime{Times: z.Time, Views: prevVt.Views - 1})
	}
}

func (cl *Chartlogger) GetChartVal(channelName string) []*ViewTime {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.v[channelName]
}
