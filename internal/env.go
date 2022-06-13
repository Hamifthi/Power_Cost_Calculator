package internal

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func InitializeEnv(envFilePath string) error {
	viper.SetConfigFile(envFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "Viper can't read the config file")
	}
	return nil
}

func GetEnv(key string) (string, error) {
	value, ok := viper.Get(key).(string)
	if !ok {
		return "", fmt.Errorf("type assertion failed for the key: %s", key)
	}
	return value, nil
}
