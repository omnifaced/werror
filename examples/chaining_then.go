package examples

import (
	"fmt"

	"github.com/omnifaced/werror"
)

func processValue(x int) werror.Result[int] {
	if x < 0 {
		return werror.Err[int](fmt.Errorf("negative value"))
	}

	return werror.Ok(x * 2)
}

func validateRange(x int) werror.Result[int] {
	if x > 100 {
		return werror.Err[int](fmt.Errorf("value too large"))
	}

	return werror.Ok(x)
}

func chainingThen() {
	result := werror.Ok(10).
		Then(processValue).
		Then(validateRange)

	value, err := result.Unwrap()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("value:", value)
	}
}
