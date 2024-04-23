package subscriber

import (
	"fmt"

	"github.com/nats-io/stan.go"
)

func NewSubscribe(sc stan.Conn) (stan.Subscription, error) {
	sub, err := sc.Subscribe("example-subject", func(msg *stan.Msg) {
		fmt.Printf("Получено сообщение: %s\n", string(msg.Data))
	}, stan.DurableName("example-durable"))
	if err != nil {
		return nil, fmt.Errorf("Ошибка создания подписчика: %v", err)
	}
	return sub, nil
}
