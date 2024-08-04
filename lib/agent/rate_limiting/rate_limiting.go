package rate_limiting

import (
	"main/globals"
	"time"
)

var (
	stop               chan struct{}
	RateLimitingTicker = time.NewTicker(globals.MinRateLimitingIntervalInMs * time.Millisecond)
)

func AdvanceRateLimitingQueues() {
	globals.RateLimitingMutex.Lock()
	defer globals.RateLimitingMutex.Unlock()

	for _, v := range globals.RateLimitingMap {
		if v.Config.WindowSizeInMinutes <= v.Status.NumberOfRequestPerWindow.Length() {
			v.Status.TotalNumberOfRequests -= v.Status.NumberOfRequestPerWindow.Pop()
		}
		v.Status.NumberOfRequestPerWindow.Push(0)
	}
}

func StartRateLimitingRoutine() {
	AdvanceRateLimitingQueues()

	stop = make(chan struct{})

	go func() {
		for {
			select {
			case <-RateLimitingTicker.C:
				AdvanceRateLimitingQueues()
			case <-stop:
				RateLimitingTicker.Stop()
				return
			}
		}
	}()
}

func StopRateLimitingRoutine() {
	close(stop)
}

func Init() {
	StartRateLimitingRoutine()
}

func Uninit() {
	StopRateLimitingRoutine()
}
