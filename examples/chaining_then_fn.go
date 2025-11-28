package examples

import (
	"fmt"

	"github.com/omnifaced/werror"
)

func double(x int) (int, error) {
	return x * 2, nil
}

func failIfOdd(x int) (int, error) {
	if x%2 == 1 {
		return 0, fmt.Errorf("odd number")
	}

	return x, nil
}

func chaining() {
	res := werror.Ok(10).
		ThenFn(double).
		ThenFn(failIfOdd)

	v, err := res.Unwrap()
	fmt.Println(v, err)
}
