package tools

import (
	"context"
	"errors"
	"testing"
)

func TestHealthCheck_Error(t *testing.T) {
	mock := &MockClient{
		GetHealthFunc: func(ctx context.Context) (map[string]any, error) {
			return nil, errors.New("connection refused")
		},
	}

	_, err := mock.GetHealth(context.Background())

	if err == nil {
		t.Error("expected error from API")
	}
}
