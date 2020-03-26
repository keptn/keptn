package docker

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/keptn/keptn/cli/pkg/logging"
)

// SplitImageName splits an image name into its name and tag
// If no tag is provided, latest is returned
func SplitImageName(imageWithTag string) (string, string) {

	// Get image name without Docker-organization
	splitsIntoImage := strings.Split(imageWithTag, "/")
	imageName := splitsIntoImage[len(splitsIntoImage)-1]

	splitsIntoTag := strings.Split(imageName, ":")
	if len(splitsIntoTag) == 2 {
		// Tag is provided in the image name
		tag := splitsIntoTag[len(splitsIntoTag)-1]
		return strings.TrimSuffix(imageWithTag, ":"+tag), tag
	}
	// Otherwise use latest tag
	return imageWithTag, "latest"
}

// CheckImageAvailability checks the availability of a image which is hosted on Docker or on Quay
func CheckImageAvailability(image, tag string) error {

	if strings.HasPrefix(image, "docker.io/") {
		resp, err := http.Get("https://index.docker.io/v1/repositories/" +
			strings.TrimPrefix(image, "docker.io/") + "/tags/" + tag)
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return errors.New("Provided image not found: " + string(body))
	} else if strings.HasPrefix(image, "quay.io/") {
		resp, err := http.Get("https://quay.io/api/v1/repository/" +
			strings.TrimPrefix(image, "quay.io/") + "/tag/" + tag + "/images")
		if err != nil {
			return err
		}
		if resp.StatusCode == http.StatusOK {
			return nil
		}
		return errors.New("Provided image not found: " + resp.Status)
	}
	logging.PrintLog("Availability of provided image cannot be checked.", logging.InfoLevel)
	return nil
}
