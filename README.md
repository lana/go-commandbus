# GO CommandBus

[![Build Status](https://img.shields.io/travis/vsmoraes/go-commandbus/master.svg?style=flat-square)](https://travis-ci.org/vsmoraes/go-commandbus)
[![Codecov branch](https://img.shields.io/codecov/c/github/vsmoraes/go-commandbus/master.svg?style=flat-square)](https://codecov.io/gh/vsmoraes/go-commandbus)
[![GoDoc](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](https://godoc.org/github.com/vsmoraes/go-commandbus)
[![Go Report Card](https://goreportcard.com/badge/github.com/vsmoraes/go-commandbus?style=flat-square)](https://goreportcard.com/report/github.com/vsmoraes/go-commandbus)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](https://github.com/vsmoraes/go-commandbus/blob/master/LICENSE)

A slight and pluggable command-bus for Go.

## Install

Use go get.
```sh
$ go get github.com/vsmoraes/go-commandbus
```

Then import the package into your own code:
```
import "github.com/vsmoraes/go-commandbus"
```

## Usage
```go
package main

import (
	"context"
	"log"

	"github.com/vsmoraes/go-commandbus"
)

type CreateUser struct {
	Name string
}

func CreateHandler(ctx context.Context, cmd *CreateUser) error {
	log.Printf("user %s created", cmd.Name)

	return nil
}

func main() {
	bus := commandbus.New()

	err := bus.Register(&CreateUser{}, CreateHandler)

	if err != nil {
		log.Fatal(err)
	}

	err = bus.Execute(context.Background(), &CreateUser{"go-commandbus"})

	if err != nil {
		log.Fatal(err)
	}
}
```

## License

This project is released under the MIT licence. See [LICENSE](https://github.com/vsmoraes/go-commandbus/blob/master/LICENSE) for more details.
