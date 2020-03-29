package mqconsumer

import (
	"context"
	"time"

	"github.com/sergazyyev/crmlibrary/ocrmerrors"
	"github.com/streadway/amqp"
)

type Handler interface {
	HandleMessage(delivery amqp.Delivery)
}

type Listener struct {
	session         *Session
	queueName       string
	exchangeName    string
	exchangeType    string
	isReady         bool
	done            chan bool
	channel         *amqp.Channel
	notifyChanClose chan *amqp.Error
	notifyConnClose chan *amqp.Error
	stream          <-chan amqp.Delivery
	handler         Handler
	handleAsc       bool
}

func (ls *Listener) init(conn *amqp.Connection) error {
	if !ls.session.isReady {
		return ocrmerrors.New(ocrmerrors.INTERNAL, "Session is unready", "Сессия еще не готова")
	}

	ch, err := conn.Channel()
	if err != nil {
		ls.session.logger.Debugf("Error when declare channel, err: %v", err)
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		ls.session.logger.Debugf("Error when set confirm mod to channel, err: %v", err)
		return err
	}

	ls.changeChannel(ch)

	if err := ls.declare(); err != nil {
		return err
	}

	if err := ls.changeStream(); err != nil {
		return err
	}

	return nil
}

func (ls *Listener) changeChannel(channel *amqp.Channel) {
	ls.channel = channel
	ls.notifyChanClose = make(chan *amqp.Error)
	ls.notifyConnClose = make(chan *amqp.Error)
	ls.channel.NotifyClose(ls.notifyChanClose)
	ls.session.connection.NotifyClose(ls.notifyConnClose)
}

func (ls *Listener) StartListen(handler Handler, async bool) error {
	if handler == nil {
		return ocrmerrors.New(ocrmerrors.ARGISNIL, "Handle interface must not be empty", "")
	}

	ls.handler = handler
	ls.handleAsc = async

	for {
		ls.isReady = false

		ls.session.logger.Debugf("Starting initialize listener %s", ls.queueName)
		err := ls.init(ls.session.connection)

		if err != nil {
			ls.session.logger.Errorf("Error when create channel, err: %v", err)
			ls.session.logger.Infof("Retry to connect...")

			select {
			case <-ls.session.done:
				ls.session.logger.Debugf("Take done signal from session")
				return nil
			case <-ls.done:
				ls.session.logger.Debugf("Take done signal from listener")
				return nil
			case <-time.After(reInitDelay):
			}
			continue
		}

		ctx, cancel := context.WithCancel(context.Background())
		go ls.consume(ctx)

		ls.session.logger.Debugf("Initialized listener %s", ls.queueName)
		if done := ls.chanIsExcept(); done {
			cancel()
			break
		}

		cancel()
	}

	return nil
}

func (ls *Listener) chanIsExcept() bool {
	for {
		select {
		case <-ls.done:
			ls.session.logger.Debugf("Take done signal from listener")
			return true
		case <-ls.notifyConnClose:
			ls.session.logger.Debugf("Take signal from session that connection is closed, reInit channel ...")
			time.Sleep(reconnectDelay)
			return false
		case <-ls.notifyChanClose:
			ls.session.logger.Infof("Channel exception, reInit channel ...")
			return false
		}
	}
}

func (ls *Listener) changeStream() error {
	st, err := ls.getStream()
	if err != nil {
		ls.session.logger.Debugf("Error when get stream for listener: %s, err: %v", ls.queueName, err)
		return err
	}
	ls.stream = st
	ls.session.logger.Debugf("Stream changed")
	return nil
}

func (ls *Listener) ready() bool {
	return ls.session != nil && ls.session.isReady && ls.isReady
}

func (ls *Listener) declare() error {
	if ls.queueName == "" {
		return ocrmerrors.New(ocrmerrors.ARGISNIL, "queueName must be", "queueName объящательно должен быть заполнен")
	}

	if ls.exchangeName != "" && ls.exchangeType != "" {
		ls.session.logger.Tracef("Declare exchange, name: %s, type: %s", ls.exchangeName, ls.exchangeType)
		if err := ls.channel.ExchangeDeclare(
			ls.exchangeName,
			ls.exchangeType,
			true,
			false,
			false,
			false,
			nil,
		); err != nil {
			ls.session.logger.Debugf("Error when declare exchange, err: %v", err)
			return err
		}
	}

	ls.session.logger.Tracef("Declare queue, name: %s", ls.queueName)
	queue, err := ls.channel.QueueDeclare(
		ls.queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ls.session.logger.Debugf("Error when declare queue, err: %v", err)
		return err
	}
	ls.queueName = queue.Name

	ls.isReady = true

	return nil
}

func (ls *Listener) BindWithKey(key string) error {
	return ls.channel.QueueBind(
		ls.queueName,
		key,
		ls.exchangeName,
		false,
		nil,
	)
}

func (ls *Listener) consume(ctx context.Context) {
	for msg := range ls.stream {
		if ls.handleAsc {
			ls.session.logger.Tracef("Start handle in new goroutine message: %s", string(msg.Body))
			go ls.handler.HandleMessage(msg)
		} else {
			ls.session.logger.Tracef("Start handle message: %s", string(msg.Body))
			ls.handler.HandleMessage(msg)
			ls.session.logger.Tracef("End handle message: %s", string(msg.Body))
		}
	}
}

func (ls *Listener) getStream() (<-chan amqp.Delivery, error) {
	if !ls.ready() {
		return nil, ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}

	st, err := ls.channel.Consume(
		ls.queueName,
		ls.queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	return st, nil
}

func (ls *Listener) Close() error {
	if !ls.ready() {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection or channel with server not initialized", "Соеденение с сервером не установлено")
	}

	ls.session.logger.Trace("Cancel consumer")
	err := ls.channel.Cancel(ls.queueName, false)
	if err != nil {
		ls.session.logger.Errorf("Error when cancel consumer, err: %s", err)
		return err
	}

	ls.session.logger.Trace("Stopping reInitialize process")
	close(ls.done)

	ls.session.logger.Trace("Close channel")
	err = ls.channel.Close()
	if err != nil {
		ls.session.logger.Errorf("Error when close channel, err: %s", err)
		return err
	}

	ls.isReady = false

	return nil
}
