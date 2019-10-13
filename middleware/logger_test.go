package middleware

import (
	"bytes"
	"context"
	"log"
	"strings"
	"testing"
)

func TestLogger(t *testing.T) {
	ctx := context.Background()
	buf := new(bytes.Buffer)

	log.SetOutput(buf)

	if err := Logger(next)(ctx, &CreateUser{true}); err == nil {
		t.Errorf("failed to execute the middleware: expected an error and got %v", err)
	}

	if err := Logger(next)(ctx, &CreateUser{}); err != nil {
		t.Errorf("failed to execute the middleware: got %v and want %v", err, nil)
	}

	if !strings.Contains(buf.String(), "Executing command") {
		t.Errorf("invalid log message: %v", buf.String())
	}
}
