package cloud

import (
	"time"
)

var (
	stop            chan struct{}
	HeartBeatTicker = time.NewTicker(10 * time.Minute)
)

func StartPollingThread() {
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

func StopPollingThread() {
	close(stop)
}

func Init() {
	SendStartEvent()
	StartPollingThread()
}

func Uninit() {
	StopPollingThread()
}
