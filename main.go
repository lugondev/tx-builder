package main

import (
	"github.com/lugondev/tx-builder/cmd"
	log "github.com/sirupsen/logrus"
)

func main() {
	command := cmd.NewCommand()
	if err := command.Execute(); err != nil {
		log.WithError(err).Fatalf("main: execution failed")
	}
	log.Infof("main: execution completed")
}
