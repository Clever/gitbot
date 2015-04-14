package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type normalizePathTestCase struct {
	configPath  string
	commandPath string
	expected    string
}

var normalizePathTests = []normalizePathTestCase{
	{"/path/to/gitbot.yml", "/absolute/path/command.py", "/absolute/path/command.py"},
	{"/path/to/gitbot.yml", "./some/rel/path.py", "/path/to/some/rel/path.py"},
	{"/path/to/gitbot.yml", "git", "git"},
}

func TestNormalizePath(t *testing.T) {
	for _, testCase := range normalizePathTests {
		actual := NormalizePath(testCase.configPath, testCase.commandPath)
		assert.Equal(t, testCase.expected, actual, "NormalizePath(%v, %v) failed", testCase.configPath, testCase.commandPath)
	}
}
