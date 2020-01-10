package mqsender

import (
	"github.com/streadway/amqp"
)

type Sender struct {
	connStr string
}

func NewSender(connectionString string) (*Sender, error) {
	conn, err := amqp.Dial(connectionString)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return &Sender{connStr: connectionString}, nil
}

func (sr *Sender) SendToQueue(queueName string, data amqp.Publishing) error {
	conn, err := sr.getConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := sr.getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish("", queueName, false, false, data)
}

func (sr *Sender) SendToExchange(exchangeName string, data amqp.Publishing) error {
	conn, err := sr.getConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := sr.getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(exchangeName, "", false, false, data)
}

func (sr *Sender) SendToExchangeWithKey(exchangeName, key string, data amqp.Publishing) error {
	conn, err := sr.getConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := sr.getChannel(conn)
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.Publish(exchangeName, key, false, false, data)
}

func (sr *Sender) getConnection() (*amqp.Connection, error) {
	return amqp.Dial(sr.connStr)
}

func (sr *Sender) getChannel(conn *amqp.Connection) (*amqp.Channel, error) {
	return conn.Channel()
}
