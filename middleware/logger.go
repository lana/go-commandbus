package middleware

import (
	"context"
	"encoding/json"
	"log"
	"reflect"

	"github.com/lana/go-commandbus"
)

func Logger(next commandbus.HandlerFunc) commandbus.HandlerFunc {
	return func(ctx context.Context, cmd interface{}) error {
		if err := next(ctx, cmd); err != nil {
			return err
		}

		n := reflect.TypeOf(cmd).String()
		buf, _ := json.Marshal(cmd)
		log.Printf("Executing command \"%s\" with params: %s", n, buf)

		return nil
	}
}
