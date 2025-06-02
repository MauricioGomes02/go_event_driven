package ports

type IProductEventBusPort interface {
	Publish(eventType string, message []byte) error
}
