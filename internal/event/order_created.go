package event

import "time"

type OrderCreated struct {
	Name     string
	DateTime time.Time
	Payload  interface{}
}

func NewOrderCreated() *OrderCreated {
	return &OrderCreated{
		Name:     "OrderCreated",
		DateTime: time.Now(),
	}
}

func (e *OrderCreated) GetName() string {
	return e.Name
}

func (e *OrderCreated) GetDateTime() time.Time {
	return e.DateTime
}

func (e *OrderCreated) GetPayload() interface{} {
	return e.Payload
}

func (e *OrderCreated) SetPayload(payload interface{}) {
	e.Payload = payload
}
