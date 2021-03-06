package main

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/uis-dat320-fall18/Aviato/chzap"
	pb "github.com/uis-dat320-fall18/Aviato/proto"
	"github.com/uis-dat320-fall18/Aviato/zlog"
)

// SubscribeServer includes a logger for zap- and statusevents
type SubscribeServer struct {
	logger zlog.AdvZapLogger
}

// recordAll processes and stores new viewers in Zaplogger
func (s *SubscribeServer) recordAll(udpServer *UDPServer) {
	for {
		eventStr, err := ReadFromUDP(udpServer)
		chZap, stChange, err := chzap.NewSTBEvent(eventStr)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
		} else if chZap != nil {
			s.logger.LogZap(*chZap) // Pass a copy of pointer value
		} else if stChange != nil {
			s.logger.LogStatus(*stChange) // Pass a copy of pointer value
		}
	}
}

func (s *SubscribeServer) top10Viewers() string {
	channels := s.logger.ChannelsViewers() // Map of all channels with number of viewers

	// Sort channels by views, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Viewers > channels[j].Viewers
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}
		top10Str += fmt.Sprintf("%v. %v, viewers: %v", count+1, v.Channel, v.Viewers)
	}
	top10Str += "\n\n" // Easy way to create space between top 10 prints
	return top10Str
}

func (s *SubscribeServer) top10Duration() string {
	channels := s.logger.ChannelsDuration() // Map of all channels with total duration

	// Sort channels by total duration, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Duration > channels[j].Duration
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}
		top10Str += fmt.Sprintf("%v. %v, total duration: %v", count+1, v.Channel, v.Duration)
	}
	top10Str += "\n\n"
	return top10Str
}

func (s *SubscribeServer) top10Mute() string {
	channels := s.logger.ChannelsMute() // Map of all channels with avg. muted duration per viewer

	// Sort channels by avg mute per viewer, descending
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].AvgMute > channels[j].AvgMute
	})

	if len(channels) > 10 { // Only want top 10 channels
		channels = channels[:10]
	}

	// Create top 10 string
	top10Str := ""
	for count, v := range channels {
		if count != 0 {
			top10Str += "\n"
		}

		top10Str += fmt.Sprintf("%v. %v, average muted duration per viewer: %v\n", count+1, v.Channel, v.AvgMute)
		t := v.MaxMuteTime
		top10Str += fmt.Sprintf("Time with highest number of muted viewers: %d-%02d-%02d %02d:%02d:%02d\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	}
	top10Str += "\n\n"
	return top10Str
}

// sma calculates the simple moving for a channel within a timeframe given by the client
func (s *SubscribeServer) sma(smaChannel string, smaLength uint64) string {
  sumViewers, count := float64(0), float64(0)
	sma := s.logger.ChannelsSMA(smaChannel) // returns []*smaStats
	if sma == nil {
		return ""
	}

	lastTime := sma[len(sma)-1].TimeAdded
	for i := len(sma) - 2; i > 0; i-- {
		// Check if the views should be included in the calculation

		if lastTime.Sub(sma[i].TimeAdded) < (time.Duration(smaLength) * time.Second) {
			sumViewers += float64(sma[i].Views)
			count++
		} else {
			break // because we loop backwords, no more times should be added
		}
	}
	if count == 0 {
		return fmt.Sprintf("Simple moving average for %s: %d\n", smaChannel, 0)
	}
	return fmt.Sprintf("Simple moving average for %s the last %d seconds is %.2f:\n", smaChannel, smaLength, sumViewers/count)
}

// Subscribe handles a client subscription request
func (s *SubscribeServer) Subscribe(stream pb.Subscription_SubscribeServer) error {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		tickChan := time.NewTicker(time.Second * time.Duration(in.RefreshRate))
		defer tickChan.Stop()
		for range tickChan.C { // Runs code inside loop ~ at specified refresh rate
			resString := ""
			if in.StatisticsType == "viewership" {
				resString = s.top10Viewers()
			} else if in.StatisticsType == "duration" {
				resString = s.top10Duration()
			} else if in.StatisticsType == "mute" {
				resString = s.top10Mute()
			} else if in.StatisticsType == "sma" {
				resString = s.sma(in.SmaChannel, in.SmaLength)
			}
			err := stream.Send(&pb.NotificationMessage{Statistics: resString})
			if err != nil {
				return err
			}
		}
	}
}
