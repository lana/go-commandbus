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
}

func TestRegisteredHandler(t *testing.T) {
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
}
