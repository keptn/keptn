package handlers

import (
	"encoding/json"
	"testing"

	"github.com/keptn/keptn/api/models"
	"github.com/kinbiko/jsonassert"
)

func TestAddingChannelInfo(t *testing.T) {

	ja := jsonassert.New(t)
	ceData := make(map[string]interface{})
	ceData["project"] = "sockshop"

	ceString, _ := json.Marshal(ceData)
	ja.Assertf(string(ceString), `{"project":"sockshop"}`)

	channelID := "id"
	token := "token"
	channelInfo := models.ChannelInfo{ChannelID: &channelID, Token: &token}

	forwardData := addChannelInfoInCE(ceData, channelInfo)
	actual, _ := json.Marshal(forwardData)
	ja.Assertf(string(actual), `{"project":"sockshop", "channelInfo":{"channelID":"id", "token":"token"}}`)
}
