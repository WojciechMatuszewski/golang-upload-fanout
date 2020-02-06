package env_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"testing-stuff/pkg/env"
)

func TestRequire(t *testing.T) {
	t.Run("panics when variables does not exist", func(t *testing.T) {
		assert.Panics(t, func() {
			env.Require("random_name")
		})
	})

	t.Run("returns correct variable value when such variable exist", func(t *testing.T) {
		k := "ENV_VARIABLE"
		v := "ENV_VARIABLE_VALUE"
		unsetFn := Setenv(t, k, v)
		defer unsetFn()

		got := env.Require(k)
		assert.Equal(t, v, got)
	})

	t.Run("panics on empty value", func(t *testing.T) {
		unsetFn := Setenv(t, "FOO", "")
		defer unsetFn()

		assert.Panics(t, func() {
			env.Require("FOO")
		})
	})
}

func Setenv(t *testing.T, key, value string) func() {
	t.Helper()

	prev := os.Getenv(key)

	err := os.Setenv(key, value)
	if err != nil {
		t.FailNow()
	}

	return func() {
		err := os.Setenv(key, prev)
		if err != nil {
			t.FailNow()
		}
	}
}
