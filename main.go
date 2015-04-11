package main

import (
	"github.com/maleck13/photoProcessor/api"
	"github.com/maleck13/photoProcessor/logger"
	"github.com/maleck13/photoProcessor/conf"
	"github.com/maleck13/photoProcessor/messaging"
)

func main() {

	logger.InitLogger()
	conf.LoadConfig()
	StartUp()

}

func StartUp() {
	messaging.StartMessaging()
	api.StartApi()
}
