package zlog

import (
	"sync"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

type Chartlogger struct {
	views map[string][]*ViewTime // Key: Channelname, value: viewtime
	lock  sync.Mutex
}
type ViewTime struct {
	times time.Time
	views float64
}

func NewChartLogger() *Chartlogger {
	cl := Chartlogger{views: make(map[string][]*ViewTime, 0)}
	return &cl
}

// LogZap updates count for the two channels in the zap
func (cl *Chartlogger) LogZap(z chzap.ChZap) {
	cl.lock.Lock()
	defer cl.lock.Unlock()

	// Log views
	vtSlice := cl.views[z.ToChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &ViewTime{times: z.Time, views: 1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &ViewTime{times: z.Time, views: prevVt.views + 1})
	}

	vtSlice = cl.views[z.FromChan]
	if len(vtSlice) <= 0 {
		vtSlice = append(vtSlice, &ViewTime{times: z.Time, views: -1})
	} else {
		prevVt := vtSlice[len(vtSlice)-1]
		vtSlice = append(vtSlice, &ViewTime{times: z.Time, views: prevVt.views - 1})
	}
}

func (cl *Chartlogger) GetChartVal(channelName string) []*ViewTime {
	cl.lock.Lock()
	defer cl.lock.Unlock()
	return cl.views[channelName]
}
