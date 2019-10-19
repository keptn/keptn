package handlers

import (
	"encoding/json"
	"testing"

	"github.com/keptn/keptn/api/models"
	"github.com/kinbiko/jsonassert"
)

func TestAddingChannelInfo(t *testing.T) {

	ja := jsonassert.New(t)
	ce := make(map[string]interface{})
	ce["data"] = make(map[string]interface{})
	ce["data"].(map[string]interface{})["project"] = "sockshop"

	ceString, _ := json.Marshal(ce)
	ja.Assertf(string(ceString), `{"data":{"project":"sockshop"}}`)

	channelID := "id"
	token := "token"
	channelInfo := models.ChannelInfo{ChannelID: &channelID, Token: &token}

	forwardData := addChannelInfoInCE(ce, channelInfo)
	actual, _ := json.Marshal(forwardData)
	ja.Assertf(string(actual), `{"data":{"project":"sockshop", "channelInfo":{"channelID":"id", "token":"token"}}}`)
}
