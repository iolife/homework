package main

import (
	"fmt"

	"github.com/pkg/errors"
)

func ErrNoRow() error {
	return errors.New("sql: no rows in result set")
}

func test() error {
	// dosomething
	e := ErrNoRow()
	if e != nil {
		return errors.Wrap(e, "not found")
	}
	//another
	return nil
}

func main() {
	e := test()
	if e != nil {
		fmt.Printf("original error: %T %v\n", errors.Cause(e), errors.Cause(e))
		fmt.Printf("stack trace:\n%+v\n", e)
	}
}
