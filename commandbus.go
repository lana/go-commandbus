package commandbus

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// HandlerFunc defines a function to execute the command handler.
// This function type is only used by `MiddlewareFunc`.
type HandlerFunc func(context.Context, interface{}) error

// MiddlewareFunc defines a function to process middleware.
// it receives the next Handler must return another handler
type MiddlewareFunc func(h HandlerFunc) HandlerFunc

type handlers map[string]interface{}

// CommandBus is the definition of how command should be handled
type CommandBus interface {
	// Register assign a command to a command handle for future executions.
	Register(interface{}, interface{}) error

	// Handlers returns all registered handlers.
	Handlers() handlers

	// Use adds middleware to the chain.
	Use(...MiddlewareFunc)

	// Execute send a given Command to its assigned command handler.
	Execute(context.Context, interface{}) error
}

type bus struct {
	handlers   handlers
	middleware []MiddlewareFunc
}

// New creates a new command bus.
func New() CommandBus {
	return &bus{
		handlers:   make(handlers),
		middleware: make([]MiddlewareFunc, 0),
	}
}

// Register assign a command to a command handle for future executions.
func (b *bus) Register(cmd interface{}, fn interface{}) error {
	if err := b.validate(cmd); err != nil {
		return err
	}

	cmdName := reflect.TypeOf(cmd).String()

	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a function", reflect.TypeOf(fn))
	}

	if _, err := b.handler(cmdName); err == nil {
		return errors.New("command already assigned to a handler")
	}

	b.handlers[cmdName] = fn

	return nil
}

// Handlers returns all registered handlers.
func (b *bus) Handlers() handlers {
	return b.handlers
}

// Use adds middleware to the chain.
func (b *bus) Use(middleware ...MiddlewareFunc) {
	b.middleware = append(b.middleware, middleware...)
}

// Execute send a given Command to its assigned command handler.
func (b *bus) Execute(ctx context.Context, cmd interface{}) error {
	if err := b.validate(cmd); err != nil {
		return err
	}

	h := applyMiddleware(b.executor, b.middleware...)

	return h(ctx, cmd)
}

func (b *bus) handler(cmdName string) (interface{}, error) {
	if h, ok := b.handlers[cmdName]; ok {
		return h, nil
	}

	return nil, errors.New("handler not found for command")
}

func applyMiddleware(h HandlerFunc, middleware ...MiddlewareFunc) HandlerFunc {
	for i := len(middleware) - 1; i >= 0; i-- {
		h = middleware[i](h)
	}

	return h
}

func (b *bus) executor(ctx context.Context, cmd interface{}) error {
	handler, err := b.handler(reflect.TypeOf(cmd).String())

	if err != nil {
		return err
	}

	fn := reflect.ValueOf(handler)
	args := make([]reflect.Value, 0)

	for _, arg := range []interface{}{ctx, cmd} {
		args = append(args, reflect.ValueOf(arg))
	}

	res := make([]interface{}, 0)

	for _, r := range fn.Call(args) {
		res = append(res, r.Interface())
	}

	if len(res) < 1 {
		return nil
	}

	r := res[:1][0]

	if err, ok := r.(error); ok {
		return err
	}

	return nil
}

func (b *bus) validate(cmd interface{}) error {
	value := reflect.ValueOf(cmd)

	if value.Kind() != reflect.Ptr || !value.IsNil() && value.Elem().Kind() != reflect.Struct {
		return errors.New("only pointer to commands are allowed")
	}

	return nil
}
