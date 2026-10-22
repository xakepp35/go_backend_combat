package memory

import (
	"context"
	"errors"
	"jobflow/pkg"
	"testing"
)

func TestCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := New().List(ctx, pkg.ListTasksRequest{}); !errors.Is(err, context.Canceled) {
		t.Fatal(err)
	}
}
