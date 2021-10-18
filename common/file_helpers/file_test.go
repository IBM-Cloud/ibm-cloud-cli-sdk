package file_helpers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileExists(t *testing.T) {
	// Setup
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal("Could not create temporary file for test")
	}

	// Assertions
	assert.True(t, FileExists(file.Name()), "file test should exist")

	// Cleanup
	os.RemoveAll(file.Name())
}

func TestFileExistsOnFolder(t *testing.T) {
	// Setup
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal("Could not create temporary directory for test")
	}

	// Assertions
	assert.False(t, FileExists(dir), "a directory should not return true")

	// Cleanup
	os.RemoveAll(dir)
}

func TestFileExistNonExistentFile(t *testing.T) {

	assert.False(t, FileExists(os.TempDir() + "/non-existent-file"), "non existent file should return false")
}

