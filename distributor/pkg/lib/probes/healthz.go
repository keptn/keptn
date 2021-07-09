package probes

import (
	"github.com/keptn/keptn/distributor/pkg/config"
	"net/http"
)

func Readyz(connectionType config.ConnectionType, env config.EnvConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {

		if connectionType == config.ConnectionTypeHTTP {
			url := config.GetPubSubRecipientURL(env)
			NewReachabilityChecker(&Config{})

		}
		w.WriteHeader(http.StatusOK)
	}
}
