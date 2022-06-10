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
