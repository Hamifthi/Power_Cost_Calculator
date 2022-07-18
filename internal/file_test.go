package internal

import (
	"encoding/csv"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCreateCSVWriterAndWriteHeader(t *testing.T) {
	fileName := "./something"
	file, _ := os.Create(fileName)
	csvWriter, err := CreateCSVWriterAndWriteHeader([]string{"something"}, file)
	assert.Nil(t, err)
	assert.IsType(t, &csv.Writer{}, csvWriter)
	_ = os.Remove(fileName)
}
