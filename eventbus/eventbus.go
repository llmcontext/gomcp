package eventbus

import (
	"sync"
)

type EventBus struct {
	subscribers map[EventBusTopic][]chan EventBusMessage
	mu          sync.Mutex
}

func NewEventBus() *EventBus {
	return &EventBus{
		subscribers: make(map[EventBusTopic][]chan EventBusMessage),
	}
}

func (bus *EventBus) Subscribe(topic EventBusTopic, ch chan EventBusMessage) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	bus.subscribers[topic] = append(bus.subscribers[topic], ch)
}

func (bus *EventBus) Publish(message EventBusMessage) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	for _, ch := range bus.subscribers[message.Topic] {
		ch <- message
	}
}

// func main() {
// 	bus := NewMessageBus()

// 	// Subscriber 1
// 	go func() {
// 		ch := make(chan Message)
// 		bus.Subscribe("topic1", ch)

// 		for msg := range ch {
// 			fmt.Println("Subscriber 1 received:", msg)
// 		}
// 	}()

// 	// Subscriber 2
// 	go func() {
// 		ch := make(chan Message)
// 		bus.Subscribe("topic2", ch)

// 		for msg := range ch {
// 			fmt.Println("Subscriber 2 received:", msg)
// 		}
// 	}()

// 	// Publisher
// 	go func() {
// 		bus.Publish("topic1", "Hello from topic 1")
// 		bus.Publish("topic2", "Hello from topic 2")

// 		time.Sleep(time.Second)
// 	}()

// 	time.Sleep(2 * time.Second)
// }
