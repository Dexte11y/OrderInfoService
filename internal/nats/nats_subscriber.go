package nats

import (
	"encoding/json"
	"fmt"
	"orderinfoservice/internal/model"
	"strings"

	"github.com/nats-io/stan.go"
)

func NewSubscribe(sc stan.Conn) (stan.Subscription, error) {
	sub, err := sc.Subscribe("example-subject", func(msg *stan.Msg) {
		fmt.Println("Subscriber: получено сообщение")
	}, stan.DurableName("example-durable"))
	if err != nil {
		return nil, fmt.Errorf("ошибка создания подписчика: %v", err)
	}
	return sub, nil
}

func GetMsgByIDSubscribe(sc stan.Conn, orderUID string, msgChan chan<- model.Order) error {
	sub, err := sc.Subscribe("example-subject", func(msg *stan.Msg) {
		keyEndIndex := strings.Index(string(msg.Data), "{")
		if keyEndIndex == -1 {
			fmt.Println("Неверный формат сообщения")
			return
		}
		key := string(msg.Data[:keyEndIndex])

		if key == orderUID {
			var order model.Order
			if err := json.Unmarshal(msg.Data[keyEndIndex:], &order); err != nil {
				fmt.Println("Ошибка при распаковке сообщения:", err)
				return
			}
			// Отправляем сообщение в канал
			msgChan <- order
		}
	}, stan.DurableName("example-durable1"))
	if err != nil {
		return fmt.Errorf("ошибка создания подписчика: %v", err)
	}
	// Закрываем подписку после использования
	defer sub.Close()
	return nil
}
