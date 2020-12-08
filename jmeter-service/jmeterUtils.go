package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

const maxAcceptedErrorRate = 0.1
const JMeterConfigDirectory = "/jmeter"

// TestInfo contains information about which test to execute
type TestInfo struct {
	Project      string
	Stage        string
	Service      string
	TestStrategy string
	Context      string
}

// Returns true if temp files should be removed. This is default - but can be changed through env variable DEBUG_KEEP_TEMP_FILES == true
func DoRemoveTempFiles() bool {
	debugFlag := os.Getenv("DEBUG_KEEP_TEMP_FILES")
	if strings.Compare(debugFlag, "true") == 0 {
		return false
	}

	return true
}

// ToString returns a string representation of a TestInfo object
func (ti *TestInfo) ToString() string {
	return fmt.Sprintf("Project: %s, Service: %s, Stage: %s, TestStrategy: %s, Context: %s", ti.Project, ti.Service, ti.Stage, ti.TestStrategy, ti.Context)
}

/**
 * Returns additoinal JMeter Command Line Parameters including additional params passed to the JMeter script
 */
func addJMeterCommandLineArguments(testInfo *TestInfo, initialList []string) []string {
	dtTenant := fmt.Sprintf("-JDT_TENANT=%s", os.Getenv("DT_TENANT"))
	dtAPIToken := fmt.Sprintf("-JDT_API_TOKEN=%s", os.Getenv("DT_API_TOKEN"))

	keptnProject := fmt.Sprintf("-JKEPTN_PROJECT=%s", testInfo.Project)
	keptnStage := fmt.Sprintf("-JKEPTN_STAGE=%s", testInfo.Stage)
	keptnService := fmt.Sprintf("-JKEPTN_SERVICE=%s", testInfo.Service)
	keptnTestStrategy := fmt.Sprintf("-JKEPTN_TESTSTRATEGY=%s", testInfo.TestStrategy)

	return append(initialList, dtTenant, dtAPIToken, keptnProject, keptnStage, keptnService, keptnTestStrategy)
}

/**
 * Parses the output of the JMEter test and returns true or false
 */
func parseJMeterResult(jmeterCommandResult string, testInfo *TestInfo, workload *Workload, funcValidation bool, logger *keptnutils.Logger) (bool, error) {

	logger.Debug(jmeterCommandResult)

	summary := getLastOccurence(strings.Split(jmeterCommandResult, "\n"), "summary =")
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

	return true, nil
}

/**
 * Executes the actual JMeter script
 * Step 1: Downloads all resources from the jmeter subfolder in the local container in a temporary folder and validates the referenced jmeter file was there
 * Step 2: Executes the JMeter script that is referenced in the workload definition
 * Step 3: Parses the response after JMeter execution is done
 * Step 4: Removes the temporary folder
 *
 * Parameters:
 * testInfo: information about the test, e.g: project, stage, service
 * workload: jmeter.conf.yaml details
 * resultsDir: resultsDir output
 * url: the full server url. It gets parsed and then passed via JMeter properties SERVER_URL, SERVER_PORT, PROTOCOL, SERVER_PROTOCAL and CHECK_PATH
 * LTN: will be passed as DT_LTN
 * funcValidation: if true the function returns false if there were any errors detected during test execution
 *
 * Return:
 * Status: true or false
 * Error: error details if status was false
 */
func executeJMeter(testInfo *TestInfo, workload *Workload, resultsDir string, url *url.URL, LTN string, funcValidation bool, logger *keptnutils.Logger) (bool, error) {
	os.RemoveAll(resultsDir)
	os.MkdirAll(resultsDir, 0644)

	// Step 1: Lets download all files that match /jmeter/ into a local temp directory
	// Due to current limitations of the REST API we also fall-back and always load a specific file referenced in workload on service, stage or project level
	// Implementing https://github.com/keptn/keptn/issues/2756
	localTempDir := testInfo.Context
	os.RemoveAll(localTempDir)
	os.MkdirAll(localTempDir, 0644)
	fileMatchPattern := JMeterConfigDirectory
	primaryScriptDownloaded, downloadedFileCount, err := GetAllKeptnResources(testInfo.Project, testInfo.Stage, testInfo.Service, true, fileMatchPattern, workload.Script, localTempDir, logger)

	if err != nil {
		err = fmt.Errorf("Error loading /jmeter/* files for %s.%s.%s: %s", testInfo.Project, testInfo.Stage, testInfo.Service, err.Error())
		return false, err
	}
	if downloadedFileCount == 0 {
		err = fmt.Errorf("No files found in /jmeter/* for %s.%s.%s", testInfo.Project, testInfo.Stage, testInfo.Service)
		return false, err
	}
	if !primaryScriptDownloaded {
		err = fmt.Errorf("Primary file %s was not found for %s.%s.%s", workload.Script, testInfo.Project, testInfo.Stage, testInfo.Service)
		return false, err
	}

	// this flag allows us to control whether files should be removed or not
	removeTempFiles := DoRemoveTempFiles()

	// Step 1a: Lets validate if the script that was referenced in the workload was downloaded
	mainScriptFileName := localTempDir + "/" + workload.Script
	if !FileExists(mainScriptFileName) {
		err = fmt.Errorf("JMeter script %s not found locally at %s for %s.%s.%s", workload.Script, mainScriptFileName, testInfo.Project, testInfo.Stage, testInfo.Service)
		if removeTempFiles {
			os.RemoveAll(localTempDir)
		}
		return false, err
	}

	// Step 2: Lets execute the script - but be aware that we launch jmeter from the localTempDir as a working directory!
	testInfoStr := testInfo.ToString() + ", scriptName: " + mainScriptFileName + ", serverURL: " + url.String()
	logger.Debug("Starting JMeter test. " + testInfoStr)

	jMeterCommandLineArgs := []string{"-n", "-t", workload.Script,
		// "-e", "-o", resultsDir,
		"-l", resultsDir + "_result.tlf",
		"-JPROTOCOL=" + url.Scheme,
		"-JSERVER_PROTOCOL=" + url.Scheme,
		"-JSERVER_URL=" + url.Hostname(),
		"-JDT_LTN=" + LTN,
		"-JVUCount=" + strconv.Itoa(workload.VUser),
		"-JLoopCount=" + strconv.Itoa(workload.LoopCount),
		"-JCHECK_PATH=" + derivePath(url),
		"-JSERVER_PORT=" + derivePort(url),
		"-JThinkTime=" + strconv.Itoa(workload.ThinkTime)}

	jMeterCommandLineArgs = addJMeterCommandLineArguments(testInfo, jMeterCommandLineArgs)
	jmeterCommandResult, err := keptnutils.ExecuteCommandInDirectory("jmeter", jMeterCommandLineArgs, localTempDir)

	// now lets remove all downloaded files
	if removeTempFiles {
		os.RemoveAll(localTempDir)
	}

	// Step 3: Parse result
	// and lets analyze the result
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}

	result, err := parseJMeterResult(jmeterCommandResult, testInfo, workload, funcValidation, logger)
	if result && err != nil {
		logger.Debug("Successfully executed JMeter test. " + testInfo.ToString())
	} else {
		logger.Error("Successfully executed JMeter test. " + testInfo.ToString())
	}

	return result, err
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
