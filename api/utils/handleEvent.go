// This file is safe to edit. Once it exists it will not be overwritten

package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

// PostToEventBroker makes a post request to the eventbroker
func PostToEventBroker(e interface{}, logger *keptnutils.Logger) error {

	data, err := json.Marshal(e)
	if err != nil {
		logger.Error(fmt.Sprintf("Error marshaling data %s", err.Error()))
		return err
	}

	url := "http://" + os.Getenv("EVENTBROKER_URI")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/cloudevents+json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		logger.Error(fmt.Sprintf("Error making POST %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		return nil
	}
	return errors.New(resp.Status)
}
