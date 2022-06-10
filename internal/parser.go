package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

func ParseTariff(lines []string) ([]entity.Tariff, error) {
	var tariffs []entity.Tariff
	for i := 0; i < len(lines); i++ {
		stringSlice := strings.Split(lines[i], ",")
		start, err := time.Parse(time.RFC3339, stringSlice[0])
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting start time string")
		}
		end, err := time.Parse(time.RFC3339, stringSlice[1])
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting end time string")
		}
		energyFee, err := strconv.ParseFloat(stringSlice[2], 64)
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting energy fee")
		}
		parkingFee, err := strconv.ParseFloat(stringSlice[3], 64)
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting parking fee")
		}
		tariff := entity.Tariff{Start: start, End: end, EnergyFee: energyFee, ParkingFee: parkingFee}
		tariffs = append(tariffs, tariff)
	}
	return tariffs, nil
}

func ParseSession(lines []string) ([]entity.Session, error) {
	var sessions []entity.Session
	for i := 0; i < len(lines); i++ {
		stringSlice := strings.Split(lines[i], ",")
		start, err := time.Parse(time.RFC3339, stringSlice[1])
		if err != nil {
			return sessions, errors.Wrap(err, "Error while parsing Sessions converting start time string")
		}
		end, err := time.Parse(time.RFC3339, stringSlice[1])
		if err != nil {
			return sessions, errors.Wrap(err, "Error while parsing Sessions converting end time string")
		}
		energy, err := strconv.ParseFloat(stringSlice[3], 64)
		if err != nil {
			return sessions, errors.Wrap(err, "Error while parsing Sessions converting energy")
		}
		session := entity.Session{ID: stringSlice[0], Start: start, End: end, Energy: energy}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
