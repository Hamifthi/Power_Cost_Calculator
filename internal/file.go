package internal

import (
	"encoding/csv"
	"os"
)

func CreateCSVWriterAndWriteHeader(header []string, file *os.File) (*csv.Writer, error) {
	writer := csv.NewWriter(file)
	err := writer.Write(header)
	if err != nil {
		return nil, err
	}
	return writer, nil
}
