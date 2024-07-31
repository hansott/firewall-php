package cloud

import (
	"main/globals"
	"time"
)

var (
	stop            chan struct{}
	HeartBeatTicker = time.NewTicker(10 * time.Minute)
)

func StartHeartbeatRoutine() {
	stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-HeartBeatTicker.C:
				SendHeartbeatEvent()
			case <-stop:
				HeartBeatTicker.Stop()
				return
			}
		}
	}()
}

func StopHeartbeatRoutine() {
	close(stop)
}

func Init() {
	SendStartEvent()
	StartHeartbeatRoutine()

	globals.StatsData.StartedAt = GetTime()
}

func Uninit() {
	StopHeartbeatRoutine()
}
