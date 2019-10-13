package middleware

import (
	"context"
	"reflect"
	"testing"

	"github.com/vsmoraes/go-dispatcher"
)

func TestDispatcher(t *testing.T) {
	d := dispatcher.New()
	d.On("createuser.received", func(ctx context.Context, e dispatcher.Event) {
		if _, ok := (e.Data()).(*CreateUser); !ok {
			t.Errorf("invalid data retrieved: got %v", reflect.TypeOf(e.Data()))
		}
	})

	mw := Dispatcher(d)
	ctx := context.Background()

	if err := mw(next)(ctx, &CreateUser{true}); err == nil {
		t.Errorf("failed to execute the middleware: expected an error and got %v", err)
	}

	if err := mw(next)(ctx, &CreateUser{}); err != nil {
		t.Errorf("failed to execute the middleware: got %v and want %v", err, nil)
	}
}
