package examples

import (
	"fmt"
	"github.com/omnifaced/werror"
)

func getConfig() (string, error) {
	return "", fmt.Errorf("config not found")
}

func mustExample() {
	value := werror.Ok("success").Must()
	fmt.Println("value:", value)
}

func orExample() {
	defaultConfig := "default.yaml"

	config := werror.Wrap(getConfig()).Or(defaultConfig)
	fmt.Println("using config:", config)

	successConfig := werror.Ok("custom.yaml").Or(defaultConfig)
	fmt.Println("using config:", successConfig)
}
