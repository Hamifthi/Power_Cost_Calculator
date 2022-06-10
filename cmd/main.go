package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "CostCalculator ", log.LstdFlags)
	tariffStrings, err := internal.ReadFile("./data/tariffs.csv")
	if err != nil {
		logger.Println(err)
	}
	tariffs, err := internal.ParseTariff(tariffStrings[1:])
	if err != nil {
		logger.Println(err)
	}
	err = pkg.ReadAndProcessEfficiently("./data/sessions.csv", tariffs)
	if err != nil {
		logger.Println(err)
	}
}
