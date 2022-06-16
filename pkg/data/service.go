package data

import (
	"encoding/csv"
	"fmt"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"log"
	"math/rand"
	"os"
	"time"
)

type FakeDataGenerator struct {
	logger *log.Logger
}

func New(logger *log.Logger) *FakeDataGenerator {
	return &FakeDataGenerator{logger: logger}
}

func createRandomID() string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"
	b := make([]byte, 32)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	dashIndexes := []int{8, 13, 18, 23}
	for _, i := range dashIndexes {
		b = append(b[:i+1], b[i:]...)
		b[i] = '-'
	}
	return string(b)
}

func writeToCsvFile(header []string, fileLocation string, data [][]string) error {
	var writer *csv.Writer
	file, err := os.Create(fileLocation)
	if err != nil {
		return err
	}
	defer func() {
		err = file.Close()
	}()
	writer = csv.NewWriter(file)
	writer, err = internal.CreateCSVWriterAndWriteHeader(header, file)
	if err != nil {
		return err
	}
	defer writer.Flush()
	err = writer.WriteAll(data)
	if err != nil {
		return err
	}
	return nil
}

func (dg *FakeDataGenerator) CreateSessions(generatedFilePath string, numberOfSessions int) error {
	var sessions [][]string
	sessionStart, _ := time.Parse(time.RFC3339, "2020-06-04T05:41:00+02:00")
	var sessionEnd time.Time
	for i := 0; i < numberOfSessions; i++ {
		sessionEnd = sessionStart.Add(time.Duration(rand.Intn(120) * 10e+9))
		energy := internal.TruncateFloat(rand.Float64() * 5)
		sessions = append(sessions, []string{createRandomID(), sessionStart.Format(time.RFC3339),
			sessionEnd.Format(time.RFC3339), fmt.Sprintf("%.2f", energy)})
		sessionStart = sessionEnd
	}
	err := writeToCsvFile([]string{"ID", "dt_start", "dt_end", "energy"}, generatedFilePath, sessions)
	if err != nil {
		dg.logger.Printf("Error writing sessions to csv output file because of %v", err)
		return err
	}
	return nil
}

func (dg *FakeDataGenerator) CreateTariffs(generatedFilePath string, numberOfTariffs int) error {
	var tariffs [][]string
	tariffStart, _ := time.Parse(time.RFC3339, "2020-06-04T22:00:00+02:00")
	var tariffEnd time.Time
	for i := 0; i < numberOfTariffs; i++ {
		tariffEnd = tariffStart.Add(time.Duration(rand.Intn(21600) * 10e+9))
		energyFee := rand.Float64()
		parkingFee := rand.Float64()
		tariffs = append(tariffs, []string{tariffStart.Format(time.RFC3339), tariffEnd.Format(time.RFC3339),
			fmt.Sprintf("%.1f", energyFee), fmt.Sprintf("%.1f", parkingFee)})
		tariffStart = tariffEnd
	}
	err := writeToCsvFile([]string{"dt_start", "dt_end", "energy_fee", "parking_fee"}, generatedFilePath, tariffs)
	if err != nil {
		dg.logger.Printf("Error writing tariffs to csv output file because of %v", err)
		return err
	}
	return nil
}
