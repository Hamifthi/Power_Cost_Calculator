package internal

import (
	"encoding/csv"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
	"time"
)

func createStartEndTimeFromString(startString, endString string) (time.Time, time.Time) {
	start, _ := time.Parse(time.RFC3339, startString)
	end, _ := time.Parse(time.RFC3339, endString)
	return start, end
}

func TestReadFileInvalidFilePathFails(t *testing.T) {
	filePath := "./tariffs.csv"
	stringSlice, err := ReadFile(filePath)
	assert.Nil(t, stringSlice)
	assert.NotNil(t, err)
	assert.EqualErrorf(t, err, err.Error(), "open %s no such file or directory", filePath)
}

func TestReadFileSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := NewEnvHandler(logger)
	envHandler.InitializeEnv("../test.env")
	tariffFileLocation := envHandler.GetEnv("INTERNAL_TARIFFS_FILE_LOCATION")
	stringSlice, err := ReadFile(tariffFileLocation)
	assert.Nil(t, err)
	assert.NotNil(t, stringSlice)
	assert.IsType(t, []string{}, stringSlice)
}

func TestParseTariffInvalidStartFails(t *testing.T) {
	tariffLines := []string{"2000,2020-06-03T00:00:00+02:00,0.5,0.5"}
	tariffs, err := ParseTariff(tariffLines)
	assert.Nil(t, tariffs)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Tariffs converting start time string")
}

func TestParseTariffInvalidEndFails(t *testing.T) {
	tariffLines := []string{"2020-06-03T00:00:00+02:00,2000,0.5,0.5"}
	tariffs, err := ParseTariff(tariffLines)
	assert.Nil(t, tariffs)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Tariffs converting end time string")
}

func TestParseTariffInvalidEnergyFeeFails(t *testing.T) {
	tariffLines := []string{"2020-06-03T00:00:00+02:00,2020-06-03T10:10:00+02:00,float,0.5"}
	tariffs, err := ParseTariff(tariffLines)
	assert.Nil(t, tariffs)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Tariffs converting energy fee")
}

func TestParseTariffInvalidParkingFeeFails(t *testing.T) {
	tariffLines := []string{"2020-06-03T00:00:00+02:00,2020-06-03T10:10:00+02:00,0.5,float"}
	tariffs, err := ParseTariff(tariffLines)
	assert.Nil(t, tariffs)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Tariffs converting parking fee")
}

func TestParseTariffSuccessfully(t *testing.T) {
	tariffLines := []string{"2020-06-03T00:00:00+02:00,2020-06-03T10:10:00+02:00,0.5,0.5"}
	tariffs, err := ParseTariff(tariffLines)
	assert.Nil(t, err)
	start, end := createStartEndTimeFromString("2020-06-03T00:00:00+02:00", "2020-06-03T10:10:00+02:00")
	assert.Equal(t, tariffs, []entity.Tariff{{Start: start, End: end, EnergyFee: 0.5, ParkingFee: 0.5}})
}

func TestParseSessionInvalidStartFails(t *testing.T) {
	sessionLines := []string{"a949d681-e12b-4d93-a3e5-e2e777e68f12,2000,2020-06-03T00:00:00+02:00,3"}
	sessions, err := ParseSession(sessionLines)
	assert.Nil(t, sessions)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Sessions converting start time string")
}

func TestParseSessionInvalidEndFails(t *testing.T) {
	sessionLines := []string{"a949d681-e12b-4d93-a3e5-e2e777e68f12,2020-06-03T00:00:00+02:00,2000,3"}
	sessions, err := ParseSession(sessionLines)
	assert.Nil(t, sessions)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Sessions converting end time string")
}

func TestParseSessionInvalidEnergyFails(t *testing.T) {
	sessionLines := []string{"a949d681-e12b-4d93-a3e5-e2e777e68f12,2020-06-03T10:00:00+02:00," +
		"2020-06-03T10:26:00+02:00,energy"}
	sessions, err := ParseSession(sessionLines)
	assert.Nil(t, sessions)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, err.Error(), "Error while parsing Sessions converting energy")
}

func TestParseSessionSuccessfully(t *testing.T) {
	sessionLines := []string{"a949d681-e12b-4d93-a3e5-e2e777e68f12,2020-06-03T10:00:00+02:00," +
		"2020-06-03T10:26:00+02:00,3"}
	sessions, err := ParseSession(sessionLines)
	assert.Nil(t, err)
	start, end := createStartEndTimeFromString("2020-06-03T10:00:00+02:00", "2020-06-03T10:26:00+02:00")
	assert.Equal(t, sessions, []entity.Session{{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12",
		Start: start, End: end, Energy: 3.0}})
}

func TestCheckTimeOverlapWithoutOverlap(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T14:49:00+02:00",
		"2020-06-03T14:55:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	assert.False(t, checkTimeOverlap(tariff, session))
}

func TestCheckTimeOverlapWithOverlapTariffBeforeSession(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T09:49:00+02:00",
		"2020-06-03T14:55:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	assert.True(t, checkTimeOverlap(tariff, session))
}

func TestCheckTimeOverlapWithOverlapSessionBeforeTariff(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-02T23:49:00+02:00",
		"2020-06-03T05:55:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	assert.True(t, checkTimeOverlap(tariff, session))
}

func TestCheckTimeOverlapWithOverlapSessionEqualTariff(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	assert.True(t, checkTimeOverlap(tariff, session))
}

func TestCalculateTimeOverlapWithOverlapTariffBeforeSession(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T09:49:00+02:00",
		"2020-06-03T14:55:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	timeOverlap := calculateTimeOverlap(tariff, session)
	assert.Equal(t, timeOverlap, 0.35)
}

func TestCalculateTimeOverlapWithOverlapSessionBeforeTariff(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-02T23:00:00+02:00",
		"2020-06-03T05:30:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	timeOverlap := calculateTimeOverlap(tariff, session)
	assert.Equal(t, timeOverlap, 5.5)
}

func TestCalculateTimeOverlapWithOverlapSessionLessThanTariff(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T06:00:00+02:00",
		"2020-06-03T09:20:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	timeOverlap := TruncateFloat(calculateTimeOverlap(tariff, session))
	assert.Equal(t, timeOverlap, 3.333)
}

func TestCalculateTimeOverlapWithOverlapSessionEqualTariff(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	timeOverlap := TruncateFloat(calculateTimeOverlap(tariff, session))
	assert.Equal(t, timeOverlap, 10.166)
}

func TestCostCalculator(t *testing.T) {
	tariffStart, tariffEnd := createStartEndTimeFromString("2020-06-03T00:00:00+02:00",
		"2020-06-03T10:10:00+02:00")
	tariff := entity.Tariff{Start: tariffStart, End: tariffEnd, EnergyFee: 0.5, ParkingFee: 0.5}
	sessionStart, sessionEnd := createStartEndTimeFromString("2020-06-03T10:00:00+02:00",
		"2020-06-03T10:26:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	costs := CostCalculator([]entity.Tariff{tariff}, []entity.Session{session})
	assert.Equal(t, costs[0], entity.Cost{SessionID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", TotalCost: 0.759})
}

func TestGetEnvSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := NewEnvHandler(logger)
	envHandler.InitializeEnv("../test.env")
	value := envHandler.GetEnv("TEST")
	assert.Equal(t, value, "test")
}

func TestCreateCSVWriterAndWriteHeader(t *testing.T) {
	fileName := "./something"
	file, _ := os.Create(fileName)
	csvWriter, err := CreateCSVWriterAndWriteHeader([]string{"something"}, file)
	assert.Nil(t, err)
	assert.IsType(t, &csv.Writer{}, csvWriter)
	_ = os.Remove(fileName)
}
