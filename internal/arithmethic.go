package internal

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"math"
	"time"
)

// check time overlap works with this rule (StartA <= EndB) and (EndA >= StartB)
// keep in mind I use tariff as time A and session as time B
func checkTimeOverlap(tariff entity.Tariff, session entity.Session) bool {
	if tariff.Start.Before(session.End) && tariff.End.After(session.Start) {
		return true
	}
	return false
}

// calculate time overlaps first calculated max of starts and then min of ends.
// Then it subtracts the min end from max start to get time overlap
func calculateTimeOverlap(tariff entity.Tariff, session entity.Session) float64 {
	var duration time.Duration
	var start time.Time
	var end time.Time
	if tariff.Start.After(session.Start) {
		start = tariff.Start
	} else if tariff.Start.Before(session.Start) {
		start = session.Start
	} else {
		start = tariff.Start
	}
	if tariff.End.Before(session.End) {
		end = tariff.End
	} else if tariff.End.After(session.End) {
		end = session.End
	} else {
		end = tariff.End
	}
	duration = end.Sub(start)
	return duration.Hours()
}

func TruncateFloat(number float64) float64 {
	return math.Trunc(number*1000) / 1000
}

func CostCalculator(tariffs []entity.Tariff, sessions []entity.Session) []entity.Cost {
	uniqueCosts := make(map[string]float64)
	var costs []entity.Cost
	for _, session := range sessions {
		sessionDuration := session.End.Sub(session.Start).Hours()
		for _, tariff := range tariffs {
			if checkTimeOverlap(tariff, session) {
				duration := calculateTimeOverlap(tariff, session)
				uniqueCosts[session.ID] += session.Energy*(duration/sessionDuration)*tariff.EnergyFee +
					tariff.ParkingFee*duration
			}
		}
	}
	for id, totalCost := range uniqueCosts {
		uniqueCosts[id] = TruncateFloat(totalCost * 1.15)
		cost := entity.Cost{SessionID: id, TotalCost: uniqueCosts[id]}
		costs = append(costs, cost)
	}
	return costs
}
