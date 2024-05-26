package main

import "C"
import (
	"encoding/json"
	"fmt"
)

type eventFunctionExecutedFn func(map[string]interface{}) string

var eventHandlers = map[string]eventFunctionExecutedFn{
	"function_executed": OnFunctionExecuted,
	"method_executed":   OnMethodExecuted,
}

//export OnEvent
func OnEvent(eventJson string) (outputJson string) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("[AIKIDO-GO] Recovered from panic:", r)
			outputJson = "{}"
		}
	}()

	fmt.Println("[AIKIDO-GO] OnEvent:", eventJson)

	var event map[string]interface{}
	err := json.Unmarshal([]byte(eventJson), &event)
	if err != nil {
		panic(fmt.Sprintf("Error parsing JSON: %s", err))
	}

	eventName := MustGetFromMap[string](event, "event")
	data := MustGetFromMap[map[string]interface{}](event, "data")

	CheckIfKeyExists(eventHandlers, eventName)

	return eventHandlers[eventName](data)
}

func main() {}
