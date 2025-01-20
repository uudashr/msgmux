package msgmux

import (
	"context"
	"fmt"
	"reflect"
)

// Message represents a specific message. It must be a struct type.
//
// Example:
//
//	type CancelOrder struct {
//		OrderID string
//		Reason string
//	}
//
//	type OrderCompleted struct {
//		OrderID string
//	}
type Message any

// MssageHandler is a function type to handle the messages.
//
// The form of function is:
//
//	func(msg Message) error
//
// or with [context.Context] parameter:
//
//	func(ctx context.Context, msg Message) error
//
// where the [Message] represents a specific message type (struct).
//
// Example:
//
//	mux := msgmux.NewDispatchMux()
//	mux.Handle(func(event OrderCompleted) error {
//		// handle the event
//		return nil
//	})
type MessageHandler any

// DispatchMux is a message multiplexer.
type DispatchMux struct {
	handlers map[reflect.Type]MessageHandler
}

// NewDispatchMux allocates and returns a new [DispatchMux].
func NewDispatchMux() *DispatchMux {
	return &DispatchMux{}
}

// Handle registers a fn as handler for the specific message type.
//
// The fn need to be valid [MessageHandler], otherwise it will panic.
// Multiple registrations for the same message type is not allowed and will cause a panic.
func (m *DispatchMux) Handle(fn MessageHandler) {
	if err := validateHandler(fn); err != nil {
		panic(err)
	}

	if m.handlers == nil {
		m.handlers = make(map[reflect.Type]MessageHandler)
	}

	fnType := reflect.TypeOf(fn)
	var msgType reflect.Type
	switch fnType.NumIn() {
	case 1:
		msgType = fnType.In(0)
	case 2:
		msgType = fnType.In(1)
	default:
		panic(fmt.Sprintf("msgmux: invalid handler function signature (got: %v)", fnType))
	}

	if _, reg := m.handlers[msgType]; reg {
		panic(fmt.Sprintf("msgmux: handler for message %v already registered", msgType.Name()))
	}

	m.handlers[msgType] = fn
}

// DispatchContext dispatches the message.
//
// The msg need to be a valid [Message], otherwise it will return an error.
func (m *DispatchMux) DispatchContext(ctx context.Context, msg Message) error {
	eventType := reflect.TypeOf(msg)
	if eventType.Kind() != reflect.Struct {
		return fmt.Errorf("msgmux: msg should be a struct (got: %v)", eventType.Kind())
	}

	if m.handlers == nil {
		return nil
	}

	handler, reg := m.handlers[eventType]
	if !reg {
		return fmt.Errorf("msgmux: no handler registered for message %v", eventType.Name())
	}

	return invokeHandler(ctx, handler, msg)
}

// Dispatch dispatches the message.
//
// A shorthand of [DispatchMux.DispatchContext] with [context.Background].
func (m *DispatchMux) Dispatch(msg Message) error {
	return m.DispatchContext(context.Background(), msg)
}

func invokeHandler(ctx context.Context, handler MessageHandler, msg Message) error {
	fnValue := reflect.ValueOf(handler)
	eventValue := reflect.ValueOf(msg)

	var args []reflect.Value
	switch fnValue.Type().NumIn() {
	case 1:
		args = []reflect.Value{eventValue}
	case 2:
		ctxValue := reflect.ValueOf(ctx)
		args = []reflect.Value{ctxValue, eventValue}
	default:
		panic(fmt.Sprintf("msgmux: invalid handler function signature (got: %v)", fnValue.Type()))
	}

	out := fnValue.Call(args)
	if out[0].IsNil() {
		return nil
	}

	return out[0].Interface().(error)
}

func validateHandler(fn MessageHandler) error {
	fnType := reflect.TypeOf(fn)
	if fnType.Kind() != reflect.Func {
		return fmt.Errorf("msgmux: fn MessageHandler is not a function (got: %v)", fnType.Kind())
	}

	switch fnType.NumIn() {
	case 1:
		if fnType.In(0).Kind() != reflect.Struct {
			return fmt.Errorf("msgmux: fn MessageHandler input parameter should be a struct (got: %v)", fnType.In(0).Kind())
		}
	case 2:
		if fnType.In(0).Kind() != reflect.Interface {
			// expect context.Context interface, but got non-interface
			return fmt.Errorf("msgmux: fn MessageHandler 1st input parameter should be an context.Context (got: %v)", fnType.In(0).Kind())
		}

		if !fnType.In(0).Implements(reflect.TypeFor[context.Context]()) {
			return fmt.Errorf("msgmux: fn MessageHandler 1st input parameter should be context.Context (got: %v)", fnType.In(0).Kind())
		}

		if fnType.In(1).Kind() != reflect.Struct {
			return fmt.Errorf("msgmux: fn MessageHandler 2nd input parameter should be a struct (got: %v)", fnType.In(1).Kind())
		}
	default:
		return fmt.Errorf("msgmux: fn MessageHandler should have 1 or 2 input parameters (got: %d)", fnType.NumIn())
	}

	if fnType.NumOut() != 1 {
		return fmt.Errorf("msgmux: fn MessageHandler should have 1 output parameter (got: %d)", fnType.NumOut())
	}

	if fnType.Out(0).Kind() != reflect.Interface {
		// expect error interface, but got non-interface
		return fmt.Errorf("msgmux: fn MessageHandler output parameter should be an error (got: %v)", fnType.Out(0).Kind())
	}

	if fnType.Out(0).Name() != "error" {
		return fmt.Errorf("msgmux: fn MessageHandler output parameter should be an error (got: %v)", fnType.Out(0).Name())
	}

	return nil
}
