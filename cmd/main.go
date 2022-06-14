package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg/calculator"
	"log"
	"os"
	"strconv"
)

func main() {
	logger := log.New(os.Stdout, "CostCalculator ", log.LstdFlags)
	// get and parse sync pool buffer size
	syncPoolSizeStr, err := internal.GetEnv("SYNC_POOL_SIZE", logger)
	syncPoolSize, _ := strconv.ParseInt(syncPoolSizeStr, 10, 64)
	if err != nil {
		logger.Fatalf("Error parsing sync pool size to int to %v", err)
	}
	// get lines pool and strings pool
	linesPool, stringsPool := calculator.CreateSyncPools(syncPoolSize)
	// output file location that service will create
	outputFileLocation, _ := internal.GetEnv("OUTPUT_FILE_LOCATION", logger)
	// create cost calculator service and process inputs that creates output
	costCalculator := calculator.New(logger, linesPool, stringsPool, outputFileLocation)
	tariffsFileLocation, _ := internal.GetEnv("TARIFFS_FILE_LOCATION", logger)
	tariffs, err := costCalculator.ReadAndParseTariffs(tariffsFileLocation)
	if err != nil {
		logger.Fatalf("can't read tariffs file due to %v", err)
	}
	sessionFileLocation, _ := internal.GetEnv("SESSIONS_FILE_LOCATION", logger)
	err = costCalculator.ReadAndProcessSessions(sessionFileLocation, tariffs)
	if err != nil {
		logger.Fatalf("can't process sessions file due to %v", err)
	}
}
