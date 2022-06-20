package internal

import (
	"fmt"
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

// use this function to compare date without considering timezone because truncate function apply
// timezone in the calculation
func dateEqual(firstDate, secondDate time.Time) bool {
	y1, m1, d1 := firstDate.Date()
	y2, m2, d2 := secondDate.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

func SearchExactApplicableTariff(tariffs []entity.Tariff, left, right int, session entity.Session) int {
	if right >= left {
		mid := left + (right-left)/2
		tariff := tariffs[mid]
		if dateEqual(tariff.Start, session.Start) || dateEqual(tariff.End, session.Start) ||
			dateEqual(tariff.Start, session.End) || dateEqual(tariff.End, session.End) {
			return mid
		}
		if session.End.Before(tariff.Start) {
			return SearchExactApplicableTariff(tariffs, left, mid-1, session)
		} else if session.Start.After(tariff.End) {
			return SearchExactApplicableTariff(tariffs, mid+1, right, session)
		}
	}
	return -1
}

func GetApplicableTariffs(tariffs []entity.Tariff, tariffIndex int, session entity.Session) []entity.Tariff {
	var applicableTariffs []entity.Tariff
	sessionOneDayBefore := session.Start.Add(-time.Hour * 24)
	sessionOneDayAfter := session.End.Add(time.Hour * 24)
	// first get applicable tariffs that are before the session
	index := tariffIndex
	for {
		if index < 0 {
			break
		}

		if dateEqual(tariffs[index].Start, session.Start) ||
			dateEqual(tariffs[index].End, session.Start) ||
			dateEqual(tariffs[index].Start, sessionOneDayBefore) ||
			dateEqual(tariffs[index].End, sessionOneDayBefore) {
			if checkTimeOverlap(tariffs[index], session) {
				applicableTariffs = append(applicableTariffs, tariffs[index])
			}
		} else {
			break
		}
		index -= 1
	}
	// then get applicable tariffs that are after the session
	index = tariffIndex
	for {
		index += 1
		if index >= len(tariffs) {
			break
		}
		if dateEqual(tariffs[index].Start, session.End) ||
			dateEqual(tariffs[index].End, session.End) ||
			dateEqual(tariffs[index].Start, sessionOneDayAfter) ||
			dateEqual(tariffs[index].End, sessionOneDayAfter) {
			if checkTimeOverlap(tariffs[index], session) {
				applicableTariffs = append(applicableTariffs, tariffs[index])
			}
		} else {
			break
		}
	}
	return applicableTariffs
}

func CostCalculator(tariffs []entity.Tariff, session entity.Session) []string {
	var cost float64
	sessionDuration := session.End.Sub(session.Start).Hours()
	for _, tariff := range tariffs {
		duration := calculateTimeOverlap(tariff, session)
		cost += session.Energy*(duration/sessionDuration)*tariff.EnergyFee +
			tariff.ParkingFee*duration
	}
	cost = TruncateFloat(cost * 1.15)
	formattedCost := []string{session.ID, fmt.Sprintf("%.3f", cost)}
	return formattedCost
}
