package main

import (
	"github.com/lugondev/tx-builder/cmd"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {

	err := LoadConfig(".")
	if err != nil {
		log.Fatal("error load config:", err)
	}
	command := cmd.NewCommand()
	if err := command.Execute(); err != nil {
		log.WithError(err).Fatalf("main: execution failed")
	}

	log.Infof("main: execution completed")
}

func LoadConfig(path string) (err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("configuration")
	viper.SetConfigType("yml")

	return viper.ReadInConfig()
}
