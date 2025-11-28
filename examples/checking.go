package examples

import (
	"fmt"

	"github.com/omnifaced/werror"
)

func processData(data string) (string, error) {
	if data == "" {
		return "", fmt.Errorf("empty data")
	}

	return "processed: " + data, nil
}

func checkingExample() {
	result := werror.Wrap(processData("hello"))

	if result.IsOk() {
		fmt.Println("success:", result.Value())
	}

	if result.IsErr() {
		fmt.Println("error:", result.Error())
	}

	failResult := werror.Wrap(processData(""))
	if failResult.IsErr() {
		fmt.Println("failed as expected:", failResult.Error())
	}
}
