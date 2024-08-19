package cloud

import (
	"main/globals"
	"main/utils"
	"time"
)

var (
	HeartbeatRoutineChannel     = make(chan struct{})
	HeartBeatTicker             = time.NewTicker(10 * time.Minute)
	ConfigPollingRoutineChannel = make(chan struct{})
	ConfigPollingTicker         = time.NewTicker(1 * time.Minute)
)

func Init() {
	SendStartEvent()
	utils.StartPollingRoutine(HeartbeatRoutineChannel, HeartBeatTicker, SendHeartbeatEvent)
	utils.StartPollingRoutine(ConfigPollingRoutineChannel, ConfigPollingTicker, CheckConfigUpdatedAt)

	globals.StatsData.StartedAt = GetTime()
}

func Uninit() {
	utils.StopPollingRouting(HeartbeatRoutineChannel)
}
