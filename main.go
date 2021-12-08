package main

import (
	"log"

	"github.com/blainemoser/gitsync/configs"
	"github.com/blainemoser/gitsync/logging"
	"github.com/blainemoser/gitsync/queue"
	"github.com/blainemoser/gitsync/utils"
)

var Queue *queue.Queue

func main() {
	pwd()
	logging.StaticWrite("Starting", "INFO")
	<-start()
}

func start() chan bool {
	configs, err := configs.NewConfigs().SetDirectories("configs.json", "")
	fatal(err)
	Queue, err = queue.NewQueue().Walk(configs)
	fatal(err)
	Queue.StandbyAll()
	return make(chan bool, 1)
}

func pwd() {
	pwd, err := utils.GetPWD()
	if err != nil {
		log.Printf("warning: could not set the working directory: %s\n", err.Error())
		return
	}
	logging.SetBaseDir(pwd)
}

func fatal(err error) {
	if err != nil {
		logging.StaticWrite(err.Error(), "ERROR")
		log.Fatal(err)
	}
}
