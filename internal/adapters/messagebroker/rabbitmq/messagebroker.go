package rabbitmq

import (
	"context"
	"log"

	"github.com/mayank18049/GoChessMoveShare/internal/domain/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQMessageBroker struct {
	connection *amqp.Connection
	logger     *log.Logger
}

func NewRabbitMQMessageBroker(ctx context.Context, URI string, logger *log.Logger) *RabbitMQMessageBroker {
	logger.Printf("%s\n", URI)
	conn, err := amqp.Dial(URI)
	if err != nil {
		logger.Fatalf("Failed DialUp on RabbitMQ: %s ", err.Error())
	}
	rmq := RabbitMQMessageBroker{connection: conn, logger: logger}
	return &rmq
}

func (rmq *RabbitMQMessageBroker) CreateQueue(ctx context.Context, QueueName string, QueueType ports.QueueType) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("CreateQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()
	queueConfig := amqp.Table{}
	if QueueType == ports.QUEUE_REPLAY {
		queueConfig["x-queue-type"] = "stream"
	}
	_, err = ch.QueueDeclare(QueueName, true, false, false, false, queueConfig)
	if err != nil {
		rmq.logger.Printf("CreateQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("CreateQueue: Created Queue %s\n", QueueName)
	return ports.MESSAGE_SERVICE_OK
}

func (rmq *RabbitMQMessageBroker) CreateExchange(ctx context.Context, ExchangeName string, ExchangeType ports.ExchangeType) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("CreateExchange: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()
	var exchangekind string
	switch ExchangeType {
	case ports.EXCHANGE_FANOUT:
		exchangekind = amqp.ExchangeFanout
	case ports.EXCHANGE_DIRECT:
		exchangekind = amqp.ExchangeDirect
	default:
		rmq.logger.Printf("CreateExchange: Incorrect exchange type %d\n", ExchangeType)
		return ports.MESSAGE_SERVICE_FAILED
	}

	err = ch.ExchangeDeclare(ExchangeName, exchangekind, false, false, false, false, amqp.Table{})
	if err != nil {
		rmq.logger.Printf("CreateExchange: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("CreateExchange: Created Exchange %s\n", ExchangeName)
	return ports.MESSAGE_SERVICE_OK
}

func (rmq *RabbitMQMessageBroker) ConnectQueue(ctx context.Context, ExchangeName string, QueueName string, key string) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("ConnectQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()
	err = ch.QueueBind(QueueName, key, ExchangeName, false, amqp.Table{})
	if err != nil {
		rmq.logger.Printf("ConnectQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("ConnectQueue: Exchange: %s connected to Queue: %s\n", ExchangeName, QueueName)
	return ports.MESSAGE_SERVICE_OK
}
func (rmq *RabbitMQMessageBroker) DisconnectQueue(ctx context.Context, ExchangeName string, QueueName string, key string) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("DisconnectQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()
	err = ch.QueueUnbind(QueueName, key, ExchangeName, amqp.Table{})
	if err != nil {
		rmq.logger.Printf("DisconnectQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("DisconnectQueue: Exchange: %s Disconnected Queue: %s\n", ExchangeName, QueueName)
	return ports.MESSAGE_SERVICE_OK
}
func (rmq *RabbitMQMessageBroker) DeleteExchange(ctx context.Context, ExchangeName string) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("DeleteExchange: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()
	err = ch.ExchangeDelete(ExchangeName, true, false)
	if err != nil {
		rmq.logger.Printf("DeleteExchange: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("DeleteExchange: Deleted Exchange %s\n", ExchangeName)
	return ports.MESSAGE_SERVICE_OK
}

func (rmq *RabbitMQMessageBroker) DeleteQueue(ctx context.Context, QueueName string) ports.BrokerStatus {
	ch, err := rmq.connection.Channel()
	if err != nil {
		rmq.logger.Printf("CreateQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	defer ch.Close()

	_, err = ch.QueueDelete(QueueName, false, false, false)
	if err != nil {
		rmq.logger.Printf("DeleteQueue: %s\n", err.Error())
		return ports.MESSAGE_SERVICE_FAILED
	}
	rmq.logger.Printf("DeleteQueue: Deleted Queue %s\n", QueueName)
	return ports.MESSAGE_SERVICE_OK
}
