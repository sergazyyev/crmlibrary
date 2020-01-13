package mq

import (
	"github.com/sergazyyev/crmlibrary/ocrmconfigs"
	"github.com/streadway/amqp"
)

func GetConnection(config *ocrmconfigs.RabbitMqConfig) (*amqp.Connection, error) {
	conn, err := amqp.Dial(config.GetConnectionUrl())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func GetChannel(connection *amqp.Connection) (*amqp.Channel, error) {
	ch, err := connection.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}