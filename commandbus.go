package commandbus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// MiddlewareFunc defines a function to process middleware.
// it receives the next Handler must return another handler
type MiddlewareFunc func(interface{}) interface{}

// CommandBus is the definition of how command should be handled
type CommandBus interface {
	// Register assign a Command to a CommandHandle for
	// future executions
	Register(interface{}, interface{}) error

	// Handlers returns all registered handlers
	Handlers() map[string]interface{}

	// Use adds middleware to the chain
	Use(...MiddlewareFunc)

	// Execute send a given Command to its assigned command handler
	Execute(context.Context, interface{}) error
}

type bus struct {
	handlers   map[string]interface{}
	middleware []MiddlewareFunc
}

func New() CommandBus {
	return &bus{
		handlers:   make(map[string]interface{}, 0),
		middleware: make([]MiddlewareFunc, 0),
	}
}

func (b *bus) Register(cmd interface{}, fn interface{}) error {
	if err := validCmd(cmd); err != nil {
		return err
	}

	cmdName := reflect.TypeOf(cmd).String()

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	if _, err := b.handler(cmdName); err == nil {
		return errors.New("command already assigned to a handler")
	}

	b.handlers[cmdName] = fn

	return nil
}

func (b *bus) Handlers() map[string]interface{} {
	return b.handlers
}

func (b *bus) Use(middleware ...MiddlewareFunc) {
	b.middleware = append(b.middleware, middleware...)
}

func (b *bus) Execute(ctx context.Context, cmd interface{}) error {
	if err := validCmd(cmd); err != nil {
		return err
	}

	handler, err := b.handler(reflect.TypeOf(cmd).String())
	if err != nil {
		return err
	}

	fn := reflect.ValueOf(chain(handler, b.middleware...))
	args := buildHandlerArgs([]interface{}{ctx, cmd})

	results := buildHandlerResults(fn.Call(args))
	if len(results) < 1 {
		return nil
	}

	r := results[:1][0]
	if err, ok := r.(error); ok {
		return err
	}

	return nil
}

func (b *bus) handler(cmdName string) (interface{}, error) {
	if h, ok := b.handlers[cmdName]; ok {
		return h, nil
	}

	return nil, errors.New("handler not found for command")
}

func buildHandlerResults(res []reflect.Value) []interface{} {
	results := make([]interface{}, 0)

	for _, r := range res {
		results = append(results, r.Interface())
	}

	return results
}

func buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}

func chain(h interface{}, middleware ...MiddlewareFunc) interface{} {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

func validCmd(cmd interface{}) error {
	value := reflect.ValueOf(cmd)
	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer to commands are allowed")
	}

	return nil
}
