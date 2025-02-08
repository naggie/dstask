package integration

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtensibility(t *testing.T) {
	repo, cleanup := makeDstaskRepo(t)
	defer cleanup()

	workingDirectory, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	originalPath := os.Getenv("PATH")
	newPath := originalPath + string(os.PathListSeparator) + workingDirectory
	os.Setenv("PATH", newPath)

	f, err := os.OpenFile("dstask-extensibility", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
	if err != nil {
		t.Fatal("Failed to open extensibility script", err)
	}

	if _, err := f.Write([]byte("#!/usr/bin/env bash\necho \"Extensibility Test\"")); err != nil {
		t.Fatal("Failed to write extensibility script", err)
	}

	cleanup = func() {
		os.Remove("dstask-extensibility")
		os.Setenv("PATH", originalPath)
	}
	defer cleanup()

	program := testCmd(repo)
	output, _, success := program("extensibility")
	assert.Equal(t, "Extensibility Test\n", string(output))
	assert.True(t, success)
	assert.True(t, true)
}
