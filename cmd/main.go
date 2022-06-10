package main

import (
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/internal"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg"
	"log"
	"os"
)

func main() {
	logger := log.New(os.Stdout, "CostCalculator ", log.LstdFlags)
	tariffs, sessions, err := pkg.ReadAndParseFiles("./data/tariffs.csv", "./data/sessions.csv")
	if err != nil {
		logger.Println(err)
	}
	costs, err := pkg.CostCalculator(tariffs, sessions)
	if err != nil {
		logger.Println(err)
	}
	err = internal.WriteFile("./data/calculated_costs.csv", costs)
	if err != nil {
		logger.Println(err)
	}
}
