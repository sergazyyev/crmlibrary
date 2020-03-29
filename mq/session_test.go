package mq_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/sergazyyev/crmlibrary/mq/mqproducer"

	"github.com/sergazyyev/crmlibrary/mq/mqconsumer"

	"github.com/streadway/amqp"

	"github.com/stretchr/testify/assert"
)

var (
	addr   = "amqp://mqadmin:mqadminpassword@localhost:5672"
	logLvl = "debug"
)

type Handler struct {
}

func (h Handler) HandleMessage(delivery amqp.Delivery) {
	defer func() {
		err := delivery.Ack(false)
		if err != nil {
			fmt.Println(fmt.Sprintf("Error when ask message: %s, err: %v", string(delivery.Body), err))
		} else {
			fmt.Println(fmt.Sprintf("Asked: %s", string(delivery.Body)))
		}
	}()
	fmt.Printf("Receive %s\n", string(delivery.Body))
}

func TestConcurrentProduceAndConsume(t *testing.T) {
	sess, err := mqconsumer.New(addr, logLvl)
	assert.NotNil(t, sess)
	assert.NoError(t, err)
	go func() {
		ls, err := sess.Listener("test_queue", "", "")
		assert.NoError(t, err)
		assert.NotNil(t, ls)
		var hr Handler
		err = ls.StartListen(hr, false)
		if err != nil {
			t.Fatal(err)
		}
	}()

	go func() {
		ls, err := sess.Listener("test_queue1", "", "")
		assert.NoError(t, err)
		assert.NotNil(t, ls)
		var hr Handler
		err = ls.StartListen(hr, false)
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(1 * time.Second)

	go func() {
		sender, err := mqproducer.New(addr, logLvl)
		if err != nil {
			t.Fatal(err)
		}

		for i := 0; i <= 10; i++ {
			if err := sender.UnsafeSendToQueue("test_queue", amqp.Publishing{
				Body: []byte(fmt.Sprintf("test_queue message: %d", i)),
			}); err != nil {
				fmt.Printf("Push failed: %s\n", err)
				t.Fatal(err)
			}
		}

		for i := 0; i <= 7; i++ {
			if err := sender.UnsafeSendToQueue("test_queue1", amqp.Publishing{
				Body: []byte(fmt.Sprintf("test_queue1 message: %d", i)),
			}); err != nil {
				fmt.Printf("Push failed: %s\n", err)
				t.Fatal(err)
			}
		}

		err = sender.Close()
		if err != nil {
			t.Fatal(err)
		}
	}()

	time.Sleep(60 * time.Second)
	err = sess.Close()
	assert.NoError(t, err)
}
