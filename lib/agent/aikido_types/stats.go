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

type RateLimitingConfig struct {
	MaxRequests         int
	WindowSizeInMinutes int
}

type RateLimitingStatus struct {
	NumberOfRequestPerWindow Queue
	TotalNumberOfRequests    int
}

type RateLimitingKey struct {
	Method string
	Route  string
}

type RateLimitingValue struct {
	Config RateLimitingConfig
	Status RateLimitingStatus
}
