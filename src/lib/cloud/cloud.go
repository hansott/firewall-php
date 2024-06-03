package cloud

import (
	"fmt"
	"io"
	"main/log"
	"net/http"
	"sync"
	"time"
)

var endpoint = ""

var (
	stop chan struct{}
	wg   sync.WaitGroup
)

func MakeGetRequest(url string) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	log.Info("Response status:", resp.Status)
	log.Info("Response body:", string(body))
}

func StartConfigThread() {
	stop = make(chan struct{})

	ticker := time.NewTicker(time.Second * 5)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				MakeGetRequest(endpoint)
			case <-stop:
				ticker.Stop()
				return
			}
		}
	}()
}

func StopConfigThread() {
	close(stop)
	wg.Wait()
}

func Init(configuredEndpoint string) {
	endpoint = configuredEndpoint
	StartConfigThread()
}

func Uninit() {
	StopConfigThread()
}
