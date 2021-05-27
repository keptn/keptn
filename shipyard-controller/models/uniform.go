package models

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

type GetUniformIntegrationParams struct {
	Name    string `form:"name" json:"name"`
	ID      string `form:"id" json:"id"`
	Project string `form:"project" json:"project"`
	Stage   string `form:"stage" json:"stage"`
	Service string `form:"service" json:"service"`
}

// TODO: delete this and use structs defined in go-utils

type Integration struct {
	ID           string       `json:"id" bson:"_id"`
	Name         string       `json:"name" bson:"name"`
	MetaData     MetaData     `json:"metadata" bson:"metadata"`
	Subscription Subscription `json:"subscription" bson:"subscription"`
}

type MetaData struct {
	Hostname           string             `json:"hostname" bson:"hostname"`
	DeploymentName     string             `json:"deplyomentname" bson:"deploymentname"`
	IntegrationVersion string             `json:"integrationversion" bson:"integrationversion"`
	DistributorVersion string             `json:"distributorversion" bson:"distributorversion"`
	Status             string             `json:"status" bson:"status"`
	Location           string             `json:"location" bson:"location"`
	KubernetesMetaData KubernetesMetaData `json:"kubernetesmetadata" bson:"kubernetesmetadata"`
}

type Subscription struct {
	Topics []string           `json:"topics" bson:"topics"`
	Status string             `json:"status" bson:"status"`
	Filter SubscriptionFilter `json:"filter" bson:"filter"`
}

type SubscriptionFilter struct {
	Project string `json:"project" bson:"project"`
	Stage   string `json:"stage" bson:"stage"`
	Service string `json:"service" bson:"service"`
}

type KubernetesMetaData struct {
	Namespace      string `json:"namespace" bson:"namespace"`
	PodName        string `json:"podname" bson:"podname"`
	DeploymentName string `json:"deploymentname" bson:"deploymentname"`
}

type IntegrationID struct {
	Name      string `json:"name" bson:"name"`
	Namespace string `json:"namespace" bson:"namespace"`
	Project   string `json:"project" bson:"project"`
	Stage     string `json:"stage" bson:"stage"`
	Service   string `json:"service" bson:"service"`
}

func (i IntegrationID) Hash() (string, error) {
	if !i.validate() {
		return "", fmt.Errorf("incomplete integration ID. At least 'name' and 'namespace' must be set.")
	}
	raw := fmt.Sprintf("%s-%s-%s-%s-%s", i.Name, i.Namespace, i.Project, i.Stage, i.Service)
	hasher := sha1.New() //nolint:gosec
	hasher.Write([]byte(raw))
	return hex.EncodeToString(hasher.Sum(nil)), nil
}

func (i IntegrationID) validate() bool {
	return i.Name != "" && i.Namespace != ""
}
