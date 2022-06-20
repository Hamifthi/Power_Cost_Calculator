package calculator

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

var testEnvPath string

func init() {
	currentDir, _ := os.Getwd()
	rootDir := filepath.Dir(filepath.Dir(currentDir))
	testEnvPath = rootDir + "/test.env"
}

func initializeCalculatorService(logger *log.Logger, envHandler *internal.EnvHandler) *CostCalculator {
	syncPoolSizeStr := envHandler.GetEnv("TEST_SYNC_POOL_SIZE")
	syncPoolSize, _ := strconv.ParseInt(syncPoolSizeStr, 10, 64)
	pools := entity.CreateSyncPools(syncPoolSize)
	CostCalculator := New(logger, pools, envHandler.GetEnv("OUTPUT_FILE_LOCATION"), 5)
	return CostCalculator
}

func TestReadAndParseTariffsSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Calculator_Test ", log.LstdFlags)
	envHandler := internal.ReadTestEnv(testEnvPath, logger)
	calculator := initializeCalculatorService(logger, envHandler)
	tariffFilePath := envHandler.GetEnv("TARIFFS_FILE_LOCATION")
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
	logger := log.New(os.Stdout, "Calculator_Test ", log.LstdFlags)
	envHandler := internal.ReadTestEnv(testEnvPath, logger)
	calculator := initializeCalculatorService(logger, envHandler)
	tariffFilePath := envHandler.GetEnv("TARIFFS_FILE_LOCATION")
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	filePath := "data/sessions.csv"
	err := calculator.ReadAndProcessSessions(filePath, tariffs)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "open %s no such file or directory", filePath)
}

func TestReadAndProcessSessionsAlreadyExistedCostsFileSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Calculator_Test ", log.LstdFlags)
	envHandler := internal.ReadTestEnv(testEnvPath, logger)
	calculator := initializeCalculatorService(logger, envHandler)
	outputFileLocation := envHandler.GetEnv("OUTPUT_FILE_LOCATION")
	_, _ = os.Create(outputFileLocation)
	tariffFilePath := envHandler.GetEnv("TARIFFS_FILE_LOCATION")
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	sessionsFilePath := envHandler.GetEnv("SESSIONS_FILE_LOCATION")
	err := calculator.ReadAndProcessSessions(sessionsFilePath, tariffs)
	assert.Nil(t, err)
	_ = os.Remove(outputFileLocation)
}

func TestReadAndProcessSessionsSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Calculator_Test ", log.LstdFlags)
	envHandler := internal.ReadTestEnv(testEnvPath, logger)
	calculator := initializeCalculatorService(logger, envHandler)
	tariffFilePath := envHandler.GetEnv("TARIFFS_FILE_LOCATION")
	tariffs, _ := calculator.ReadAndParseTariffs(tariffFilePath)
	sessionsFilePath := envHandler.GetEnv("SESSIONS_FILE_LOCATION")
	err := calculator.ReadAndProcessSessions(sessionsFilePath, tariffs)
	assert.Nil(t, err)
	actualCostFilePath := envHandler.GetEnv("ACTUAL_COSTS_FILE_LOCATION")
	actualCosts, _ := internal.ReadFile(actualCostFilePath)
	outputFileLocation := envHandler.GetEnv("OUTPUT_FILE_LOCATION")
	expectedCosts, _ := internal.ReadFile(outputFileLocation)
	assert.ElementsMatch(t, actualCosts, expectedCosts)
	_ = os.Remove(outputFileLocation)
}
