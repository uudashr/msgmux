package msgmux_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/uudashr/msgmux"
)

func TestDispatch(t *testing.T) {
	type OrderCompleted struct {
		OrderID string
	}

	mux := msgmux.NewDispatchMux()
	mux.Handle(func(ctx context.Context, event OrderCompleted) error {
		return nil
	})

	err := mux.DispatchContext(context.Background(), OrderCompleted{OrderID: "order-123"})
	require.NoError(t, err)
}

func TestDispatchMux_Handle(t *testing.T) {
	type OrderCompleted struct {
		OrderID string
	}

	tests := []struct {
		name        string
		handleFunc  msgmux.MessageHandler
		expectPanic bool
	}{
		{
			name: "valid handler without context",
			handleFunc: func(event OrderCompleted) error {
				return nil
			},
		},
		{
			name: "valid handler with context",
			handleFunc: func(ctx context.Context, event OrderCompleted) error {
				return nil
			},
		},
		{
			name:        "string value as handler",
			handleFunc:  "invalid",
			expectPanic: true,
		},
		{
			name:        "int value as handler",
			handleFunc:  123,
			expectPanic: true,
		},
		{
			name:        "nil value as handler",
			handleFunc:  nil,
			expectPanic: true,
		},
		{
			name: "non struct arg handler without context",
			handleFunc: func(event string) error {
				return nil
			},
			expectPanic: true,
		},
		{
			name: "non struct arg handler with context",
			handleFunc: func(ctx context.Context, event string) error {
				return nil
			},
			expectPanic: true,
		},
		{
			name: "no return error handler without context",
			handleFunc: func(event OrderCompleted) {
			},
			expectPanic: true,
		},
		{
			name: "no return error handler with context",
			handleFunc: func(ctx context.Context, event OrderCompleted) {
			},
			expectPanic: true,
		},
		{
			name: "empty arg handler without error",
			handleFunc: func() {
			},
			expectPanic: true,
		},
		{
			name: "empty arg handler with error",
			handleFunc: func() error {
				return nil
			},
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mux := msgmux.NewDispatchMux()

			if tt.expectPanic {
				require.Panics(t, func() {
					mux.Handle(tt.handleFunc)
				})
				return
			}

			require.NotPanics(t, func() {
				mux.Handle(tt.handleFunc)
			})
		})
	}
}
