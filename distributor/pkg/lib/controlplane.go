package lib

import logger "github.com/sirupsen/logrus"

func Register() {
	logger.Info("Registering integration")
}

func Unregister() {
	logger.Info("Unregistering integration")
}
