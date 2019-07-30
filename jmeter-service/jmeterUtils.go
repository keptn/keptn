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

func executeJMeter(testInfo string, scriptName string, resultsDir string, serverURL string, serverPort int, checkPath string, vuCount int,
	loopCount int, thinkTime int, LTN string, funcValidation bool, avgRtValidation int, logger *keptnutils.Logger) (bool, error) {

	os.RemoveAll(resultsDir)
	os.MkdirAll(resultsDir, 0644)

	testInfo = testInfo + ", scriptName: " + scriptName + ", serverURL: " + serverURL
	logger.Debug("Starting JMeter test. " + testInfo)
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

	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	// Parse result
	summary := getLastOccurence(strings.Split(res, "\n"), "summary =")
	if summary == "" {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo)
	}

	space := regexp.MustCompile(`\s+`)
	splits := strings.Split(space.ReplaceAllString(summary, " "), " ")
	runs, err := strconv.Atoi(splits[2])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo)
	}

	errorCount, err := strconv.Atoi(splits[14])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo)
	}

	avg, err := strconv.Atoi(splits[8])
	if err != nil {
		return false, errors.New("Cannot parse jmeter-result. " + testInfo)
	}

	if funcValidation && errorCount > 0 {
		logger.Debug(fmt.Sprintf("Function validation failed because we got %d errors.", errorCount) + ". " + testInfo)
		return false, nil
	}

	maxAcceptedErrors := maxAcceptedErrorRate * float64(runs)
	if errorCount > int(maxAcceptedErrors) {
		logger.Debug(fmt.Sprintf("jmeter test failed because we got a too high error rate of %.2f.", float64(errorCount)/float64(runs)) + ". " + testInfo)
		return false, nil
	}

	if avgRtValidation > 0 && avg > avgRtValidation {
		logger.Debug(fmt.Sprintf("Avg rt validation failed because we got an avg rt of %d", avgRtValidation) + ". " + testInfo)
		return false, nil
	}

	logger.Debug("Successfully executed JMeter test. " + testInfo)
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
