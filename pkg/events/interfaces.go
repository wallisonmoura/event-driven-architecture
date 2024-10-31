package events

import (
	"sync"
	"time"
)

type EventInteface interface {
	GetName() string
	GetDateTime() time.Time
	getPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	Handle(event EventInteface, wg *sync.WaitGroup)
}

type EventDispatcherInterface interface {
	Register(eventName string, handler EventHandlerInterface) error
	Dispatch(event EventInteface) error
	Remove(eventName string, handler EventHandlerInterface) error
	Has(eventName string, handler EventHandlerInterface) bool
	Clear()
}
