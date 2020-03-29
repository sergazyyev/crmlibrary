package mqconsumer

import (
	"time"

	"github.com/sergazyyev/crmlibrary/ocrmerrors"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

type Session struct {
	connection      *amqp.Connection
	logger          *logrus.Logger
	done            chan bool
	notifyConnClose chan *amqp.Error
	isReady         bool
	listeners       map[string]*Listener
}

const (
	reconnectDelay      = 5 * time.Second
	reInitDelay         = 5 * time.Second
	resendDelay         = 5 * time.Second
	waitTime            = 1 * time.Second
	sleepWhenNotMessage = 10 * time.Second
)

func (s *Session) handleReconnect(addr string) {
	for {
		s.isReady = false
		s.logger.Debugf("Trying to connect %s", addr)

		_, err := s.connect(addr)

		if err != nil {
			s.logger.Errorf("Error when create connection to %s, err: %v", addr, err)
			s.logger.Infof("Retry to connect...")

			select {
			case <-s.done:
				return
			case <-time.After(reconnectDelay):
			}
			continue
		}

		s.logger.Debugf("Connection setup")
		if done := s.connIsClose(); done {
			break
		}
	}
}

func (s *Session) connIsClose() bool {
	for {
		select {
		case <-s.done:
			s.logger.Debugf("Take done signal, closing connection")
			return true
		case <-s.notifyConnClose:
			s.logger.Errorf("Connection closed. Reconnecting...")
			return false
		}
	}
}

//create new amqp connection for addr
func (s *Session) connect(addr string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(addr)

	if err != nil {
		return nil, err
	}

	s.changeConnection(conn)
	s.isReady = true

	return conn, nil
}

//Register new connection to object
//and new trigger channel for close connection
func (s *Session) changeConnection(conn *amqp.Connection) {
	s.connection = conn
	s.notifyConnClose = make(chan *amqp.Error)
	s.connection.NotifyClose(s.notifyConnClose)
}

func (s *Session) newListener(queueName, exchangeName, exchangeType string) (*Listener, error) {
	if !s.isReady {
		return nil, ocrmerrors.New(ocrmerrors.INVALID, "Session is unready", "Не создана сессия для сервера amqp")
	}

	listener := &Listener{
		session:      s,
		done:         make(chan bool),
		queueName:    queueName,
		exchangeName: exchangeName,
		exchangeType: exchangeType,
	}

	s.listeners[queueName] = listener

	return listener, nil
}

func (s *Session) Listener(queueName, exchangeName, exchangeType string) (*Listener, error) {
	if !s.isReady {
		return nil, ocrmerrors.New(ocrmerrors.INVALID, "Session is unready", "Не создана сессия для сервера amqp")
	}

	if queueName == "" {
		return nil, ocrmerrors.New(ocrmerrors.INVALID, "QueueName cant be empty", "Имя очереди не может быть пустым")
	}

	listener, ok := s.listeners[queueName]
	if !ok {
		return s.newListener(queueName, exchangeName, exchangeType)
	} else {
		return listener, nil
	}
}

func New(addr string, loggerLvl string) (*Session, error) {
	conn, err := amqp.Dial(addr)
	if err != nil {
		return nil, err
	}

	err = conn.Close()
	if err != nil {
		return nil, err
	}

	lvl, err := logrus.ParseLevel(loggerLvl)
	if err != nil {
		return nil, err
	}

	logger := logrus.New()
	logger.SetLevel(lvl)

	s := &Session{
		logger:    logger,
		listeners: make(map[string]*Listener),
		done:      make(chan bool),
	}

	go s.handleReconnect(addr)

	time.Sleep(waitTime)

	return s, nil
}

func (s *Session) Close() error {
	if !s.isReady {
		return ocrmerrors.New(ocrmerrors.INVALID, "Connection is not created for close", "")
	}

	for _, ls := range s.listeners {
		err := ls.Close()
		if err != nil {
			return err
		}
	}

	close(s.done)

	err := s.connection.Close()
	if err != nil {
		return err
	}

	s.isReady = false
	return nil
}
