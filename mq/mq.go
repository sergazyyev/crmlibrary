package mq

import (
	"github.com/streadway/amqp"
	"gitlab.alfa-bank.kz/crmw/library/ocrmconfigs"
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