package middleware

import (
	"context"
	"encoding/json"
	"log"
	"reflect"
)

func Logger(next interface{}) interface{} {
	return func(ctx context.Context, cmd interface{}) interface{} {
		n := reflect.TypeOf(cmd).String()
		buf, _ := json.Marshal(cmd)
		log.Printf("Executing command \"%s\" with params: %s", n, buf)

		return next
	}
}
