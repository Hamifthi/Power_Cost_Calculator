package main

import (
	"fmt"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/pkg/calculator"
	"log"
	"os"
	"sync"
)

func main() {
	logger := log.New(os.Stdout, "CostCalculator ", log.LstdFlags)
	linesPool := sync.Pool{New: func() interface{} {
		lines := make([]byte, 500*1024)
		return lines
	}}
	stringsPool := sync.Pool{New: func() interface{} {
		strs := ""
		return strs
	}}
	costCalculator := calculator.New(logger, &linesPool, &stringsPool, "./data/calculated_costs.csv")
	tariffs, err := costCalculator.ReadAndParseTariffs("./data/tariffs.csv")
	if err != nil {
		fmt.Println("can't read tariffs file")
		os.Exit(1)
	}
	err = costCalculator.ReadAndProcessSessions("./data/sessions.csv", tariffs)
	if err != nil {
		fmt.Println("can't process sessions")
		os.Exit(1)
	}
}
