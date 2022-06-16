package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg/calculator"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg/data"
	"log"
	"os"
	"strconv"
	"time"
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
	fakeTariffsFIleLocation := envHandler.GetEnv("FAKE_TARIFFS_FILE_LOCATION")
	fakeSessionsFIleLocation := envHandler.GetEnv("FAKE_SESSIONS_FILE_LOCATION")
	// get and parse number of fake sessions and tariffs
	numberOFFakeSessionsStr := envHandler.GetEnv("NUMBER_OF_SESSIONS")
	numberOFFakeSessions, err := strconv.Atoi(numberOFFakeSessionsStr)
	if err != nil {
		logger.Fatalf("Error parsing number of fake sessions to int due to %v", err)
	}
	numberOFFakeTariffsStr := envHandler.GetEnv("NUMBER_OF_TARIFFS")
	numberOFFakeTariffs, err := strconv.Atoi(numberOFFakeTariffsStr)
	if err != nil {
		logger.Fatalf("Error parsing number of fake tariffs to int due to %v", err)
	}
	// get and parse sync pool buffer size
	syncPoolSizeStr := envHandler.GetEnv("SYNC_POOL_SIZE")
	syncPoolSize, err := strconv.ParseInt(syncPoolSizeStr, 10, 64)
	if err != nil {
		logger.Fatalf("Error parsing sync pool size to int due to %v", err)
	}
	// get lines pool and strings pool
	// use sync pools to reduce pressure on garbage collector
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
	// create fake data and process it
	start := time.Now()
	dataGenerator := data.New(logger)
	err = dataGenerator.CreateSessions(fakeSessionsFIleLocation, numberOFFakeSessions)
	if err != nil {
		logger.Fatalf("Error creating fake sessions due to %v", err)
	}
	err = dataGenerator.CreateTariffs(fakeTariffsFIleLocation, numberOFFakeTariffs)
	if err != nil {
		logger.Fatalf("Error creating fake tariffs due to %v", err)
	}
	tariffs, err = costCalculator.ReadAndParseTariffs(fakeTariffsFIleLocation)
	if err != nil {
		logger.Fatalf("can't read tariffs file due to %v", err)
	}
	err = costCalculator.ReadAndProcessSessions(fakeSessionsFIleLocation, tariffs)
	if err != nil {
		logger.Fatalf("can't process sessions file due to %v", err)
	}
	logger.Printf("time to create and process fake data is %v", time.Since(start))
}
