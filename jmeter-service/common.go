package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	configutils "github.com/keptn/go-utils/pkg/api/utils"
	keptn "github.com/keptn/go-utils/pkg/lib"
)

var runlocal = (os.Getenv("env") == "runlocal")

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

func GetConfigurationServiceURL() string {
	if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE_URL") == "" {
		return "configuration-service:8080"
	} else if os.Getenv("env") == "production" && os.Getenv("CONFIGURATION_SERVICE_URL") != "" {
		return os.Getenv("CONFIGURATION_SERVICE_URL")
	}
	return "localhost:8080"
}

/**
 * Loads a Resource from the Keptn configuration repository
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

//
// Loads jmeter.conf for the current service
//
func getJMeterConf(project string, stage string, service string, logger *keptn.Logger) (*JMeterConf, error) {
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

		keptnResourceContent, err := GetKeptnResource(project, stage, service, JMeterConfFilename)

		/* resourceHandler := keptnapi.NewResourceHandler(GetConfigurationServiceURL())
		keptnResourceContent, err := resourceHandler.GetServiceResource(project, stage, service, JMeterConfFilename)*/
		if err != nil {
			logMessage := fmt.Sprintf("No %s file found for service %s on stage %s or project-level %s", JMeterConfFilename, service, stage, project)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
		fileContent = []byte(keptnResourceContent /*keptnResourceContent.ResourceContent*/)
	}

	var jmeterConf *JMeterConf
	jmeterConf, err = parseJMeterConf(fileContent)

	if err != nil {
		logMessage := fmt.Sprintf("Couldn't parse %s file found for service %s in stage %s in project %s. Error: %s", JMeterConfFilename, service, stage, project, err.Error())
		logger.Error(logMessage)
		return nil, errors.New(logMessage)
	}

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
