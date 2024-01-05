package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	t.Run("Success read", func(t *testing.T) {
		path := "testdata/env"

		env, err := ReadDir(path)

		assert.NotNil(t, env)
		assert.NoError(t, err)

		assert.Equal(t, "\"hello\"", env["HELLO"].Value)
		assert.Equal(t, "bar", env["BAR"].Value)
		assert.Equal(t, "   foo\nwith new line", env["FOO"].Value)
		assert.Equal(t, "", env["UNSET"].Value)
		assert.Equal(t, "", env["EMPTY"].Value)
	})

	t.Run("Dir not found", func(t *testing.T) {
		path := "testdata/notfound"

		env, err := ReadDir(path)

		assert.Nil(t, env)
		assert.Error(t, err)
	})
}
