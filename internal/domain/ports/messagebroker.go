package ports

import "context"

type BrokerStatus int
type QueueType int
type ExchangeType int

const (
	MESSAGE_SERVICE_OK BrokerStatus = -iota
	MESSAGE_SERVICE_FAILED
	MESSAGE_SERVICE_NOT_IMPLEMENTED
)
const (
	QUEUE_REPLAY QueueType = iota + 1
	QUEUE_NO_REPLAY
)
const (
	EXCHANGE_FANOUT ExchangeType = iota + 1
	EXCHANGE_DIRECT
)

type MessageBroker interface {
	CreateQueue(ctx context.Context, QueueName string, Queuetype QueueType) BrokerStatus
	CreateExchange(ctx context.Context, ExchangeName string, Exchangetype ExchangeType) BrokerStatus
	ConnectQueue(ctx context.Context, ExchangeName string, QueueName string, TopicKey string) BrokerStatus
	DisconnectQueue(ctx context.Context, ExchangeName string, QueueName string, key string) BrokerStatus
	DeleteQueue(ctx context.Context, QueueName string) BrokerStatus
	DeleteExchange(ctx context.Context, ExchangeName string) BrokerStatus
}
