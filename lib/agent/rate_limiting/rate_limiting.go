package rate_limiting

import (
	"main/globals"
	"main/log"
	"time"
)

var (
	stop               chan struct{}
	RateLimitingTicker = time.NewTicker(globals.MinRateLimitingIntervalInMs * time.Millisecond)
)

func AdvanceRateLimitingQueues() {
	globals.RateLimitingMutex.Lock()
	defer globals.RateLimitingMutex.Unlock()

	for _, endpoint := range globals.RateLimitingMap {
		if endpoint.Config.WindowSizeInMinutes <= endpoint.Status.NumberOfRequestsPerWindow.Length() {
			// Sliding window is moving, need to substract the entry that goes out of the window
			// Ex: if the window is set to 10 minutes, when another minute passes,
			//     need to remove the number of requests from the entry of 11 minutes ago

			// Get the number of requests for the entry that just dropped out of the sliding window
			numberOfRequestToSubstract := endpoint.Status.NumberOfRequestsPerWindow.Pop()
			if endpoint.Status.TotalNumberOfRequests < numberOfRequestToSubstract {
				// This should never happen, but better to have a check in place
				log.Warnf("More requests to substract (%d) than total number of requests (%d) for endpoint (%v)",
					numberOfRequestToSubstract, endpoint.Status.TotalNumberOfRequests, endpoint)
			} else {
				// Remove the number of requests for the entry that just dropped out of the sliding window from total
				endpoint.Status.TotalNumberOfRequests -= numberOfRequestToSubstract
			}
		}

		// Create a new entry in queue for the current minute
		endpoint.Status.NumberOfRequestsPerWindow.Push(0)
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
