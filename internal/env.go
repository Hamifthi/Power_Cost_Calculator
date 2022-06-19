package internal

import (
	"github.com/spf13/viper"
	"log"
)

type EnvHandler struct {
	logger *log.Logger
}

func NewEnvHandler(logger *log.Logger) *EnvHandler {
	return &EnvHandler{logger: logger}
}

func (eh *EnvHandler) InitializeEnv(envFilePath string) {
	viper.SetConfigFile(envFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		eh.logger.Fatalf("Viper can't read the config file because of %v", err)
	}
	return
}

func (eh *EnvHandler) GetEnv(key string) string {
	value, ok := viper.Get(key).(string)
	if !ok {
		eh.logger.Fatalf("type assertion failed for the key: %s", key)
	}
	return value
}

func ReadTestEnv(filePath string, logger *log.Logger) *EnvHandler {
	envHandler := NewEnvHandler(logger)
	envHandler.InitializeEnv(filePath)
	return envHandler
}
