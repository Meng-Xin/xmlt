package utils

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Client     *amqp.Connection
	Ch         *amqp.Channel
	QueueName  string
	Exchange   string
	RoutingKey string
}

func NewRabbitMQ(queueName, exchange, routingKey, dsn string) *RabbitMQ {
	rabbitmq := RabbitMQ{
		QueueName:  queueName,
		Exchange:   exchange,
		RoutingKey: routingKey,
	}
	var err error
	//创建rabbitmq连接
	rabbitmq.Client, err = amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}

	//创建Channel
	rabbitmq.Ch, err = rabbitmq.Client.Channel()
	if err != nil {
		panic(err.Error())
	}

	// 创建队列
	_, err = rabbitmq.Ch.QueueDeclare(queueName, false, false, false, true, nil)
	if err != nil {
		panic(err.Error())
	}
	return &rabbitmq
}

func (r *RabbitMQ) PublishOnQueue(ctx context.Context, msg []byte) error {
	return r.Ch.PublishWithContext(ctx, r.Exchange, r.QueueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
}

func (r *RabbitMQ) SubscribeToQueue(consumerName string) ([]byte, error) {
	msgs, err := r.Ch.Consume(r.QueueName, consumerName, false, false, false, false, nil)
	if err != nil {
		return nil, err
	}
	select {
	case res := <-msgs:
		return res.Body, nil
	}
}

type User struct {
	ID   int
	Name string
}
