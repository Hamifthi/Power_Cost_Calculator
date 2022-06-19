package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
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
	// get and parse sync pool buffer size
	syncPoolSizeStr := envHandler.GetEnv("SYNC_POOL_SIZE")
	syncPoolSize, err := strconv.ParseInt(syncPoolSizeStr, 10, 64)
	if err != nil {
		logger.Printf("Error parsing sync pool size to int due to %v", err)
		syncPoolSize = 10240000
	}
	numberOfThreadsStr := envHandler.GetEnv("NUMBER_OF_THREADS")
	numberOfThreads, err := strconv.Atoi(numberOfThreadsStr)
	if err != nil {
		logger.Printf("Error parsing number of threads to int due to %v", err)
		numberOfThreads = 8
	}
	// get lines pool and strings pool
	// use sync pools to reduce pressure on garbage collector
	pools := entity.CreateSyncPools(syncPoolSize)
	// create cost calculator service and process inputs that creates output
	costCalculator := calculator.New(logger, pools, outputFileLocation, numberOfThreads)
	tariffs, err := costCalculator.ReadAndParseTariffs(tariffsFileLocation)
	if err != nil {
		logger.Fatalf("can't read tariffs file due to %v", err)
	}
	err = costCalculator.ReadAndProcessSessions(sessionFileLocation, tariffs)
	if err != nil {
		logger.Fatalf("can't process sessions file due to %v", err)
	}
	// now calculate inputs data
	inputTestOutputFileLocation := envHandler.GetEnv("INPUT_TEST_OUTPUT_FILE_LOCATION")
	inputTestTariffsFileLocation := envHandler.GetEnv("INPUT_TEST_TARIFFS_FILE_LOCATION")
	inputTestSessionFileLocation := envHandler.GetEnv("INPUT_TEST_SESSIONS_FILE_LOCATION")
	costCalculator.OutputFileLocation = inputTestOutputFileLocation
	inputTariffs, err := costCalculator.ReadAndParseTariffs(inputTestTariffsFileLocation)
	if err != nil {
		logger.Fatalf("can't read input tariffs file due to %v", err)
	}
	err = costCalculator.ReadAndProcessSessions(inputTestSessionFileLocation, inputTariffs)
	if err != nil {
		logger.Fatalf("can't process input sessions file due to %v", err)
	}
	// create fake data and process it
	start := time.Now()
	fakeTariffsFIleLocation := envHandler.GetEnv("FAKE_TARIFFS_FILE_LOCATION")
	fakeSessionsFIleLocation := envHandler.GetEnv("FAKE_SESSIONS_FILE_LOCATION")
	fakeCostsFIleLocation := envHandler.GetEnv("FAKE_COSTS_FILE_LOCATION")
	//get and parse number of fake sessions and tariffs
	numberOFFakeSessionsStr := envHandler.GetEnv("NUMBER_OF_SESSIONS")
	numberOFFakeSessions, err := strconv.Atoi(numberOFFakeSessionsStr)
	if err != nil {
		logger.Printf("Error parsing number of fake sessions to int due to %v", err)
		numberOFFakeSessions = 3000000
	}
	numberOFFakeTariffsStr := envHandler.GetEnv("NUMBER_OF_TARIFFS")
	numberOFFakeTariffs, err := strconv.Atoi(numberOFFakeTariffsStr)
	if err != nil {
		logger.Printf("Error parsing number of fake tariffs to int due to %v", err)
		numberOFFakeTariffs = 14000
	}
	dataGenerator := data.New(logger)
	err = dataGenerator.CreateSessions(fakeSessionsFIleLocation, numberOFFakeSessions)
	if err != nil {
		logger.Fatalf("Error creating fake sessions due to %v", err)
	}
	err = dataGenerator.CreateTariffs(fakeTariffsFIleLocation, numberOFFakeTariffs)
	if err != nil {
		logger.Fatalf("Error creating fake tariffs due to %v", err)
	}
	costCalculator.OutputFileLocation = fakeCostsFIleLocation
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
