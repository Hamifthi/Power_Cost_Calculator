package calculator

import "github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"

type CostCalculatorInterface interface {
	ReadAndProcessSessions(sessionFilePath string, tariffs []entity.Tariff) error
	processChunk(chunk []byte, tariffs []entity.Tariff)
	WriteCosts(costs []entity.Cost) error
	ReadAndParseTariffs(filePath string) ([]entity.Tariff, error)
}
