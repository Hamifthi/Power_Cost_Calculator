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

func (env *EnvHandler) InitializeEnv(envFilePath string) {
	viper.SetConfigFile(envFilePath)
	err := viper.ReadInConfig()
	if err != nil {
		env.logger.Fatalf("Viper can't read the config file because of %v", err)
	}
	return
}

func (env *EnvHandler) GetEnv(key string) string {
	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("type assertion failed for the key: %s", key)
	}
	return value
}
