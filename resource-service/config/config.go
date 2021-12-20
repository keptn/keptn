package config

import "os"

// ConfigDir specifies the base path of the config repos in the file system
var ConfigDir = getConfigDir()

func getConfigDir() string {
	if os.Getenv("env") == "production" {
		return "/data/config"
	}
	return "./debug/config"
}
