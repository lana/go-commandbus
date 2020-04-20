package middleware

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/lana/go-commandbus"
	"github.com/lana/go-dispatcher"
)

const (
	// DispatcherReceived event raised when a dispatcher has been received.
	DispatcherReceived = "received"

	// DispatcherExecuted event raised when a dispatcher has been executed.
	DispatcherExecuted = "executed"
)

type event struct {
	name string
	cmd  interface{}
}

// Type returns the event type.
func (e event) Type() dispatcher.EventType {
	name := reflect.TypeOf(e.cmd).Elem().Name()
	name = fmt.Sprintf("%s.%s", strings.ToLower(name), e.name)

	return dispatcher.EventType(name)
}

// Data returns the event command data.
func (e event) Data() interface{} {
	return e.cmd
}

// Dispatcher will dispatch events when a command was received and executed.
func Dispatcher(d dispatcher.Dispatcher) commandbus.MiddlewareFunc {
	return func(next commandbus.HandlerFunc) commandbus.HandlerFunc {
		return func(ctx context.Context, cmd interface{}) error {
			d.Dispatch(ctx, event{DispatcherReceived, cmd})

			if err := next(ctx, cmd); err != nil {
				return err
			}

			d.Dispatch(ctx, event{DispatcherExecuted, cmd})

			return nil
		}
	}
}
