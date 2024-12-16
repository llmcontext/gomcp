package eventbus

type EventBusTopic string

const (
	EventBusTopicProxyClientConnected    EventBusTopic = "proxy_client_connected"
	EventBusTopicProxyClientDisconnected EventBusTopic = "proxy_client_disconnected"
)

type EventBusMessage struct {
	Topic   EventBusTopic
	Payload interface{}
}
