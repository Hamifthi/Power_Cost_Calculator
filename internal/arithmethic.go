package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"math"
)

func CheckTimeOverlap(tariff entity.Tariff, session entity.Session) bool {
	if tariff.Start.Before(session.End) && tariff.End.After(session.Start) {
		return true
	}
	return false
}

func TruncateFloat(number float64) float64 {
	return math.Trunc(number*1000) / 1000
}

func CostCalculator(tariffs []entity.Tariff, sessions []entity.Session) ([]entity.Cost, error) {
	uniqueCosts := make(map[string]float64)
	var costs []entity.Cost
	for _, session := range sessions {
		for _, tariff := range tariffs {
			if CheckTimeOverlap(tariff, session) {
				tariffDuration := tariff.End.Sub(tariff.Start)
				uniqueCosts[session.ID] = session.Energy*tariff.EnergyFee +
					tariff.ParkingFee*tariffDuration.Hours()
			}
		}
	}
	for id, totalCost := range uniqueCosts {
		uniqueCosts[id] = TruncateFloat(totalCost * 1.15)
		cost := entity.Cost{SessionID: id, TotalCost: uniqueCosts[id]}
		costs = append(costs, cost)
	}
	return costs, nil
}
