package entity

type Tariff struct {
	Start      string
	End        string
	EnergyFee  float32
	ParkingFee float32
}

type Session struct {
	ID     string
	Start  string
	End    string
	Energy float32
}

type Cost struct {
	SessionID string
	TotalCost float32
}
