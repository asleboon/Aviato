package zlog

import (
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
)

type AdvZapLogger interface {
	LogZap(z chzap.ChZap)
	LogStatus(s chzap.StatusChange)
	Entries() int
	Viewers(channelName string) int
	Channels() []string
	ChannelsViewers() []*AdvChannelViewers
	ChannelsDuration() []*AdvChannelDuration
	ChannelsMute() []*AdvChannelMute
	ChannelsSMA(channelName string) []*SmaStats
}

type SmaStats struct {
	Views     int
	TimeAdded time.Time
}

type AdvChannelViewers struct {
	Channel string
	Viewers int
}

type AdvChannelDuration struct {
	Channel  string
	Duration time.Duration
}

type AdvChannelMute struct {
	Channel     string
	AvgMute     time.Duration
	MaxMuteTime time.Time
}
