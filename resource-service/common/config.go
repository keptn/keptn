package common

import "os"

// ConfigDir specifies the base path of the config repos in the file system
const defaultConfigDir = "/data/config"

func GetConfigDir() string {
	if configDir, ok := os.LookupEnv("CONFIG_DIR"); ok {
		return configDir
	}
	return defaultConfigDir
}
