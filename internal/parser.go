package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/pkg/errors"
	"strconv"
	"strings"
)

func ParseTariff(lines []string) ([]entity.Tariff, error) {
	var tariffs []entity.Tariff
	for i := 0; i < len(lines); i++ {
		stringSlice := strings.Split(lines[i], ",")
		energyFee, err := strconv.ParseFloat(stringSlice[2], 32)
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting energy fee")
		}
		parkingFee, err := strconv.ParseFloat(stringSlice[3], 32)
		if err != nil {
			return tariffs, errors.Wrap(err, "Error while parsing Tariffs converting parking fee")
		}
		tariff := entity.Tariff{Start: stringSlice[0], End: stringSlice[1],
			EnergyFee: float32(energyFee), ParkingFee: float32(parkingFee)}
		tariffs = append(tariffs, tariff)
	}
	return tariffs, nil
}

func ParseSession(lines []string) ([]entity.Session, error) {
	var sessions []entity.Session
	for i := 0; i < len(lines); i++ {
		stringSlice := strings.Split(lines[i], ",")
		energy, err := strconv.ParseFloat(stringSlice[3], 32)
		if err != nil {
			return sessions, errors.Wrap(err, "Error while parsing Sessions converting energy")
		}
		session := entity.Session{ID: stringSlice[0], Start: stringSlice[1],
			End: stringSlice[2], Energy: float32(energy)}
		sessions = append(sessions, session)
	}
	return sessions, nil
}
