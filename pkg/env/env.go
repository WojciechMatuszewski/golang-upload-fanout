package env

import (
	"fmt"
	"os"
)

func Require(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok || value == "" {
		panic(fmt.Sprintf("%v enviroment variable does not exist", key))
	}
	return value
}
