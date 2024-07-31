package aikido_types

import "sync"

type StatsDataType struct {
	StatsMutex sync.Mutex

	StartedAt       int64
	Requests        int
	RequestsAborted int
	Attacks         int
	AttacksBlocked  int
}
