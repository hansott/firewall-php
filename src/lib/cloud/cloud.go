package cloud

import (
	"sync"
	"time"
)

var (
	stop chan struct{}
	wg   sync.WaitGroup
)

func StartPollingThread() {
	stop = make(chan struct{})

	ticker := time.NewTicker(time.Second * 30)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				SendHeartbeatEvent()
				return
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func StopPollingThread() {
	close(stop)
	wg.Wait()
}

func Init() {
	SendStartEvent()
	StartPollingThread()
}

func Uninit() {
	StopPollingThread()
}
