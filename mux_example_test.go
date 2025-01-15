package msgmux_test

import (
	"fmt"

	"github.com/uudashr/msgmux"
)

func ExampleDispatchMux() {
	type CancelOrder struct {
		OrderID string
		Reason  string
	}

	type OrderCompleted struct {
		OrderID string
	}

	mux := msgmux.NewDispatchMux()
	mux.Handle(func(e CancelOrder) error {
		fmt.Printf("CancelOrder{OrderID:%s Reason:%s}\n", e.OrderID, e.Reason)
		return nil
	})

	mux.Handle(func(e OrderCompleted) error {
		fmt.Printf("OrderCompleted{OrderID:%s}\n", e.OrderID)
		return nil
	})

	mux.Dispatch(CancelOrder{
		OrderID: "order-123",
		Reason:  "Change my mind",
	})

	mux.Dispatch(OrderCompleted{
		OrderID: "order-123",
	})

	// Output:
	// CancelOrder{OrderID:order-123 Reason:Change my mind}
	// OrderCompleted{OrderID:order-123}
}
