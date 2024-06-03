package cloud

import (
	"fmt"
	"io/ioutil"
	"main/log"
	"net/http"
	"sync"
	"time"
)

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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading body:", err)
		return
	}

	log.Info("Response status:", resp.Status)
	log.Info("Response body:", string(body))
}

func StartConfigThread(url string) {
	stop = make(chan struct{})

	ticker := time.NewTicker(time.Second * 5)
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ticker.C:
				MakeGetRequest(url)
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
