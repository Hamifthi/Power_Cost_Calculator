package pkg

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/pkg/errors"
)

func ReadAndParseFiles(tariffFilePath, sessionFilePath string) ([]entity.Tariff, []entity.Session, error) {
	var tariffs []entity.Tariff
	var sessions []entity.Session
	tariffStrings, err := internal.ReadFile(tariffFilePath)
	if err != nil {
		return tariffs, sessions, errors.Wrap(err, "Error reading tariffs file")
	}
	sessionStrings, err := internal.ReadFile(sessionFilePath)
	if err != nil {
		return tariffs, sessions, errors.Wrap(err, "Error reading sessions file")
	}
	tariffs, err = internal.ParseTariff(tariffStrings[1:])
	if err != nil {
		return tariffs, sessions, err
	}
	sessions, err = internal.ParseSession(sessionStrings[1:])
	if err != nil {
		return tariffs, sessions, err
	}
	return tariffs, sessions, nil
}

func CostCalculator(tariffs []entity.Tariff, sessions []entity.Session) ([]entity.Cost, error) {
	uniqueCosts := make(map[string]float64)
	var costs []entity.Cost
	for _, session := range sessions {
		for _, tariff := range tariffs {
			if internal.CheckTimeOverlap(tariff, session) {
				tariffDuration := tariff.End.Sub(tariff.Start)
				uniqueCosts[session.ID] = session.Energy*tariff.EnergyFee +
					tariff.ParkingFee*tariffDuration.Hours()
			}
		}
	}
	for id, totalCost := range uniqueCosts {
		uniqueCosts[id] = internal.TruncateFloat(totalCost * 1.15)
		cost := entity.Cost{SessionID: id, TotalCost: uniqueCosts[id]}
		costs = append(costs, cost)
	}
	return costs, nil
}
