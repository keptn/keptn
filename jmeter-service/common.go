package main

import (
	"errors"
	"fmt"
	keptncommon "github.com/keptn/go-utils/pkg/lib/keptn"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ghodss/yaml"
	models "github.com/keptn/go-utils/pkg/api/models"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
)

var runlocal = (os.Getenv("env") == "runlocal")

// ErrPrimaryFileNotAvailable indicates that the primary test file is not available
var ErrPrimaryFileNotAvailable = errors.New("primary test file not available")

//
// Iterates through the JMeterConf and returns the workload configuration matching the testStrategy
// If no config is found in JMeterConf it falls back to the defaults
//
func getWorkload(jmeterconf *JMeterConf, teststrategy string) (*Workload, error) {
	// get the entry for the passed strategy
	if jmeterconf != nil && jmeterconf.Workloads != nil {
		for _, workload := range jmeterconf.Workloads {
			if workload.TestStrategy == teststrategy {
				return workload, nil
			}
		}
	}

	// if we didnt find it in the config go through the defaults
	for _, workload := range defaultWorkloads {
		if workload.TestStrategy == teststrategy {
			return &workload, nil
		}
	}

	return nil, errors.New("No workload configuration found for teststrategy: " + teststrategy)
}

// GetConfigurationServiceURL returns the URL of the configuration service
func GetConfigurationServiceURL() string {
	if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE") == "" {
		return "configuration-service:8080"
	} else if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE") != "" {
		return os.Getenv("CONFIGURATION_SERVICE")
	}
	return "localhost:8080"
}

/** GetKeptnResource Loads a Resource from the Keptn configuration repository
 * returns:
 * - fileContent if found or "" if no file found at all
 * - error: in case there was an error
 */
func GetKeptnResource(project string, stage string, service string, resourceUri string) (string, error) {
	resourceHandler := configutils.NewResourceHandler(GetConfigurationServiceURL())
	resource, err := resourceHandler.GetServiceResource(project, stage, service, resourceUri)
	if err != nil && err == configutils.ResourceNotFoundError {
		// if not found on serivce level - lets try it on stage level
		resource, err = resourceHandler.GetStageResource(project, stage, resourceUri)

		if err != nil && err == configutils.ResourceNotFoundError {
			// if not found on stage level we try project level
			resource, err = resourceHandler.GetProjectResource(project, resourceUri)

			if err != nil && err == configutils.ResourceNotFoundError {
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

/*
 * GetAllKeptnResources This function will download ALL Resources from Keptn's Configuration Repository where the name starts with 'resourceUriStartWith'. This for instance allows us to download all files in the /jmeter folders
 *
 * Parameters:
 * project, stage, string: reference the keptn repo
 * inheritResources: if true it will download all resources from service, stage and project level - otherwise just from service level
 * resourceUriFolderOfInterest: will only download resources where the resourceUri contains that value, e.g: "/jmeter" and then also stores the downloaded files under that prefix
 * primaryTestFileName: if specified - the implementation makes sure to download this file locally
 * localDirectory: the local directory to store these downloaded files
 *
 * Return:
 * foundPrimaryFile: true if it was downloaded
 * no of resources: total number of downloaded resources
 * error: any error that occured
 */
func GetAllKeptnResources(project string, stage string, service string, inheritResources bool, resourceURIFolderOfInterest string, primaryTestFileName string, localDirectory string, logger *keptncommon.Logger) (bool, int, error) {

	resourceHandler := configutils.NewResourceHandler(GetConfigurationServiceURL())

	// Lets first get the servcie resources
	// TODO: This endpoint is not yet implemented and therefore this always fails - https://github.com/keptn/keptn/issues/1924
	/* resourceList, err := resourceHandler.GetAllServiceResources(project, stage, service)
	if err != nil {
		return 0, err
	}*/

	resourceList := []*models.Resource{}

	// Next - lets get stage and project resources!
	// if inheritResources == true we also get the list of resources from stage and project level
	if inheritResources {
		stageResources, err := resourceHandler.GetAllStageResources(project, stage)
		if err != nil {
			return false, 0, err
		}
		resourceList = append(resourceList, stageResources...)

		// TODO: missing configutils.GetAllProjectResources(project)
		/* projectResources, err := resourceHandler.GetAllProjectResources(project)
		if err != nil {
			return 0, err
		}
		resourceList = append(resourceList, projectResources...)*/
	}

	fileCount := 0
	skippedFileCount := 0
	foundPrimaryFile := false

	// Download Files
	// now lets iterate through all resources and download those that match the resourceURIFolderOfInterest and that havent already been downloaded
	// as we download files from project, service and stage level we have different file structures, e.g:
	// Project: /jmeter/myjmeter.jmx
	// Stage: /jmeter/myjmeter2.jmx
	// Stage: /myservice/jmeter/myjmeter3.jmx
	// When we store it locally we have to store all these files in /jmeter/filename.jmx
	for _, resource := range resourceList {
		isPrimaryFile := strings.Contains(*resource.ResourceURI, primaryTestFileName)
		startingIndex := strings.Index(*resource.ResourceURI, resourceURIFolderOfInterest)

		// store to local directory if it doesnt already exist
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
			downloadedResource, err := resourceHandler.GetStageResource(project, stage, *resource.ResourceURI)
			if err != nil {
				return false, fileCount, err
			}

			logger.Debug(fmt.Sprintf("Storing %s to %s/%s - size (%d)", *resource.ResourceURI, localDirectory, targetFileName, len(downloadedResource.ResourceContent)))
			if err := storeFile(localDirectory, targetFileName, downloadedResource.ResourceContent, true); err != nil {
				return false, fileCount, err
			}
			fileCount = fileCount + 1
		} else {
			skippedFileCount = skippedFileCount + 1
			// 	logger.Debug(fmt.Sprintf("Not storing %s as it doesnt match %s or %s", *resource.ResourceURI, primaryTestFileName, resourceURIFolderOfInterest))
		}
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
			return false, fileCount, fmt.Errorf("could not store primary file in %s/%s: %s", localDirectory, primaryTestFileName, err.Error())
		}
		fileCount = fileCount + 1
		foundPrimaryFile = true
	}

	logger.Debug(fmt.Sprintf("Downloaded %d and skipped %d files for %s in %s.%s.%s", fileCount, skippedFileCount, resourceURIFolderOfInterest, project, stage, service))

	return foundPrimaryFile, fileCount, nil
}

/**
 * FileExists just returns whether the file exists
 */
func FileExists(filename string) bool {
	// lets first check if the file exists and if we should not overwrite it
	if _, err := os.Stat(filename); err == nil {
		return true
	}
	return false
}

/**
 * Stores the content to the local file system under the targetFileName (can also contain directories)
 * Returns:
 * 1: true if file was actually written, e.g: will be false if file exists and overwriteIfExists==False
 * 2: error if an error occured
 */
func storeFile(localDirectory string, targetFileName string, resourceContent string, overwriteIfExists bool) error {

	// lets construct the final directory name
	if !strings.HasSuffix(localDirectory, "/") {
		localDirectory = localDirectory + "/"
	}
	directory := localDirectory
	finalLocalFilename := localDirectory + targetFileName

	// first lets first check if the file exists and if we should not overwrite it
	if FileExists(finalLocalFilename) && !overwriteIfExists {
		return nil
	}

	// add every single piece of the path excluding the filename itself to the directory
	pathArr := strings.Split(targetFileName, "/")
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	// now lets create that directory if it doesnt exist
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

//
// Loads jmeter.conf for the current service
//
func getJMeterConf(project string, stage string, service string, logger *keptncommon.Logger) (*JMeterConf, error) {
	// if we run in a runlocal mode we are just getting the file from the local disk
	var fileContent []byte
	var err error
	if runlocal {
		fileContent, err = ioutil.ReadFile(JMeterConfFilename)
		if err != nil {
			logMessage := fmt.Sprintf("No %s file found LOCALLY for service %s in stage %s in project %s", JMeterConfFilename, service, stage, project)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
	} else {

		logger.Info(fmt.Sprintf("Loading %s for %s.%s.%s", JMeterConfFilename, project, stage, service))

		keptnResourceContent, err := GetKeptnResource(project, stage, service, JMeterConfFilename)

		if err != nil {
			logMessage := fmt.Sprintf("error when trying to load %s file for service %s on stage %s or project-level %s", JMeterConfFilename, service, stage, project)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
		if keptnResourceContent == "" {
			// if no jmeter.conf file is available, this is not an error, as the service will proceed with the default workload
			logger.Info(fmt.Sprintf("no %s found", JMeterConfFilename))
			return nil, nil
		}
		fileContent = []byte(keptnResourceContent)
	}

	var jmeterConf *JMeterConf
	jmeterConf, err = parseJMeterConf(fileContent)

	if err != nil {
		logMessage := fmt.Sprintf("Couldn't parse %s file found for service %s in stage %s in project %s. Error: %s", JMeterConfFilename, service, stage, project, err.Error())
		logger.Error(logMessage)
		return nil, errors.New(logMessage)
	}

	logger.Debug(fmt.Sprintf("Successfully loaded jmeter.conf.yaml with %d workloads", len(jmeterConf.Workloads)))

	return jmeterConf, nil
}

//
// parses content and maps it to the JMeterConf struct
//
func parseJMeterConf(input []byte) (*JMeterConf, error) {
	jmeterconf := &JMeterConf{}
	err := yaml.Unmarshal([]byte(input), &jmeterconf)

	if err != nil {
		return nil, err
	}

	return jmeterconf, nil
}
