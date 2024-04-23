package nats

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

// ConnectAndSubscribe подключается к серверу NATS и подписывается на тему
func ConnectToNats() (stan.Conn, error) {
	sc, err := stan.Connect("test-cluster", "example-client")
	if err != nil {
		return nil, fmt.Errorf("ошибка при подключении к серверу NATS: %w", err)
	}

	return sc, nil
}
