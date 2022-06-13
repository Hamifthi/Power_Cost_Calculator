package calculator

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/pkg/errors"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

type CostCalculator struct {
	logger             *log.Logger
	mu                 sync.Mutex
	linesPool          *sync.Pool
	stringsPool        *sync.Pool
	outputFileLocation string
}

func New(logger *log.Logger, linesPool, stringsPool *sync.Pool, outputFileLocation string) CostCalculator {
	return CostCalculator{
		logger:             logger,
		linesPool:          linesPool,
		stringsPool:        stringsPool,
		outputFileLocation: outputFileLocation,
	}
}

func (cc *CostCalculator) ReadAndProcessSessions(sessionFilePath string, tariffs []entity.Tariff) error {
	f, err := os.Open(sessionFilePath)
	if err != nil {
		cc.logger.Printf("Error reading session files because of %v", err)
		return err
	}
	defer func() {
		err = f.Close()
	}()
	r := bufio.NewReader(f)
	// skip the header
	_, err = r.ReadBytes('\n')
	if err != nil {
		cc.logger.Printf("Error reading the header of the session file because of %v", err)
		return err
	}
	var wg sync.WaitGroup
	for {
		buf := cc.linesPool.Get().([]byte)
		n, err := r.Read(buf)
		buf = buf[:n]
		if n == 0 {
			if err == io.EOF {
				break
			}
			if err != nil {
				cc.logger.Printf("Error reading the file into buffer because of %v", err)
				return err
			}
		}
		nextUntilNewline, err := r.ReadBytes('\n')
		if err != io.EOF {
			cc.logger.Printf("Error reading the remaining line until next line because of %v", err)
			buf = append(buf, nextUntilNewline...)
		}
		wg.Add(1)
		go func() {
			cc.processChunk(buf, tariffs)
			wg.Done()
		}()
	}
	wg.Wait()
	return nil
}

func (cc *CostCalculator) processChunk(chunk []byte, tariffs []entity.Tariff) {
	sessionStrings := cc.stringsPool.Get().(string)
	sessionStrings = string(chunk)

	cc.linesPool.Put(chunk)

	sessionsSlice := strings.Split(sessionStrings, "\n")
	cc.stringsPool.Put(sessionStrings)

	sessions, err := internal.ParseSession(sessionsSlice)
	if err != nil {
		cc.logger.Printf("Error parsing sessions because of %v", err)
		return
	}
	costs := internal.CostCalculator(tariffs, sessions)
	err = cc.WriteCosts(costs)
	if err != nil {
		cc.logger.Printf("Error writing costs to output file because of %v", err)
		return
	}
	return
}

func (cc *CostCalculator) WriteCosts(costs []entity.Cost) error {
	var writer *csv.Writer
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if _, err := os.Stat(cc.outputFileLocation); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(cc.outputFileLocation)
		if err != nil {
			cc.logger.Printf("Error creating output file because of %v", err)
			return err
		}
		defer func() {
			err = file.Close()
		}()
		writer = csv.NewWriter(file)
		err = writer.Write([]string{"session_id", "total_cost"})
		if err != nil {
			cc.logger.Printf("Error writing output file header because of %v", err)
			return err
		}
	} else {
		file, err := os.OpenFile(cc.outputFileLocation, os.O_APPEND|os.O_WRONLY, 0744)
		if err != nil {
			cc.logger.Printf("Error opening output file because of %v", err)
			return err
		}
		defer func() {
			err = file.Close()
		}()
		writer = csv.NewWriter(file)
	}
	defer writer.Flush()
	for _, cost := range costs {
		row := []string{cost.SessionID, fmt.Sprintf("%.3f", cost.TotalCost)}
		if err := writer.Write(row); err != nil {
			cc.logger.Printf("Error writing to output file because of %v", err)
			return err
		}
	}
	return nil
}

func (cc *CostCalculator) ReadAndParseTariffs(filePath string) ([]entity.Tariff, error) {
	var tariffs []entity.Tariff
	tariffStrings, err := internal.ReadFile(filePath)
	if err != nil {
		cc.logger.Printf("Error reading tariffs file because of %v", err)
		return tariffs, err
	}
	tariffs, err = internal.ParseTariff(tariffStrings[1:])
	if err != nil {
		cc.logger.Printf("Error parsing tariffs file because of %v", err)
		return tariffs, err
	}
	return tariffs, err
}
