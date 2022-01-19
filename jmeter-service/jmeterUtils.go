package main

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/keptn/go-utils/pkg/common/fileutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"

	keptnutils "github.com/keptn/go-utils/pkg/lib"
)

// JMeterConfigDirectory defines the default jmeter config directory
const JMeterConfigDirectory = "/jmeter"

// TestInfo contains information about which test to execute
type TestInfo struct {
	Project           string
	Stage             string
	Service           string
	TestStrategy      string
	Context           string
	TriggeredID       string
	TestTriggeredData v0_2_0.TestTriggeredEventData
	ServiceURL        *url.URL
}

func (ti TestInfo) String() string {
	return fmt.Sprintf("Project: %s, Service: %s, Stage: %s, TestStrategy: %s, Context: %s", ti.Project, ti.Service, ti.Stage, ti.TestStrategy, ti.Context)
}

// shouldRemoveTempFiles Returns true if temp files should be removed. This is default - but can be changed through env variable DEBUG_KEEP_TEMP_FILES == true
func shouldRemoveTempFiles() bool {
	debugFlag := os.Getenv("DEBUG_KEEP_TEMP_FILES")
	return strings.Compare(debugFlag, "true") != 0
}

// createJMeterCLIArguments create the base arguments for the JMeter call
func createJMeterCLIArguments(workload *Workload, url *url.URL, resultsDir string, loadTestName string) []string {
	return []string{"-n", "-t", workload.Script,
		// "-e", "-o", resultsDir,
		"-l", resultsDir + "_result.tlf",
		"-JPROTOCOL=" + url.Scheme,
		"-JSERVER_PROTOCOL=" + url.Scheme,
		"-JSERVER_URL=" + url.Hostname(),
		"-JDT_LTN=" + loadTestName,
		"-JVUCount=" + strconv.Itoa(workload.VUser),
		"-JLoopCount=" + strconv.Itoa(workload.LoopCount),
		"-JCHECK_PATH=" + derivePath(url),
		"-JSERVER_PORT=" + derivePort(url),
		"-JThinkTime=" + strconv.Itoa(workload.ThinkTime)}
}

// addJMeterCommandLineArguments returns additional JMeter Command Line Parameters including additional params passed to the JMeter script
func addJMeterCommandLineArguments(testInfo TestInfo, initialList []string) []string {
	dtTenant := fmt.Sprintf("-JDT_TENANT=%s", os.Getenv("DT_TENANT"))
	dtAPIToken := fmt.Sprintf("-JDT_API_TOKEN=%s", os.Getenv("DT_API_TOKEN"))

	keptnProject := fmt.Sprintf("-JKEPTN_PROJECT=%s", testInfo.Project)
	keptnStage := fmt.Sprintf("-JKEPTN_STAGE=%s", testInfo.Stage)
	keptnService := fmt.Sprintf("-JKEPTN_SERVICE=%s", testInfo.Service)
	keptnTestStrategy := fmt.Sprintf("-JKEPTN_TESTSTRATEGY=%s", testInfo.TestStrategy)

	return append(initialList, dtTenant, dtAPIToken, keptnProject, keptnStage, keptnService, keptnTestStrategy)
}

// parseJMeterResult parses the output of the JMEter test and returns true or false
func parseJMeterResult(jmeterCommandResult string, testInfo TestInfo, workload *Workload, funcValidation bool) (bool, error) {
	summary := getLastOccurrence(strings.Split(jmeterCommandResult, "\n"), "summary =")
	if summary == "" {
		return false, fmt.Errorf("cannot parse jmeter-result. %v", testInfo)
	}

	space := regexp.MustCompile(`\s+`)
	splits := strings.Split(space.ReplaceAllString(summary, " "), " ")
	runs, err := strconv.Atoi(splits[2])
	if err != nil {
		return false, fmt.Errorf("cannot parse jmeter-result. %v", testInfo)
	}

	errorCount, err := strconv.Atoi(splits[14])
	if err != nil {
		return false, fmt.Errorf("cannot parse jmeter-result. %v", testInfo)
	}

	avg, err := strconv.Atoi(splits[8])
	if err != nil {
		return false, fmt.Errorf("cannot parse jmeter-result. %v", testInfo)
	}

	if funcValidation && errorCount > 0 {
		logger.Debugf("Function validation failed because we got %d errors. %v", errorCount, testInfo)
		return false, nil
	}

	maxAcceptedErrors := float64(workload.AcceptedErrorRate) * float64(runs)
	if errorCount > int(maxAcceptedErrors) {
		logger.Debugf("Jmeter test failed because we got a too high error rate of %.2f. %v", float64(errorCount)/float64(runs), testInfo)
		return false, nil
	}

	if workload.AvgRtValidation > 0 && avg > workload.AvgRtValidation {
		logger.Debugf("Avg rt validation failed because we got an avg rt of %d. %v", workload.AvgRtValidation, testInfo)
		return false, nil
	}

	return true, nil
}

// executeJMeter executes the actual JMeter script
// Step 1: Downloads all resources from the jmeter subfolder in the local container in a temporary folder and validates the referenced jmeter file was there
// Step 2: Executes the JMeter script that is referenced in the workload definition
// Step 3: Parses the response after JMeter execution is done
// Step 4: Removes the temporary folder
//
// Parameters:
// testInfo: information about the test, e.g: project, stage, service
// workload: jmeter.conf.yaml details
// resultsDir: resultsDir output
// url: the full server url. It gets parsed and then passed via JMeter properties SERVER_URL, SERVER_PORT, PROTOCOL, SERVER_PROTOCOL and CHECK_PATH
// LTN: will be passed as DT_LTN
// funcValidation: if true the function returns false if there were any errors detected during test execution
//
// Return:
// Status: true or false
// Error: error details if status was false
func executeJMeter(testInfo TestInfo, workload *Workload, resultsDir string, url *url.URL, loadTestName string, funcValidation bool) (bool, error) {
	if err := createDir(resultsDir); err != nil {
		return false, err
	}
	// Step 1: Lets download all files that match /jmeter/ into a local temp directory
	localTempDir := testInfo.Context
	if err := createDir(localTempDir); err != nil {
		return false, err
	}
	primaryScriptDownloaded, downloadedFileCount, err := DownloadAndStoreResources(testInfo.Project, testInfo.Stage, testInfo.Service, JMeterConfigDirectory, workload.Script, localTempDir)
	if err != nil {
		if errors.Is(err, ErrPrimaryFileNotAvailable) {
			// if no .jmx file is available -> skip the tests
			logger.Debug("skipping test execution because " + workload.Script + " not found on service, stage or project level.")
			return true, nil
		}
		err = fmt.Errorf("JMeter could not fetch the script %s for service %s in stage %s of project %s: %w",
			workload.Script,
			testInfo.Service,
			testInfo.Stage,
			testInfo.Project,
			err)
		return false, err
	}
	if downloadedFileCount == 0 {
		err = fmt.Errorf(`JMeter could not find default script %s for service %s in stage %s of project %s.
			Please upload that file to the jmeter subfolder or specify the correct script file in your
			jmeter.conf.yaml`,
			workload.Script,
			testInfo.Service,
			testInfo.Stage,
			testInfo.Project)
		return false, err
	}
	if !primaryScriptDownloaded {
		err = fmt.Errorf("JMeter could not find the primary file %s for service %s in stage %s of project %s",
			workload.Script,
			testInfo.Service,
			testInfo.Stage,
			testInfo.Project)
		return false, err
	}
	// this flag allows us to control whether files should be removed or not
	removeTempFiles := shouldRemoveTempFiles()

	// Step 1a: Lets validate if the script that was referenced in the workload was downloaded
	mainScriptFileName := localTempDir + "/" + workload.Script
	if !fileutils.FileExists(mainScriptFileName) {
		err = fmt.Errorf("JMeter script %s could not be found locally at %s for service %s in stage %s of project %s",
			workload.Script,
			mainScriptFileName,
			testInfo.Service,
			testInfo.Stage,
			testInfo.Project)
		if removeTempFiles {
			if err := os.RemoveAll(localTempDir); err != nil {
				return false, err
			}
		}
		return false, err
	}
	// Step 2: Lets execute the script - but be aware that we launch jmeter from the localTempDir as a working directory!
	jMeterCommandLineArgs := addJMeterCommandLineArguments(testInfo, createJMeterCLIArguments(workload, url, resultsDir, loadTestName))
	jmeterCommandResult, err := keptnutils.ExecuteCommandInDirectory("jmeter", jMeterCommandLineArgs, localTempDir)
	if err != nil {
		logger.Error(err.Error())
		return false, err
	}
	// now lets remove all downloaded files
	if removeTempFiles {
		if err := os.RemoveAll(localTempDir); err != nil {
			return false, err
		}
	}
	// Step 3: Parse result and lets analyze the result
	result, err := parseJMeterResult(jmeterCommandResult, testInfo, workload, funcValidation)
	if result && err != nil {
		logger.Debugf("Successfully executed JMeter test: %v", testInfo)
	} else {
		logger.Errorf("Successfully executed JMeter test: %v", testInfo)
	}
	return result, err
}

func createDir(dir string) error {
	if err := os.RemoveAll(dir); err != nil {
		return err
	}
	return os.MkdirAll(dir, 0644)
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

func getLastOccurrence(vs []string, prefix string) string {
	for i := len(vs) - 1; i >= 0; i-- {
		if strings.HasPrefix(vs[i], prefix) {
			return vs[i]
		}
	}
	return ""
}
