package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
