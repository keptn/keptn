package lib

import (
	"errors"

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
	Type           string    `yaml:"type"`
	SubscriptionID string    `yaml:"subscriptionID"`
	SendFinished   bool      `yaml:"sendFinished"`
	SendStarted    *bool     `yaml:"sendStarted,omitempty"`
	EnvFrom        []EnvFrom `yaml:"envFrom"`
	Requests       []string  `yaml:"requests"`
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

// DecodeWebHookConfigYAML takes a webhook config string formatted as YAML and decodes it to
// Shipyard value
func DecodeWebHookConfigYAML(webhookConfigYaml []byte) (*WebHookConfig, error) {
	webHookConfig := &WebHookConfig{}

	if err := yaml.Unmarshal(webhookConfigYaml, webHookConfig); err != nil {
		return nil, err
	}

	if len(webHookConfig.Spec.Webhooks) == 0 {
		return nil, errors.New("Webhook configuration invalid: missing 'webhooks[]' part")
	}

	for _, webhook := range webHookConfig.Spec.Webhooks {
		if webhook.Type == "" {
			return nil, errors.New("Webhook configuration invalid: missing 'webhooks[].Type' part")
		}

		if webhook.SubscriptionID == "" {
			return nil, errors.New("Webhook configuration invalid: missing 'webhooks[].SubscriptionID' part")
		}

		if len(webhook.Requests) == 0 {
			return nil, errors.New("Webhook configuration invalid: missing 'webhooks[].Requests[]' part")
		}
	}

	return webHookConfig, nil
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
