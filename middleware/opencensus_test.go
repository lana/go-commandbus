package middleware

import (
	"context"
	"testing"
)

func TestOpenCensus(t *testing.T) {
	ctx := context.Background()

	mw := OpenCensus()

	if err := mw(next)(ctx, &CreateUser{true}); err == nil {
		t.Errorf("failed to execute the middleware: expected an error and got %v", err)
	}

	if err := mw(next)(ctx, &CreateUser{}); err != nil {
		t.Errorf("failed to execute the middleware: got %v and want %v", err, nil)
	}

	mw = OpenCensusWithConfig(OpenCensusConfig{})

	if err := mw(next)(ctx, &CreateUser{}); err != nil {
		t.Errorf("failed to execute the middleware: got %v and want %v", err, nil)
	}
}
