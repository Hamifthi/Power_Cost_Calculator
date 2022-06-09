package internal

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ReadFile(filePath string) ([]string, error) {
	var stringSlice []string
	f, err := os.Open(filePath)
	defer func() {
		err = f.Close()
	}()
	if err != nil {
		return stringSlice, err
	}
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		stringSlice = append(stringSlice, line)
		if err != nil {
			if err == io.EOF {
				break
			}
			return stringSlice, err
		}
	}
	return stringSlice, nil
}
