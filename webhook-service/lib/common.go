package lib

import "os"

func GetNamespaceFromEnvVar() string {
	return os.Getenv("POD_NAMESPACE")
}
