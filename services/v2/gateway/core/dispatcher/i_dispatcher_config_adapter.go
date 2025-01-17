package dispatcher

import "time"

type IDispatcherConfigAdapter interface {
	GetMachineMonitorTm() time.Duration
	GetMachineExpiredTime() time.Duration

	RootTopicPub() string
	RootTopicSub() string
	DefaultQOS() byte
}
