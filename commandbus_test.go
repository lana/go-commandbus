package commandbus

import (
	"context"
	"errors"
	"reflect"
	"testing"
)

type Cmd struct {
	c int
}

func CmdHandler(ctx context.Context, cmd *Cmd) error {
	cmd.c++
	return nil
}

type ErrorCmd struct {
	err error
}

func ErrorCmdHandler(ctx context.Context, cmd *ErrorCmd) error {
	return cmd.err
}

func TestRegisterHandlers(t *testing.T) {
	bus := New()
	if err := bus.Register(&Cmd{}, CmdHandler); err != nil {
		t.Errorf("Failed to register command: %v", err)
	}

	if len(bus.Handlers()) != 1 {
		t.Errorf("Expected to have 1 command, got %d", len(bus.Handlers()))
	}

	for _, handler := range bus.Handlers() {
		if reflect.ValueOf(handler) != reflect.ValueOf(CmdHandler) {
			t.Error("Registered handler is different from the expected one")
		}
	}

	// invalid handler
	fooCmd := struct{}{}
	invalidHandler := struct{}{}
	if err := bus.Register(fooCmd, invalidHandler); err == nil {
		t.Error("Invalid handler was accepted. Register must only accept functions")
	}

	// duplicated
	if err := bus.Register(&Cmd{}, CmdHandler); err == nil {
		t.Error("Bus must not accept duplicated commands")
	}
}

func TestExecuteRegisteredHandler(t *testing.T) {
	bus := New()
	if err := bus.Register(&Cmd{}, CmdHandler); err != nil {
		t.Errorf("Failed to register command: %v", err)
	}

	c := &Cmd{}
	expected := 10

	for i := 0; i < expected; i++ {
		if err := bus.Execute(context.Background(), c); err != nil {
			t.Errorf("Failed to execute command: %v", err)
		}
	}

	if c.c != expected {
		t.Errorf("Execution number is wrong. Expected %d, got %d", expected, c.c)
	}

	// try to execute an unregistered command
	invalidCommand := struct{}{}
	if err := bus.Execute(context.Background(), invalidCommand); err == nil {
		t.Error("Invalid command was executed without errors")
	}
}

func TestExecutionReturnedError(t *testing.T) {
	bus := New()
	if err := bus.Register(&ErrorCmd{}, ErrorCmdHandler); err != nil {
		t.Errorf("Failed to register command: %v", err)
	}

	cmd := &ErrorCmd{err: errors.New("this is an error")}
	err := bus.Execute(context.Background(), cmd)
	if err != cmd.err {
		t.Errorf("Failed to assert command error. Expected \"%v\", got \"%v\"", cmd.err, err)
	}
}

func TestExecutionWithMiddleware(t *testing.T) {
	bus := New()
	if err := bus.Register(&Cmd{}, CmdHandler); err != nil {
		t.Errorf("Failed to register command: %v", err)
	}

	// this middleware will be executed before the handler,
	// making the command counter equals to 10.
	bus.Use(func(next interface{}) interface{} {
		return func(ctx context.Context, cmd interface{}) interface{} {
			if c, ok := cmd.(*Cmd); ok {
				c.c += 10 // add 10 to the command value
			}

			return next
		}
	})

	c := &Cmd{}
	expected := 10
	if err := bus.Execute(context.Background(), c); err != nil {
		t.Errorf("Failed to execute command: %v", err)
	}

	if c.c != expected {
		t.Errorf("Execution number is wrong. Expected %d, got %d", expected, c.c)
	}
}
