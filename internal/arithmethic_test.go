package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

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
	cost := CostCalculator([]entity.Tariff{tariff}, session)
	assert.Equal(t, cost, []string{"a949d681-e12b-4d93-a3e5-e2e777e68f12", "0.759"})
}

func TestSearchExactApplicableTariffIndexNotFound(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := ReadTestEnv(envPath, logger)
	tariffLines, _ := ReadFile(envHandler.GetEnv("TEST_TARIFF_FILE_LOCATION"))
	tariffs, _ := ParseTariff(tariffLines[1:])
	sessionStart, sessionEnd := createStartEndTimeFromString("2023-06-03T10:00:00+02:00",
		"2023-06-03T10:26:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	value := SearchExactApplicableTariff(tariffs, 0, len(tariffs)-1, session)
	assert.Equal(t, value, -1)
}

func TestSearchExactApplicableTariffSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := ReadTestEnv(envPath, logger)
	tariffLines, _ := ReadFile(envHandler.GetEnv("TEST_TARIFF_FILE_LOCATION"))
	tariffs, _ := ParseTariff(tariffLines[1:])
	sessionStart, sessionEnd := createStartEndTimeFromString("2017-06-01T10:00:00+02:00",
		"2017-06-01T10:26:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	value := SearchExactApplicableTariff(tariffs, 0, len(tariffs)-1, session)
	assert.Equal(t, value, 12)
}

func TestGetApplicableTariffsSuccessfully(t *testing.T) {
	logger := log.New(os.Stdout, "Test ", log.LstdFlags)
	envHandler := ReadTestEnv(envPath, logger)
	tariffLines, _ := ReadFile(envHandler.GetEnv("TEST_TARIFF_FILE_LOCATION"))
	tariffs, _ := ParseTariff(tariffLines[1:])
	sessionStart, sessionEnd := createStartEndTimeFromString("2017-06-01T00:00:00+02:00",
		"2017-06-01T23:26:00+02:00")
	session := entity.Session{ID: "a949d681-e12b-4d93-a3e5-e2e777e68f12", Start: sessionStart,
		End: sessionEnd, Energy: 3}
	value := SearchExactApplicableTariff(tariffs, 0, len(tariffs)-1, session)
	applicableTariffs := GetApplicableTariffs(tariffs, value, session)
	assert.Equal(t, len(applicableTariffs), 9)
}

func TestDateEqual(t *testing.T) {
	firstDate, secondDate := createStartEndTimeFromString("2017-06-01T00:00:00+02:00",
		"2017-06-01T23:26:00+02:00")
	assert.True(t, dateEqual(firstDate, secondDate))
}
