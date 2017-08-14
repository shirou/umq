package umq

import "context"

type TransportType string

const (
	TransportAny    TransportType = "any"
	TransportMemory               = "memory"
	TransportSQS                  = "sqs"
	TransportRedis                = "redis"
)

type Transport interface {
	Connect(string) error
	GetQueue(string, ...Option) (Queue, error)
}

type Option interface {
	Target() string
	Apply(Queue) error
}

type Queue interface {
	Close() error

	Receive(...Option) (Message, error)
	ReceiveWithContext(context.Context, ...Option) (Message, error)

	Delete(Message, ...Option) error
	DeleteWithContext(context.Context, Message, ...Option) error

	Send(Message, ...Option) error
	SendWithContext(context.Context, Message, ...Option) error
}

func optionAvailable(target, opt string) bool {
	if opt == "any" || target == opt {
		return true
	}
	return false
}
