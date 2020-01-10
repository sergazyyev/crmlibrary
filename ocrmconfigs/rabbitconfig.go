package ocrmconfigs

import (
	"fmt"
	"strings"
)

type RabbitMqConfig struct {
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	Username    string `toml:"username"`
	Password    string `toml:"password"`
	VirtualHost string `toml:"virtual_host"`
}

type QueueConfig struct {
	QueueName   string `toml:"queue_name"`
	BindKeysStr string `toml:"bind_keys"`
}

func (qc *QueueConfig) GetBindKeys() []string {
	return strings.Split(qc.BindKeysStr, ",")
}

func (mq *RabbitMqConfig) GetConnectionUrl() string {
	if mq.VirtualHost != "" {
		return fmt.Sprintf("amqp://%s:%s@%s:%d/%s", mq.Username, mq.Password, mq.Host, mq.Port, mq.VirtualHost)
	} else {
		return fmt.Sprintf("amqp://%s:%s@%s:%d", mq.Username, mq.Password, mq.Host, mq.Port)
	}
}
