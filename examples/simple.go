package examples

import (
	"fmt"
	"github.com/omnifaced/werror"
)

func simpleExample() {
	r := werror.Ok(10)

	value, err := r.Unwrap()
	if err != nil {
		fmt.Println("error:", err)
	} else {
		fmt.Println("value:", value)
	}
}
