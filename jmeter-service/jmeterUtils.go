package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	configutils "github.com/keptn/go-utils/pkg/configuration-service/utils"
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

func executeJMeter(testInfo *TestInfo, workload *Workload, resultsDir string, url *url.URL, LTN string, funcValidation bool, logger *keptnutils.Logger) (bool, error) {
	os.RemoveAll(resultsDir)
	os.MkdirAll(resultsDir, 0644)

	resourceHandler := configutils.NewResourceHandler(getConfigurationServiceURL())
	testScriptResource, err := resourceHandler.GetServiceResource(testInfo.Project, testInfo.Stage, testInfo.Service, workload.Script)

	// if no test file has been found, we assume that no tests should be executed
	if err != nil || testScriptResource == nil || testScriptResource.ResourceContent == "" {
		logger.Debug("Skipping test execution because no tests have been defined.")
		return true, nil
	}

	os.RemoveAll(workload.Script)
	pathArr := strings.Split(workload.Script, "/")
	directory := ""
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	err = os.MkdirAll(directory, os.ModePerm)
	if err != nil {
		return false, err
	}
	testScriptFile, err := os.Create(workload.Script)
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

	testInfoStr := testInfo.ToString() + ", scriptName: " + workload.Script + ", serverURL: " + url.String()
	logger.Debug("Starting JMeter test. " + testInfoStr)
	res, err := keptnutils.ExecuteCommand("jmeter", []string{"-n", "-t", "./" + workload.Script,
		// "-e", "-o", resultsDir,
		"-l", resultsDir + "_result.tlf",
		"-JPROTOCOL=" + url.Scheme,
		"-JSERVER_URL=" + url.Hostname(),
		"-JDT_LTN=" + LTN,
		"-JVUCount=" + strconv.Itoa(workload.VUser),
		"-JLoopCount=" + strconv.Itoa(workload.LoopCount),
		"-JCHECK_PATH=" + derivePath(url),
		"-JSERVER_PORT=" + derivePort(url),
		"-JThinkTime=" + strconv.Itoa(workload.ThinkTime)})

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

	maxAcceptedErrors := float64(workload.AcceptedErrorRate) * float64(runs)
	if errorCount > int(maxAcceptedErrors) {
		logger.Debug(fmt.Sprintf("jmeter test failed because we got a too high error rate of %.2f.", float64(errorCount)/float64(runs)) + ". " + testInfo.ToString())
		return false, nil
	}

	if workload.AvgRtValidation > 0 && avg > workload.AvgRtValidation {
		logger.Debug(fmt.Sprintf("Avg rt validation failed because we got an avg rt of %d", workload.AvgRtValidation) + ". " + testInfo.ToString())
		return false, nil
	}

	logger.Debug("Successfully executed JMeter test. " + testInfo.ToString())
	return true, nil
}

func derivePort(url *url.URL) string {
	if url.Port() != "" {
		return url.Port()
	}
	switch strings.ToLower(url.Scheme) {
	case "http":
		return "80"
	case "https":
		return "443"
	}
	return ""
}

func derivePath(url *url.URL) string {
	if url.Path != "" {
		return url.Path
	}
	return "/health"
}

func getLastOccurence(vs []string, prefix string) string {
	for i := len(vs) - 1; i >= 0; i-- {
		if strings.HasPrefix(vs[i], prefix) {
			return vs[i]
		}
	}
	return ""
}
