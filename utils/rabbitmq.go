package utils

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	client     *amqp.Connection
	ch         *amqp.Channel
	queueName  string
	exchange   string
	routingKey string
}

func NewRabbitMQ(queueName, exchange, routingKey, dsn string) *RabbitMQ {
	rabbitmq := RabbitMQ{
		queueName:  queueName,
		exchange:   exchange,
		routingKey: routingKey,
	}
	var err error
	//创建rabbitmq连接
	rabbitmq.client, err = amqp.Dial(dsn)
	if err != nil {
		panic(err.Error())
	}

	//创建Channel
	rabbitmq.ch, err = rabbitmq.client.Channel()
	if err != nil {
		panic(err.Error())
	}

	// 创建队列
	_, err = rabbitmq.ch.QueueDeclare(queueName, false, false, false, true, nil)
	if err != nil {
		panic(err.Error())
	}
	return &rabbitmq
}

func (r *RabbitMQ) PublishOnQueue(ctx context.Context, msg []byte) error {
	return r.ch.PublishWithContext(ctx, r.exchange, r.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
}

func (r *RabbitMQ) SubscribeToQueue(consumerName string) ([]byte, error) {
	msgs, err := r.ch.Consume(r.queueName, consumerName, false, false, false, false, nil)
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
