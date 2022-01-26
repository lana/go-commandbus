package middleware

import (
	"context"

	"go.lana.tech/errors"
)

// CreateUser command for testing.
type CreateUser struct {
	Err bool
}

// next command handler for testing.
func next(ctx context.Context, cmd interface{}) error {
	v := cmd.(*CreateUser)

	if v.Err {
		return errors.New("command was failed")
	}

	return nil
}
