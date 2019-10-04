package commandbus

import (
	"context"
	"reflect"
	"testing"
)

type Cmd struct {
	c int
}

func CmdHandler(ctx context.Context, cmd *Cmd) {
	cmd.c++
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
	fooCmd := struct {}{}
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
	invalidCommand := struct {}{}
	if err := bus.Execute(context.Background(), invalidCommand); err == nil {
		t.Error("Invalid command was executed without errors")
	}
}
