package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPathValidator_IsValidDirectory(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
		message  string
	}{
		{path: ".", expected: true, message: "current directory should be valid"},
		{path: "/", expected: true, message: "root directory should be valid"},
		{path: "./resolver_test.go", expected: false, message: "file is not a directory"},
		{path: "directory-in-a-galaxy-far-far-away", expected: false, message: "directory is not expected to exist"},
	}
	validator := NewPathValidator()
	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, validator.IsValidDirectory(testCase.path), testCase.message)
	}
}

func TestPathValidator_IsValidFile(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
		message  string
	}{
		{path: "./resolver_test.go", expected: true, message: "file should be a file"},
		{path: "file-in-a-galaxy-far-far-away", expected: false, message: "file is not expected to exist"},
		{path: ".", expected: false, message: "current directory is not a file"},
		{path: "/", expected: false, message: "root directory is not a file"},
	}
	validator := NewPathValidator()
	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, validator.IsValidFile(testCase.path), testCase.message)
	}
}
