package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg/calculator"
	"log"
	"os"
	"strconv"
)

func main() {
	// first initialize and read all environmental variables
	logger := log.New(os.Stdout, "CostCalculator ", log.LstdFlags)
	envHandler := internal.NewEnvHandler(logger)
	envHandler.InitializeEnv("./.env")
	// file locations that service will need
	outputFileLocation := envHandler.GetEnv("OUTPUT_FILE_LOCATION")
	tariffsFileLocation := envHandler.GetEnv("TARIFFS_FILE_LOCATION")
	sessionFileLocation := envHandler.GetEnv("SESSIONS_FILE_LOCATION")
	// get and parse sync pool buffer size
	syncPoolSizeStr := envHandler.GetEnv("SYNC_POOL_SIZE")
	syncPoolSize, err := strconv.ParseInt(syncPoolSizeStr, 10, 64)
	if err != nil {
		logger.Fatalf("Error parsing sync pool size to int to %v", err)
	}
	// get lines pool and strings pool
	linesPool, stringsPool := calculator.CreateSyncPools(syncPoolSize)
	// create cost calculator service and process inputs that creates output
	costCalculator := calculator.New(logger, linesPool, stringsPool, outputFileLocation)
	tariffs, err := costCalculator.ReadAndParseTariffs(tariffsFileLocation)
	if err != nil {
		logger.Fatalf("can't read tariffs file due to %v", err)
	}
	err = costCalculator.ReadAndProcessSessions(sessionFileLocation, tariffs)
	if err != nil {
		logger.Fatalf("can't process sessions file due to %v", err)
	}
}
