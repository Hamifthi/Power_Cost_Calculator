package calculator

import (
	"bufio"
	"encoding/csv"
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
	pools              *entity.SyncPools
	OutputFileLocation string
	numberOfThreads    int
}

func New(logger *log.Logger, pools *entity.SyncPools, outputFileLocation string, numberOfThreads int) *CostCalculator {
	return &CostCalculator{
		logger:             logger,
		pools:              pools,
		OutputFileLocation: outputFileLocation,
		numberOfThreads:    numberOfThreads,
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
	// maximum number of concurrent go routines
	waitChan := make(chan struct{}, cc.numberOfThreads)
	for {
		waitChan <- struct{}{}
		buf := cc.pools.LinesPool.Get().([]byte)
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
			buf = append(buf, nextUntilNewline[:len(nextUntilNewline)-1]...)
		}
		wg.Add(1)
		go func() {
			cc.processChunk(buf, tariffs)
			wg.Done()
			<-waitChan
		}()
	}
	wg.Wait()
	return nil
}

func (cc *CostCalculator) processChunk(chunk []byte, tariffs []entity.Tariff) {
	sessionStrings := cc.pools.StringsPool.Get().(string)
	sessionStrings = string(chunk)

	cc.pools.LinesPool.Put(chunk)

	sessionsSlice := strings.Split(sessionStrings, "\n")
	cc.pools.StringsPool.Put(sessionStrings)

	sessions, err := internal.ParseSession(sessionsSlice)
	if err != nil {
		cc.logger.Printf("Error parsing sessions because of %v", err)
		return
	}
	var costs [][]string
	for _, session := range sessions {
		applicableTariffIndex := internal.SearchExactApplicableTariff(tariffs, 0, len(tariffs)-1, session)
		applicableTariffs := internal.GetApplicableTariffs(tariffs, applicableTariffIndex, session)
		sessionCost := internal.CostCalculator(applicableTariffs, session)
		costs = append(costs, sessionCost)
	}
	err = cc.WriteCosts(costs)
	if err != nil {
		cc.logger.Printf("Error writing costs to output file because of %v", err)
		return
	}
	return
}

func (cc *CostCalculator) WriteCosts(costs [][]string) error {
	var writer *csv.Writer
	cc.mu.Lock()
	defer cc.mu.Unlock()
	if _, err := os.Stat(cc.OutputFileLocation); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(cc.OutputFileLocation)
		if err != nil {
			cc.logger.Printf("Error creating output file because of %v", err)
			return err
		}
		defer func() {
			err = file.Close()
		}()
		writer, err = internal.CreateCSVWriterAndWriteHeader([]string{"session_id", "total_cost"}, file)
		if err != nil {
			cc.logger.Printf("Error writing output file header because of %v", err)
			return err
		}
	} else {
		file, err := os.OpenFile(cc.OutputFileLocation, os.O_APPEND|os.O_WRONLY, 0744)
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
	if err := writer.WriteAll(costs); err != nil {
		cc.logger.Printf("Error writing to output file because of %v", err)
		return err
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
