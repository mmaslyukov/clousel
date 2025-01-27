package operator

type IPortOperatorAdapterConfig interface {
	RootTopicPub() string
	RootTopicSub() string
	DefaultQOS() byte
}
