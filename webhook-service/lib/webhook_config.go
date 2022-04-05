package lib

import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

type WebHookConfig struct {
	ApiVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   Metadata          `yaml:"metadata"`
	Spec       WebHookConfigSpec `yaml:"spec"`
}

// Metadata contains meta-data of a webhook config
type Metadata struct {
	Name string `json:"name" yaml:"name"`
}

type WebHookConfigSpec struct {
	Webhooks []Webhook `yaml:"webhooks"`
}

type Webhook struct {
	Type           string        `yaml:"type"`
	SubscriptionID string        `yaml:"subscriptionID"`
	SendFinished   bool          `yaml:"sendFinished"`
	SendStarted    *bool         `yaml:"sendStarted,omitempty"`
	EnvFrom        []EnvFrom     `yaml:"envFrom"`
	Requests       []interface{} `yaml:"requests"`
}

type EnvFrom struct {
	SecretRef WebHookSecretRef `yaml:"secretRef"`
	Name      string           `yaml:"name"`
}

type Request struct {
	URL     string   `yaml:"url"`
	Method  string   `yaml:"method"`
	Headers []Header `yaml:"headers,omitempty"`
	Payload string   `yaml:"payload,omitempty"`
	Options string   `yaml:"options,omitempty"`
}

type Header struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type WebHookSecretRef struct {
	Key  string `yaml:"key"`
	Name string `yaml:"name"`
}

const webhookConfInvalid = "Webhook configuration invalid: "
const betaApiVersion = "webhookconfig.keptn.sh/v1beta1"

// DecodeWebHookConfigYAML takes a webhook config string formatted as YAML and decodes it to
// Shipyard value
func DecodeWebHookConfigYAML(webhookConfigYaml []byte) (*WebHookConfig, error) {
	webHookConfig := &WebHookConfig{}

	if err := yaml.Unmarshal(webhookConfigYaml, webHookConfig); err != nil {
		return nil, err
	}

	if len(webHookConfig.Spec.Webhooks) == 0 {
		return nil, errors.New(webhookConfInvalid + "missing 'webhooks[]' part")
	}

	for _, webhook := range webHookConfig.Spec.Webhooks {
		if webhook.Type == "" {
			return nil, errors.New(webhookConfInvalid + "missing 'webhooks[].Type' part")
		}

		if webhook.SubscriptionID == "" {
			return nil, errors.New(webhookConfInvalid + "missing 'webhooks[].SubscriptionID' part")
		}

		if len(webhook.Requests) == 0 {
			return nil, errors.New(webhookConfInvalid + "missing 'webhooks[].Requests[]' part")
		}
	}

	if webHookConfig.ApiVersion == betaApiVersion {
		for _, webhook := range webHookConfig.Spec.Webhooks {
			for _, request := range webhook.Requests {
				requestStruct := Request{}
				mapstructure.Decode(request, &requestStruct)
				if err := verifyBeta1Request(requestStruct); err != nil {
					return nil, err
				}
			}
		}
	}

	return webHookConfig, nil
}

func verifyBeta1Request(request Request) error {
	if request.URL == "" {
		//+ validate URL
		return fmt.Errorf(webhookConfInvalid + "webhook request URL empty")
	}
	if request.Method == "" {
		//+validate method
		return fmt.Errorf(webhookConfInvalid + "webhook request method empty")
	}
	if len(request.Headers) > 0 {
		for _, header := range request.Headers {
			if header.Key == "" || header.Value == "" {
				return fmt.Errorf(webhookConfInvalid + "webhook request header or value empty")
			}
		}
	}
	return nil
}

func (wh Webhook) ShouldSendStartedEvent() bool {
	if wh.SendStarted == nil {
		return true
	}
	return *wh.SendStarted
}

func (wh Webhook) ShouldSendFinishedEvent() bool {
	return wh.SendFinished
}
