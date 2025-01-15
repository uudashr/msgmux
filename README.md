# msgmux
Message multiplexer.

It can be use to process any kind of message, such as event, command, etc.

## Usage

```go
// Message definitions
type CancelOrder struct {
    OrderID string
    Reason  string
}

type OrderCompleted struct {
    OrderID string
}

// Create the multiplexer
mux := msgmux.NewDispatchMux()

// Define handlers
mux.Handle(func(e CancelOrder) error {
    fmt.Printf("CancelOrder{OrderID:%s Reason:%s}\n", e.OrderID, e.Reason)
    return nil
})

mux.Handle(func(e OrderCompleted) error {
    fmt.Printf("OrderCompleted{OrderID:%s}\n", e.OrderID)
    return nil
})

// Dispatch messages
mux.Dispatch(CancelOrder{
    OrderID: "order-123",
    Reason:  "Change my mind",
})

mux.Dispatch(OrderCompleted{
    OrderID: "order-123",
})
```
