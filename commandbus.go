package main

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

// CommandBus is the definition of how command should be handled
type CommandBus interface {
	// Register assign a Command to a CommandHandle for
	// future executions
	Register(string, interface{}) error

	// Execute send a given Command to its assigned CommandHandler
	Execute(context.Context, interface{}) error
}

type bus struct {
	handlers map[string]interface{}
}

func (b *bus) Register(cmdName string, fn interface{}) error {
	if reflect.TypeOf(fn).Kind() != reflect.Func {
		return fmt.Errorf("%s is not a reflect.Func", reflect.TypeOf(fn))
	}

	if _, err := b.handler(cmdName); err == nil {
		return errors.New("command already asigned to a handler")
	}

	b.handlers[cmdName] = fn

	return nil
}

func (b *bus) Execute(ctx context.Context, cmd interface{}) error {
	handler, err := b.handler(reflect.TypeOf(cmd).String())
	if err != nil {
		return err
	}

	fn := reflect.ValueOf(handler)
	args := b.buildHandlerArgs([]interface{}{ctx, cmd})

	fn.Call(args) // todo return errors here
	return nil
}

func (b *bus) handler(cmdName string) (interface{}, error) {
	if h, ok := b.handlers[cmdName]; ok {
		return h, nil
	}

	return nil, errors.New("handler not found for command")
}

func (b *bus) buildHandlerArgs(args []interface{}) []reflect.Value {
	reflectedArgs := make([]reflect.Value, 0)

	for _, arg := range args {
		reflectedArgs = append(reflectedArgs, reflect.ValueOf(arg))
	}

	return reflectedArgs
}

func New() CommandBus {
	return &bus{
		handlers: make(map[string]interface{}, 0),
	}
}
