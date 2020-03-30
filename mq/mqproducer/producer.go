package mqproducer

import (
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sergazyyev/crmlibrary/ocrmerrors"

	"github.com/streadway/amqp"
)

type Sender struct {
	connection      *amqp.Connection
	channel         *amqp.Channel
	logger          *logrus.Logger
	isConfirmMode   bool
	notifyChanClose chan *amqp.Error
	notifyConnClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isReady         bool
	done            chan bool
}

const (
	reconnectDelay = 5 * time.Second
	reInitDelay    = 2 * time.Second
	resendDelay    = 5 * time.Second
	waitTime       = 1 * time.Second
)

func (sr *Sender) handleReconnect(addr string) {
	for {
		sr.isReady = false

		sr.logger.Debugf("Creating connection to %s", addr)
		conn, err := sr.connect(addr)

		if err != nil {
			sr.logger.Errorf("Error when create connection to %s, err: %v", addr, err)
			sr.logger.Infof("Retry to connect...")

			select {
			case <-sr.done:
				sr.logger.Debugf("Take done signal for producer")
				return
			case <-time.After(reInitDelay):
			}
			continue
		}

		if done := sr.handleReInit(conn); done {
			break
		}
	}
}

func (sr *Sender) handleReInit(conn *amqp.Connection) bool {
	for {
		sr.isReady = false

		err := sr.init(conn)
		if err != nil {
			sr.logger.Errorf("Error when init channel for producer, err: %v", err)
			sr.logger.Infof("ReInit channel...")

			select {
			case <-sr.done:
				return true
			case <-time.After(reInitDelay):
			}
			continue
		}

		select {
		case <-sr.done:
			sr.logger.Debugf("Take done signal from producer")
			return true
		case <-sr.notifyConnClose:
			sr.logger.Errorf("Take signal that connection is closed, reconnect ...")
			return false
		case <-sr.notifyChanClose:
			sr.logger.Errorf("Channel exception, reInit channel ...")
			return false
		}
	}
}

func (sr *Sender) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	sr.changeConnection(conn)
	sr.logger.Debugf("Connection to addr %s setup", addr)
	return conn, nil
}

func (sr *Sender) init(conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		sr.logger.Debugf("Error when declare channel, err: %v", err)
		return err
	}

	if sr.isConfirmMode {
		err = ch.Confirm(false)
		if err != nil {
			sr.logger.Debugf("Error when set confirm mod to channel, err: %v", err)
			return err
		}
	}

	sr.changeChannel(ch)
	sr.isReady = true

	return nil
}

func (sr *Sender) changeChannel(channel *amqp.Channel) {
	sr.channel = channel
	sr.notifyChanClose = make(chan *amqp.Error)
	sr.channel.NotifyClose(sr.notifyChanClose)
	if sr.isConfirmMode {
		sr.notifyConfirm = make(chan amqp.Confirmation, 1)
		sr.channel.NotifyPublish(sr.notifyConfirm)
	}
}

func (sr *Sender) changeConnection(connection *amqp.Connection) {
	sr.connection = connection
	sr.notifyConnClose = make(chan *amqp.Error)
	sr.connection.NotifyClose(sr.notifyConnClose)
}

func (sr *Sender) Close() error {
	if !sr.isReady {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}

	close(sr.done)

	err := sr.channel.Close()
	if err != nil {
		return err
	}

	err = sr.connection.Close()
	if err != nil {
		return err
	}

	sr.isReady = false
	return nil
}

func New(addr string, logLvl string, confirmMode bool) (*Sender, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	err = conn.Close()
	if err != nil {
		return nil, err
	}

	lvl, err := logrus.ParseLevel(logLvl)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(lvl)

	sr := &Sender{
		logger:        logger,
		done:          make(chan bool),
		isConfirmMode: confirmMode,
	}

	go sr.handleReconnect(addr)
	time.Sleep(waitTime)

	return sr, nil
}

func (sr *Sender) sendToQueue(queueName string, data amqp.Publishing) error {
	sr.logger.Tracef("Trying send message: %s", string(data.Body))
	sr.logger.Debugf("Trying send message")
	return sr.channel.Publish("", queueName, false, false, data)
}

func (sr *Sender) SendToQueue(queueName string, data amqp.Publishing) error {
	if !sr.isReady {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}
	if !sr.isConfirmMode {
		return sr.sendToQueue(queueName, data)
	} else {
		for {
			err := sr.sendToQueue(queueName, data)
			if err != nil {
				sr.logger.Errorf("Push failed. Retrying...")
				select {
				case <-sr.done:
					return ocrmerrors.New(ocrmerrors.INTERNAL, "Shutdown command takes", "Получена команда остановки")
				case <-time.After(resendDelay):
				}
				continue
			}
			select {
			case confirm := <-sr.notifyConfirm:
				if confirm.Ack {
					sr.logger.Debugf("Push message confirmed!")
					return nil
				}
			case <-time.After(resendDelay):
			}
			sr.logger.Debugf("Push didn't confirm. Retrying...")
		}
	}
}

func (sr *Sender) sendToExchange(exchangeName string, data amqp.Publishing) error {
	sr.logger.Tracef("Trying send message: %s", string(data.Body))
	sr.logger.Debugf("Trying send message")
	return sr.channel.Publish(exchangeName, "", false, false, data)
}

func (sr *Sender) SafeSendToExchange(exchangeName string, data amqp.Publishing) error {
	if !sr.isReady {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}
	if !sr.isConfirmMode {
		return sr.sendToExchange(exchangeName, data)
	} else {
		for {
			err := sr.sendToExchange(exchangeName, data)
			if err != nil {
				sr.logger.Errorf("Push failed. Retrying...")
				select {
				case <-sr.done:
					return ocrmerrors.New(ocrmerrors.INTERNAL, "Shutdown command takes", "Получена команда остановки")
				case <-time.After(resendDelay):
				}
				continue
			}
			select {
			case confirm := <-sr.notifyConfirm:
				if confirm.Ack {
					sr.logger.Debugf("Push message confirmed!")
					return nil
				}
			case <-time.After(resendDelay):
			}
			sr.logger.Debugf("Push didn't confirm. Retrying...")
		}
	}
}

func (sr *Sender) sendToExchangeWithKey(exchangeName, key string, data amqp.Publishing) error {
	sr.logger.Tracef("Trying send message: %s", string(data.Body))
	sr.logger.Debugf("Trying send message")
	return sr.channel.Publish(exchangeName, key, false, false, data)
}

func (sr *Sender) SendToExchangeWithKey(exchangeName, key string, data amqp.Publishing) error {
	if !sr.isReady {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}
	if !sr.isConfirmMode {
		return sr.sendToExchangeWithKey(exchangeName, key, data)
	} else {
		for {
			err := sr.sendToExchangeWithKey(exchangeName, key, data)
			if err != nil {
				sr.logger.Errorf("Push failed. Retrying...")
				select {
				case <-sr.done:
					return ocrmerrors.New(ocrmerrors.INTERNAL, "Shutdown command takes", "Получена команда остановки")
				case <-time.After(resendDelay):
				}
				continue
			}
			select {
			case confirm := <-sr.notifyConfirm:
				if confirm.Ack {
					sr.logger.Debugf("Push message confirmed!")
					return nil
				}
			case <-time.After(resendDelay):
			}
			sr.logger.Debugf("Push didn't confirm. Retrying...")
		}
	}
}
