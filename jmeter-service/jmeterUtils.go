package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/utils"
)

const maxAcceptedErrorRate = 0.1

// TestInfo contains information about which test to execute
type TestInfo struct {
	Project      string
	Stage        string
	Service      string
	TestStrategy string
}

// ToString returns a string representation of a TestInfo object
func (ti *TestInfo) ToString() string {
	return "Project: " + ti.Project + ", Service: " + ti.Service + ", Stage: " + ti.Stage + ", TestStrategy: " + ti.TestStrategy
}

func getConfigurationServiceURL() string {
	if os.Getenv("env") == "production" {
		return "configuration-service:8080"
	}
	return "localhost:8080"
}

func executeJMeter(testInfo *TestInfo, scriptName string, resultsDir string, serverURL string, serverPort int, checkPath string, vuCount int,
	loopCount int, thinkTime int, LTN string, funcValidation bool, avgRtValidation int, logger *keptnutils.Logger) (bool, error) {
	os.RemoveAll(resultsDir)
	os.MkdirAll(resultsDir, 0644)

	resourceHandler := keptnutils.NewResourceHandler(getConfigurationServiceURL())
	testScriptResource, err := resourceHandler.GetServiceResource(testInfo.Project, testInfo.Stage, testInfo.Service, scriptName)

	// if no test file has been found, we assume that no tests should be executed
	if err != nil || testScriptResource == nil || testScriptResource.ResourceContent == "" {
		logger.Debug("Skipping test execution because no tests have been defined.")
		return true, nil
	}

	os.RemoveAll(scriptName)
	pathArr := strings.Split(scriptName, "/")
	directory := ""
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return false, err
	}
	testScriptFile, err := os.Create(scriptName)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	defer testScriptFile.Close()

	_, err = testScriptFile.Write([]byte(testScriptResource.ResourceContent))

	if err != nil {
		logger.Error(err.Error())
		return false, err
	}

	testInfoStr := testInfo.ToString() + ", scriptName: " + scriptName + ", serverURL: " + serverURL
	logger.Debug("Starting JMeter test. " + testInfoStr)
	res, err := keptnutils.ExecuteCommand("jmeter", []string{"-n", "-t", "./" + scriptName,
		// "-e", "-o", resultsDir,
		"-l", resultsDir + "_result.tlf",
		"-JSERVER_URL=" + serverURL,
		"-JDT_LTN=" + LTN,
		"-JVUCount=" + strconv.Itoa(vuCount),
		"-JLoopCount=" + strconv.Itoa(loopCount),
		"-JCHECK_PATH=" + checkPath,
		"-JSERVER_PORT=" + strconv.Itoa(serverPort),
		"-JThinkTime=" + strconv.Itoa(thinkTime)})

	logger.Info(res)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}

	// Parse result
	summary := getLastOccurence(strings.Split(res, "\n"), "summary =")
	if summary == "" {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo.ToString())
	}

	space := regexp.MustCompile(`\s+`)
	splits := strings.Split(space.ReplaceAllString(summary, " "), " ")
	runs, err := strconv.Atoi(splits[2])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo.ToString())
	}

	errorCount, err := strconv.Atoi(splits[14])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo.ToString())
	}

	avg, err := strconv.Atoi(splits[8])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo.ToString())
	}

	if funcValidation && errorCount > 0 {
		logger.Debug(fmt.Sprintf("Function validation failed because we got %d errors.", errorCount) + ". " + testInfo.ToString())
		return false, nil
	}

	maxAcceptedErrors := maxAcceptedErrorRate * float64(runs)
	if errorCount > int(maxAcceptedErrors) {
		logger.Debug(fmt.Sprintf("jmeter test failed because we got a too high error rate of %.2f.", float64(errorCount)/float64(runs)) + ". " + testInfo.ToString())
		return false, nil
	}

	if avgRtValidation > 0 && avg > avgRtValidation {
		logger.Debug(fmt.Sprintf("Avg rt validation failed because we got an avg rt of %d", avgRtValidation) + ". " + testInfo.ToString())
		return false, nil
	}

	logger.Debug("Successfully executed JMeter test. " + testInfo.ToString())
	return true, nil
}

func getLastOccurence(vs []string, prefix string) string {
	for i := len(vs) - 1; i >= 0; i-- {
		if strings.HasPrefix(vs[i], prefix) {
			return vs[i]
		}
	}
	return ""
}
