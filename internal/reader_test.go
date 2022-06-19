package internal

import (
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestReadFileInvalidFilePathFails(t *testing.T) {
	filePath := "./tariffs.csv"
	stringSlice, err := ReadFile(filePath)
	assert.Nil(t, stringSlice)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "open %s no such file or directory", filePath)
}

func TestReadFileSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := ReadTestEnv(envPath, logger)
	tariffFileLocation := envHandler.GetEnv("INTERNAL_TARIFFS_FILE_LOCATION")
	stringSlice, err := ReadFile(tariffFileLocation)
	assert.Nil(t, err)
	assert.NotNil(t, stringSlice)
	assert.IsType(t, []string{}, stringSlice)
}
