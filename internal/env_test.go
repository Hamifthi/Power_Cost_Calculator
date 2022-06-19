package internal

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"testing"
)

var envPath string

func init() {
	currentDir, _ := os.Getwd()
	rootDir := filepath.Dir(currentDir)
	envPath = rootDir + "/test.env"
}

func TestGetEnvSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := ReadTestEnv(envPath, logger)
	value := envHandler.GetEnv("TEST")
	assert.Equal(t, value, "test")
}
