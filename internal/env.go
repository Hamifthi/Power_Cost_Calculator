package internal

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"log"
)

func InitializeEnv(envFilePath string) error {
	viper.SetConfigFile(envFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Viper can't read the config file")
	}
	return nil
}

func GetEnv(key string, logger *log.Logger) (string, error) {
	value, ok := viper.Get(key).(string)
	if !ok {
		logger.Printf("type assertion failed for the key: %s", key)
	}
	return value, nil
}
