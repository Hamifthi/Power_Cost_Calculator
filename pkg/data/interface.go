package data

type FakeDataGeneratorInterface interface {
	CreateSessions(generatedFilePath string, numberOfSessions int) error
	CreateTariffs(generatedFilePath string, numberOfTariffs int) error
}
