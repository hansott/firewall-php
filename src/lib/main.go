package main

import "C"
import (
	"encoding/json"
	"fmt"
	"log"
)

type eventFunctionExecutedFn func(map[string]interface{}) string

var eventHandlers = map[string]eventFunctionExecutedFn{
	"function_executed": OnFunctionExecuted,
}

//export OnEvent
func OnEvent(eventJson string) string {
	fmt.Println("[AIKIDO-GO] OnEvent:", eventJson)

	var event map[string]interface{}
	err := json.Unmarshal([]byte(eventJson), &event)
	if err != nil {
		log.Fatalf("Error parsing JSON: %s", err)
	}

	eventName := MustGetFromMap[string](event, "event")
	data := MustGetFromMap[map[string]interface{}](event, "data")

	ExitIfKeyDoesNotExistInMap(eventHandlers, eventName)

	return eventHandlers[eventName](data)
}

func main() {}
