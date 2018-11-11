package zlog

import (
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
}

type AdvChannelViewers struct {
	Channel string
	Viewers int
}

type AdvChannelDuration struct {
	Channel  string
	Duration int
}

type AdvChannelMute struct {
	Channel string
	Mute    int
}
