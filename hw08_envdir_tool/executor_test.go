package main

import (
	"testing"
)

func TestRunCmd(t *testing.T) {
	testCases := []struct {
		name     string
		env      Environment
		cmd      []string
		expected int
	}{
		{
			name:     "empty command",
			env:      Environment{},
			cmd:      []string{},
			expected: 1,
		},
		{
			name:     "existing command",
			env:      Environment{},
			cmd:      []string{"echo", "hello"},
			expected: 0,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			got := RunCmd(tt.cmd, tt.env)

			if got != tt.expected {
				t.Errorf("RunCmd() got = %v, expected %v", got, tt.expected)
			}
		})
	}
}
