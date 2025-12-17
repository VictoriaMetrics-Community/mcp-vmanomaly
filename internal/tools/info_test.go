package tools

import (
	"context"
	"errors"
	"testing"
)

func TestGetBuildinfo_Error(t *testing.T) {
	mock := &MockClient{
		GetBuildInfoFunc: func(ctx context.Context) (map[string]any, error) {
			return nil, errors.New("API error")
		},
	}

	_, err := mock.GetBuildInfo(context.Background())

	if err == nil {
		t.Error("expected error from API")
	}
}
