package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Consume starts consuming messages from a queue until the context is canceled
func (r *rabbit) Consume(ctx context.Context, config ConfigConsume, f func(*amqp.Delivery)) (err error) {
	r.wgChannel.Add(1)
	defer r.wgChannel.Done()
	var msgs <-chan amqp.Delivery
	msgs, err = r.chConsumer.Consume(
		config.QueueName,
		config.Consumer,
		config.AutoAck,
		config.Exclusive,
		config.NoLocal,
		config.NoWait,
		config.Args,
	)
	if err != nil {
		return
	}
	var allCanceled bool
	for {
		select {
		case msg, ok := <-msgs:
			if !ok {
				return
			}
			r.wgChannel.Add(1)
			if config.ExecuteConcurrent {
				go func() {
					f(&msg)
					r.wgChannel.Done()
				}()
			} else {
				f(&msg)
				r.wgChannel.Done()
			}
		case <-ctx.Done():
			if allCanceled {
				continue
			}
			err = r.chConsumer.Cancel(config.Consumer, false)
			allCanceled = true
			continue
		}
	}
}
