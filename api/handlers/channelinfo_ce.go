package handlers

import "github.com/keptn/keptn/api/models"

type EnrichedCEData struct {
	Data        interface{}        `json:",inline"`
	ChannelInfo models.ChannelInfo `json:"channelInfo"`
}
