package ports

import "context"

type IProductEventBusPort interface {
	Publish(_context context.Context, eventType string, message []byte) error
}
