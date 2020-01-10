package mqlistener

import (
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"gitlab.alfa-bank.kz/crmw/library/ocrmerrors"
)

type Listener struct {
	ch           *amqp.Channel
	log          *logrus.Logger
	exchangeName string
	exchangeType string
	queueName    string
}

//Создает новго лиснера на данную queue по имени в данном exchange.
//Если необходимо будет отправлять сообщения в очередь по key,
//необходимо по отдельность вызвать функцию BindWithKey для кождого ключа
func NewListener(exchangeName, exchangeType, queueName string, channel *amqp.Channel, logger *logrus.Logger) (*Listener, error) {
	if queueName == "" {
		return nil, ocrmerrors.New(ocrmerrors.ARGISNIL, "queueueName must be", "queueueName объящательно должен быть заполнен")
	}
	if exchangeName != "" && exchangeType != "" {
		if err := channel.ExchangeDeclare(
			exchangeName, // name
			exchangeType, // type
			true,         // durable
			false,        // auto-deleted
			false,        // internal
			false,        // no-wait
			nil,          // arguments
		); err != nil {
			return nil, err
		}
	}
	queue, err := channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}
	return &Listener{
		ch:           channel,
		log:          logger,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
		queueName:    queue.Name,
	}, nil
}

func (ls *Listener) BindWithKey(key string) error {
	return ls.ch.QueueBind(
		ls.queueName,    // queue name
		key,             // routing key
		ls.exchangeName, // exchange
		false,
		nil)
}

func (ls *Listener) StartListen(handler MqHandler) error {
	ls.log.Tracef("Preparing listen queue: %s", ls.queueName)
	msgs, err := ls.ch.Consume(
		ls.queueName, // queue
		"",           // consumer
		false,        // auto ack
		false,        // exclusive
		false,        // no local
		false,        // no wait
		nil,          // args
	)
	if err != nil {
		//May be i must close msgs channel?, but now i dont know how do it
		ls.log.Errorf("Cant consume queue %s, err: %v", ls.queueName, err)
		return err
	}

	ls.log.Tracef("Start consume queue: %s", ls.queueName)
	for msg := range msgs {
		ls.log.Tracef("Start handle message: %s", string(msg.Body))
		go handler.HandleMqMessage(msg)
	}
	return nil

}
