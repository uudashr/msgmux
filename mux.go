package msgmux

import (
	"fmt"
	"reflect"
)

// Message represent specific message. It must be struct type.
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

// MssageHandler is a function type to handle the message.
//
// The form of function is:
//
//	func(msg Message) error
//
// where the [Message] represent the specific message type (struct).
//
// Example:
//
//	mux := &msgmux.DispatchMux{}
//	mux.Handle(func(event OrderCompleted) error {
//		// handle the event
//		return nil
//	})
type MessageHandler any

type DispatchMux struct {
	handlers map[reflect.Type]MessageHandler
}

func NewDispatchMux() *DispatchMux {
	return &DispatchMux{}
}

func (m *DispatchMux) Handle(fn MessageHandler) {
	if err := validateHandler(fn); err != nil {
		panic(err)
	}

	if m.handlers == nil {
		m.handlers = make(map[reflect.Type]MessageHandler)
	}

	fnType := reflect.TypeOf(fn)
	fnTypeIn := fnType.In(0)

	if _, reg := m.handlers[fnTypeIn]; reg {
		panic(fmt.Sprintf("msgmux: handler for message %v already registered", fnTypeIn.Name()))
	}

	m.handlers[fnTypeIn] = fn
}

func (m *DispatchMux) Dispatch(msg Message) error {
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

	return invokeHandler(handler, msg)
}

func invokeHandler(handler MessageHandler, msg Message) error {
	fnValue := reflect.ValueOf(handler)
	eventValue := reflect.ValueOf(msg)
	out := fnValue.Call([]reflect.Value{eventValue})
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

	if fnType.NumIn() != 1 {
		return fmt.Errorf("msgmux: fn MessageHandler should have 1 input parameter (got: %d)", fnType.NumIn())
	}

	if fnType.NumOut() != 1 {
		return fmt.Errorf("msgmux: fn MessageHandler should have 1 output parameter (got: %d)", fnType.NumOut())
	}

	if fnType.In(0).Kind() != reflect.Struct {
		return fmt.Errorf("msgmux: fn MessageHandler input parameter should be a struct (got: %v)", fnType.In(0).Kind())
	}

	if fnType.Out(0).Kind() != reflect.Interface {
		return fmt.Errorf("msgmux: fn MessageHandler output parameter should be an interface (got: %v)", fnType.Out(0).Kind())
	}

	if fnType.Out(0).Name() != "error" {
		return fmt.Errorf("msgmux: fn MessageHandler output parameter should be an error (got: %v)", fnType.Out(0).Name())
	}

	return nil
}
