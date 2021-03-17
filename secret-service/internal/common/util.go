package common

import "os"

type StringSupplier func() string

func EnvBasedStringSupplier(envVarName, defaultVal string) StringSupplier {
	return func() string {
		ns := os.Getenv(envVarName)
		if ns != "" {
			return ns
		}
		return defaultVal
	}
}
