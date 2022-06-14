package calculator

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"reflect"
	"sync"
	"testing"
	"time"
)

func initializeCalculatorService() *CostCalculator {
	logger := log.New(os.Stdout, "CostCalculator_Test ", log.LstdFlags)
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 200*1024)
		return lines
	}}
	stringsPool := sync.Pool{New: func() interface{} {
		strs := ""
		return strs
	}}
	CostCalculator := New(logger, &linesPool, &stringsPool, "./costs_test.csv")
	return CostCalculator
}

func TestReadAndParseTariffsSuccessfully(t *testing.T) {
	calculator := initializeCalculatorService()
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	_ = internal.InitializeEnv("../../test.env")
	tariffFilePath, _ := internal.GetEnv("TARIFFS_FILE_LOCATION", logger)
	tariffs, err := calculator.ReadAndParseTariffs(tariffFilePath)
	assert.Nil(t, err)
	assert.NotNil(t, tariffs)
	tariff := tariffs[0]
	assert.IsType(t, time.Time{}, tariff.Start)
	assert.IsType(t, time.Time{}, tariff.End)
	assert.Equal(t, tariff.EnergyFee, 0.5)
	assert.Equal(t, tariff.ParkingFee, 0.5)
}

func TestReadAndProcessSessionsInvalidLocationFails(t *testing.T) {
	calculator := initializeCalculatorService()
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	_ = internal.InitializeEnv("../../test.env")
	tariffFilePath, _ := internal.GetEnv("TARIFFS_FILE_LOCATION", logger)
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	filePath := "data/sessions.csv"
	err := calculator.ReadAndProcessSessions(filePath, tariffs)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "open %s no such file or directory", filePath)
}

func TestReadAndProcessSessionsAlreadyExistedCostsFileSuccessfully(t *testing.T) {
	calculator := initializeCalculatorService()
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	_ = internal.InitializeEnv("../../test.env")
	outputFileLocation, _ := internal.GetEnv("OUTPUT_FILE_LOCATION", logger)
	_, _ = os.Create(outputFileLocation)
	tariffFilePath, _ := internal.GetEnv("TARIFFS_FILE_LOCATION", logger)
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	sessionsFilePath, _ := internal.GetEnv("SESSIONS_FILE_LOCATION", logger)
	err := calculator.ReadAndProcessSessions(sessionsFilePath, tariffs)
	assert.Nil(t, err)
	_ = os.Remove("./costs_test.csv")
}

func TestReadAndProcessSessionsSuccessfully(t *testing.T) {
	calculator := initializeCalculatorService()
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	_ = internal.InitializeEnv("../../test.env")
	tariffFilePath, _ := internal.GetEnv("TARIFFS_FILE_LOCATION", logger)
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	sessionsFilePath, _ := internal.GetEnv("SESSIONS_FILE_LOCATION", logger)
	err := calculator.ReadAndProcessSessions(sessionsFilePath, tariffs)
	assert.Nil(t, err)
	actualCostFilePath, _ := internal.GetEnv("ACTUAL_COSTS_FILE_LOCATION", logger)
	actualCosts, _ := internal.ReadFile(actualCostFilePath)
	outputFileLocation, _ := internal.GetEnv("OUTPUT_FILE_LOCATION", logger)
	expectedCosts, _ := internal.ReadFile(outputFileLocation)
	assert.True(t, reflect.DeepEqual(actualCosts, expectedCosts))
	_ = os.Remove(outputFileLocation)
}
