package internal

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ShellRechargeSolutionsEU/codechallenge-go-hamed-fathi/entity"
	"os"
)

func WriteFile(fileLocation string, costs []entity.Cost) error {
	var writer *csv.Writer
	if _, err := os.Stat(fileLocation); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(fileLocation)
		if err != nil {
			return err
		}
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			return err
		}
		writer = csv.NewWriter(file)
		err = writer.Write([]string{"session_id", "total_cost"})
		return err
	} else {
		file, err := os.OpenFile(fileLocation, os.O_APPEND|os.O_WRONLY, 0744)
		defer func() {
			err = file.Close()
		}()
		if err != nil {
			return err
		}
		writer = csv.NewWriter(file)
	}
	defer writer.Flush()
	for _, cost := range costs {
		row := []string{cost.SessionID, fmt.Sprintf("%.3f", cost.TotalCost)}
		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
