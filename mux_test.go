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
