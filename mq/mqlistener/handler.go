package mqlistener

import "github.com/streadway/amqp"

type MqHandler interface {
	HandleMqMessage(delivery amqp.Delivery)
}
