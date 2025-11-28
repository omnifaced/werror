package examples

import (
	"fmt"
	"github.com/omnifaced/werror"
)

func divide(a, b int) (int, error) {
	if b == 0 {
		return 0, fmt.Errorf("division by zero")
	}

	return a / b, nil
}

func sideEffectsExample() {
	werror.Wrap(divide(10, 2)).
		OnSuccess(func(value int) {
			fmt.Println("success:", value)
		}).
		OnError(func(err error) {
			fmt.Println("error occurred:", err)
		}).
		Always(func() {
			fmt.Println("operation completed")
		})

	fmt.Println("---")

	werror.Wrap(divide(10, 0)).
		OnSuccess(func(value int) {
			fmt.Println("success:", value)
		}).
		OnError(func(err error) {
			fmt.Println("error occurred:", err)
		}).
		Always(func() {
			fmt.Println("operation completed")
		})
}
