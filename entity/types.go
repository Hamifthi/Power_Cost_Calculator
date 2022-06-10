package entity

import "time"

type Tariff struct {
	Start      time.Time
	End        time.Time
	EnergyFee  float64
	ParkingFee float64
}

type Session struct {
	ID     string
	Start  time.Time
	End    time.Time
	Energy float64
}

type Cost struct {
	SessionID string
	TotalCost float64
}
