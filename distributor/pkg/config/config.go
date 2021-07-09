package config

import (
	"net/url"
	"strings"
)

func GetHTTPPollingEndpoint(env EnvConfig) string {
	endpoint := env.KeptnAPIEndpoint
	if endpoint == "" {
		if endpoint == "" {
			return DefaultEventsEndpoint
		}
	} else {
		endpoint = strings.TrimSuffix(env.KeptnAPIEndpoint, "/") + "/controlPlane/v1/event/triggered"
	}

	parsedURL, _ := url.Parse(endpoint)

	if parsedURL.Scheme == "" {
		parsedURL.Scheme = "http"
	}
	if parsedURL.Path == "" {
		parsedURL.Path = "v1/event/triggered"
	}

	return parsedURL.String()
}

func GetPubSubRecipientURL(env EnvConfig) string {
	recipientService := env.PubSubRecipient

	if !strings.HasPrefix(recipientService, "https://") && !strings.HasPrefix(recipientService, "http://") {
		recipientService = "http://" + recipientService
	}

	path := ""
	if env.PubSubRecipientPath != "" {
		path = "/" + strings.TrimPrefix(env.PubSubRecipientPath, "/")
	}
	return recipientService + ":" + env.PubSubRecipientPort + path
}
