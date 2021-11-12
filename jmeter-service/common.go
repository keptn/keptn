package main

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/common/fileutils"
	logger "github.com/sirupsen/logrus"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	configutils "github.com/keptn/go-utils/pkg/api/utils"
)

// Opaque key type used for graceful shutdown context value
type gracefulShutdownKeyType struct{}

var gracefulShutdownKey = gracefulShutdownKeyType{}

// Opaque key type used for graceful shutdown context value
type keptnQuitType struct{}

var keptnQuit = keptnQuitType{}

// ErrPrimaryFileNotAvailable indicates that the primary test file is not available
var ErrPrimaryFileNotAvailable = errors.New("primary test file not available")

// GetConfigurationServiceURL returns the URL of the configuration service
func GetConfigurationServiceURL() string {
	if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE") == "" {
		return "configuration-service:8080"
	} else if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE") != "" {
		return os.Getenv("CONFIGURATION_SERVICE")
	}

	return "localhost:8080"
}

// GetKeptnResource Loads a Resource from the Keptn configuration repository
// returns:
// - fileContent if found or "" if no file found at all
// - error: in case there was an error
func GetKeptnResource(project string, stage string, service string, resourceURI string) (string, error) {
	resourceHandler := configutils.NewResourceHandler(GetConfigurationServiceURL())
	resource, err := resourceHandler.GetServiceResource(project, stage, service, resourceURI)
	if err != nil && errors.Is(err, configutils.ResourceNotFoundError) {
		// if not found on service level - lets try it on stage level
		resource, err = resourceHandler.GetStageResource(project, stage, resourceURI)

		if err != nil && errors.Is(err, configutils.ResourceNotFoundError) {
			// if not found on stage level we try project level
			resource, err = resourceHandler.GetProjectResource(project, resourceURI)

			if err != nil && errors.Is(err, configutils.ResourceNotFoundError) {
				return "", nil
			} else if err != nil {
				return "", err
			}
		} else if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	return resource.ResourceContent, nil
}

// DownloadAndStoreResources downloads all resources from Keptn's Configuration Repository where the name starts with 'resourceUriStartWith'.
// This for instance allows us to download all files in the /jmeter folders
//
// Parameters:
// project, stage, service: reference the keptn repo
// resourceUriFolderOfInterest: will only download resources where the resourceUri contains that value, e.g: "/jmeter" and then also stores the downloaded files under that prefix
// primaryTestFileName: if specified - the implementation makes sure to download this file locally
// localDirectory: the local directory to store these downloaded files
//
// Return:
// foundPrimaryFile: true if it was downloaded
// number of resources: total number of downloaded resources
// error: any error that occurred
func DownloadAndStoreResources(project string, stage string, service string, resourceURIFolderOfInterest string, primaryTestFileName string, localDirectory string) (bool, int, error) {
	foundPrimaryFile, fileCount, err := getAllKeptnResources(project, stage, resourceURIFolderOfInterest, primaryTestFileName, localDirectory)
	if err != nil {
		return foundPrimaryFile, fileCount, err
	}
	// Fallback if primary file wasn't loaded yet
	// last effort - if we couldn't download the specific test file because, e.g: limitations of our current API - then simply go back to download this specific file
	if !foundPrimaryFile {
		primaryTestFileContent, err := GetKeptnResource(project, stage, service, primaryTestFileName)
		if err != nil {
			return false, fileCount, err
		}
		if primaryTestFileContent == "" {
			return false, fileCount, ErrPrimaryFileNotAvailable
		}
		logger.Debug(fmt.Sprintf("Storing primary file in %s/%s - size(%d)", localDirectory, primaryTestFileName, len(primaryTestFileContent)))
		if err := storeFile(localDirectory, primaryTestFileName, primaryTestFileContent, true); err != nil {
			return false, fileCount, fmt.Errorf("could not store primary file in %s/%s: %w", localDirectory, primaryTestFileName, err)
		}
		fileCount++
		foundPrimaryFile = true
	}
	return foundPrimaryFile, fileCount, nil
}

func getAllKeptnResources(project string, stage string, resourceURIFolderOfInterest string, primaryTestFileName string, localDirectory string) (bool, int, error) {
	// NOTE: we should also implement and use missing configutils.GetAllProjectResources(project) & configutils.GetAllServiceResources(project,service)
	// and merge them into the resource list
	resourceHandler := configutils.NewResourceHandler(GetConfigurationServiceURL())
	resourceList, err := resourceHandler.GetAllStageResources(project, stage)
	if err != nil {
		return false, 0, err
	}
	// iterate over all resources and download those that match the resourceURIFolderOfInterest
	// when we store it locally we have to store all these files in /jmeter/filename.jmx
	var fileCount, skippedFileCount int
	var foundPrimaryFile bool
	for _, resource := range resourceList {
		isPrimaryFile := strings.Contains(*resource.ResourceURI, primaryTestFileName)
		startingIndex := strings.Index(*resource.ResourceURI, resourceURIFolderOfInterest)

		// now lets strip off the any prepending directory names prior to resourceURIFolderOfInterest
		targetFileName := ""
		if startingIndex >= 0 {
			targetFileName = (*resource.ResourceURI)[startingIndex:]
		}
		if isPrimaryFile {
			targetFileName = primaryTestFileName
			foundPrimaryFile = true
		}
		// only store it if we really know whether and where we have to store it to!
		if targetFileName != "" {
			// now we have to download that resource first as so far we only have the resourceURI
			downloadedResource, err := resourceHandler.GetStageResource(project, stage, strings.TrimPrefix(*resource.ResourceURI, "/"))
			if err != nil {
				return false, fileCount, err
			}
			logger.Debugf("Storing %s to %s/%s - size (%d)", *resource.ResourceURI, localDirectory, targetFileName, len(downloadedResource.ResourceContent))
			if err := storeFile(localDirectory, targetFileName, downloadedResource.ResourceContent, true); err != nil {
				return false, fileCount, err
			}
			fileCount++
		} else {
			skippedFileCount++
			logger.Debugf("Not storing %s as it doesn't match %s or %s", *resource.ResourceURI, primaryTestFileName, resourceURIFolderOfInterest)
		}
	}
	return foundPrimaryFile, fileCount, nil
}

// storeFile stores the content to the local file system under the targetFileName (can also contain directories)
// Returns:
// 1: true if file was actually written, e.g: will be false if file exists and overwriteIfExists==False
// 2: error if an error occurred
func storeFile(localDirectory string, targetFileName string, resourceContent string, overwriteIfExists bool) error {
	// lets construct the final directory name
	if !strings.HasSuffix(localDirectory, "/") {
		localDirectory += "/"
	}
	directory := localDirectory
	finalLocalFilename := localDirectory + targetFileName

	// first lets first check if the file exists and if we should not overwrite it
	if fileutils.FileExists(finalLocalFilename) && !overwriteIfExists {
		return nil
	}

	// add every single piece of the path excluding the filename itself to the directory
	pathArr := strings.Split(targetFileName, "/")
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	// now lets create that directory if it doesn't exist
	err := os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return err
	}

	// now we store the file
	writeToFile, err := os.Create(finalLocalFilename)
	if err != nil {
		return err
	}
	defer writeToFile.Close()

	_, err = writeToFile.Write([]byte(resourceContent))
	if err != nil {
		return err
	}

	return nil
}

// getJMeterConf Loads jmeter.conf for the current service
func getJMeterConf(testInfo TestInfo) (*JMeterConf, error) {
	// if we run in a runlocal mode we are just getting the file from the local disk
	var fileContent []byte
	var err error
	logger.Info(fmt.Sprintf("Loading %s for %s.%s.%s", JMeterConfFilename, testInfo.Project, testInfo.Stage, testInfo.Service))
	keptnResourceContent, err := GetKeptnResource(testInfo.Project, testInfo.Stage, testInfo.Service, JMeterConfFilename)

	if err != nil {
		logMessage := fmt.Sprintf("error when trying to load %s file for service %s on stage %s or project-level %s", JMeterConfFilename, testInfo.Service, testInfo.Stage, testInfo.Project)
		logger.Info(logMessage)
		return nil, errors.New(logMessage)
	}
	if keptnResourceContent == "" {
		// if no jmeter.conf file is available, this is not an error, as the service will proceed with the default workload
		logger.Info(fmt.Sprintf("no %s found", JMeterConfFilename))
		return nil, nil
	}
	fileContent = []byte(keptnResourceContent)

	var jmeterConf *JMeterConf
	jmeterConf, err = decodeJmeterConf(fileContent)
	if err != nil {
		logMessage := fmt.Sprintf("Couldn't parse %s file found for service %s in stage %s in project %s. Error: %s", JMeterConfFilename, testInfo.Service, testInfo.Stage, testInfo.Project, err.Error())
		logger.Error(logMessage)
		return nil, errors.New(logMessage)
	}

	logger.Debug(fmt.Sprintf("Successfully loaded jmeter.conf.yaml with %d workloads", len(jmeterConf.Workloads)))

	return jmeterConf, nil
}

func decodeJmeterConf(input []byte) (*JMeterConf, error) {
	jmeterconf := &JMeterConf{}
	err := yaml.Unmarshal(input, &jmeterconf)
	if err != nil {
		return nil, err
	}
	return jmeterconf, nil
}
