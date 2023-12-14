package rabbitmq

import (
	"context"
	"github.com/Verce11o/yata-auth/config"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"time"
)

type EmailPublisher struct {
	AmqpConn *amqp.Connection
	log      *zap.SugaredLogger
	trace    trace.Tracer
	cfg      config.RabbitMQ
}

func NewEmailPublisher(amqpConn *amqp.Connection, log *zap.SugaredLogger, trace trace.Tracer, cfg config.RabbitMQ) *EmailPublisher {
	return &EmailPublisher{AmqpConn: amqpConn, log: log, trace: trace, cfg: cfg}
}

func (c *EmailPublisher) createChannel(exchangeName string, queueName string, bindingKey string) *amqp.Channel {

	ch, err := c.AmqpConn.Channel()

	if err != nil {
		panic(err)
	}

	err = ch.ExchangeDeclare(
		exchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	queue, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	err = ch.QueueBind(
		queue.Name,
		bindingKey,
		exchangeName,
		false,
		nil,
	)

	if err != nil {
		panic(err)
	}

	return ch

}

func (c *EmailPublisher) Publish(ctx context.Context, message []byte) error {

	ch := c.createChannel(c.cfg.ExchangeName, c.cfg.QueueName, c.cfg.BindingKey)

	if err := ch.PublishWithContext(
		ctx,
		c.cfg.ExchangeName,
		c.cfg.BindingKey,
		false,
		false,
		amqp.Publishing{
			ContentType:  "text/plain",
			DeliveryMode: amqp.Persistent,
			MessageId:    uuid.New().String(),
			Timestamp:    time.Now(),
			Body:         message,
		}); err != nil {

		return err
	}

	return nil

}
