package event

import "time"

type OrdersListed struct {
	Name     string
	DateTime time.Time
	Payload  interface{}
}

func NewOrdersListed() *OrdersListed {
	return &OrdersListed{
		Name:     "OrdersListed",
		DateTime: time.Now(),
	}
}

func (e *OrdersListed) GetName() string {
	return e.Name
}

func (e *OrdersListed) GetDateTime() time.Time {
	return e.DateTime
}

func (e *OrdersListed) GetPayload() interface{} {
	return e.Payload
}

func (e *OrdersListed) SetPayload(payload interface{}) {
	e.Payload = payload
}
