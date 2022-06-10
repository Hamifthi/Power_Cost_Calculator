package pkg

import (
	"bufio"
	"fmt"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"io"
	"os"
	"strings"
	"sync"
)

func ReadAndProcessEfficiently(sessionFilePath string, tariffs []entity.Tariff) error {
	f, err := os.Open(sessionFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err = f.Close()
	}()
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 500*1024)
		return lines
	}}
	stringsPool := sync.Pool{New: func() interface{} {
		strs := ""
		return strs
	}}
	r := bufio.NewReader(f)
	// skip the header
	_, err = r.ReadBytes('\n')
	if err != nil {
		return err
	}
	var wg sync.WaitGroup
	for {
		buf := linesPool.Get().([]byte)
		n, err := r.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
		}
		nextUntilNewline, err := r.ReadBytes('\n')
		if err != io.EOF {
			buf = append(buf, nextUntilNewline...)
		}
		wg.Add(1)
		go func() {
			processChunk(buf, &linesPool, &stringsPool, tariffs)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func processChunk(chunk []byte, linesPool, stringPool *sync.Pool, tariffs []entity.Tariff) {
	sessionStrings := stringPool.Get().(string)
	sessionStrings = string(chunk)

	linesPool.Put(chunk)

	sessionsSlice := strings.Split(sessionStrings, "\n")
	stringPool.Put(sessionStrings)

	sessions, err := internal.ParseSession(sessionsSlice)
	if err != nil {
		fmt.Println(err)
		return
	}
	costs, err := internal.CostCalculator(tariffs, sessions)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = internal.WriteFile("./data/calculated_costs.csv", costs)
	if err != nil {
		fmt.Println(err)
		return
	}
	return
}
