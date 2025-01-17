package machine

import "time"

type IMachineConfigAdapter interface {
	GameStartTimeout() time.Duration
	GameStartMonitor() time.Duration
	GetNatsUrl() string
}
