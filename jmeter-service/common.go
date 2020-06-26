package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/ghodss/yaml"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
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
		resourceHandler := keptnapi.NewResourceHandler("configuration-service:8080")
		keptnResourceContent, err := resourceHandler.GetServiceResource(project, stage, service, JMeterConfFilename)
		if err != nil {
			logMessage := fmt.Sprintf("No %s file found for service %s in stage %s in project %s", JMeterConfFilename, service, stage, project)
			logger.Info(logMessage)
			return nil, errors.New(logMessage)
		}
		fileContent = []byte(keptnResourceContent.ResourceContent)
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
