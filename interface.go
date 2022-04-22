package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ combines all the interfaces of the package
type RabbitMQ interface {
	Connector
	Closer
	QueueCreator
	Consumer
}

// Connector is an interface for connecting to a RabbitMQ server
type Connector interface {
	Connect(config ConfigConnection) (notify chan *amqp.Error, err error)
}

// Closer is an interface for closing a RabbitMQ connection
type Closer interface {
	Close() (err error)
}

// QueueCreator is the interface for creating, binding and unbinding queues
type QueueCreator interface {
	CreateQueue(config ConfigQueue) (queue amqp.Queue, err error)
	BindQueueExchange(config ConfigBindQueue) (err error)
	UnbindQueueExchange(config ConfigBindQueue) (err error)
}

// Consumer is the interface for consuming messages from a queue
type Consumer interface {
	Consume(ctx context.Context, config ConfigConsume, f func(*amqp.Delivery)) (err error)
}

// Publisher is the interface for publishing messages to an exchange
type Publisher interface {
	Publish(ctx context.Context, obj interface{}, config ConfigPublish) (err error)
}
