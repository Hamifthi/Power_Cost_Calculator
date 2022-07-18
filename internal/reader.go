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
	if err != nil {
		return stringSlice, err
	}
	defer func() {
		err = f.Close()
	}()
	buf := bufio.NewReader(f)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		// it's because not adding empty line to the stringSlice
		if len(line) != 0 {
			stringSlice = append(stringSlice, line)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return stringSlice, err
		}
	}
	return stringSlice, err
}
