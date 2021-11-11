package main

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	logger "github.com/sirupsen/logrus"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

type TestRunner struct {
	eventSender *keptnv2.HTTPEventSender
}

func NewTestRunner(eventSender *keptnv2.HTTPEventSender) *TestRunner {
	return &TestRunner{eventSender}
}

func (tr *TestRunner) RunTests(testInfo TestInfo) error {
	if err := tr.sendTestsStartedEvent(testInfo); err != nil {
		logger.WithError(err).Error("Unable to send test '.started' event")
		return err
	}
	testStartedAt := time.Now()

	jmeterConfig, err := getJMeterConf(testInfo)
	if err != nil {
		return err
	}

	if err := tr.runHealthCheck(testInfo, testStartedAt, jmeterConfig); err != nil {
		return err
	}
	if err := tr.runTests(testInfo, testStartedAt, jmeterConfig); err != nil {
		return err
	}

	return nil
}

func (tr *TestRunner) runTests(testInfo TestInfo, testStartedAt time.Time, jmeterConf *JMeterConf) error {
	// now lets execute the workload based on the passed testStrategy
	res := false
	var testStrategy = strings.ToLower(testInfo.TestStrategy)

	if testStrategy == "" {
		// no testStrategy passed at all -> we just send a successful test finished event!
		logger.Info("No testStrategy specified therefore skipping further test execution and sending back success")
		res = true
	} else {
		// lets get the workload configuration for the test strategy
		teststrategyWorkload, err := getWorkload(jmeterConf, testStrategy)
		if teststrategyWorkload != nil {
			res, err = tr.runWorkload(testInfo, teststrategyWorkload)
			if err != nil {
				msg := fmt.Sprintf("could not run test workload: %s", err.Error())
				logger.Error(msg)
				if err := tr.sendErroredTestsFinishedEvent(testInfo, testStartedAt, msg); err != nil {
					logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
				}
				return nil
			} else {
				logger.Infof("Tests for %s with status = %s.%s", testStrategy, strconv.FormatBool(res), testInfo.ToString())
			}
		} else {
			// no workload for that test strategy!!
			res = false
			logger.Errorf("No workload definition found for testStrategy %s", testStrategy)
		}
	}

	// now lets send the test finished event
	msg := fmt.Sprintf("Tests for %s with status = %s.%s", testStrategy, strconv.FormatBool(res), testInfo.ToString())
	if !res {
		if err := tr.sendTestsFinishedEvent(testInfo, testStartedAt, msg, keptnv2.ResultFailed); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return nil
	}

	if err := tr.sendTestsFinishedEvent(testInfo, testStartedAt, msg, keptnv2.ResultPass); err != nil {
		logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
	}
	return nil
}

func (tr *TestRunner) runHealthCheck(testInfo TestInfo, testStartedAt time.Time, jmeterConf *JMeterConf) error {
	healthCheckWorkload, err := getWorkload(jmeterConf, TestStrategy_HealthCheck)
	if err != nil {
		return err
	}
	if healthCheckWorkload != nil {
		logger.Info("No Health Check test workload configuration found. Skipping Health Check")
		return nil
	}
	if err := checkEndpointAvailable(5*time.Second, testInfo.ServiceURL); err != nil {
		msg := fmt.Sprintf("Jmeter-service cannot reach URL %s: %s", testInfo.Service, err.Error())
		logger.Error(msg)
		if err := tr.sendTestsFinishedEvent(testInfo, testStartedAt, msg, keptnv2.ResultFailed); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return nil
	}
	res, err := tr.runWorkload(testInfo, healthCheckWorkload)
	if err != nil {
		msg := fmt.Sprintf("could not run test workload: %s", err.Error())
		logger.Error(msg)
		if err := tr.sendErroredTestsFinishedEvent(testInfo, testStartedAt, msg); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return nil
	}
	if !res {
		msg := fmt.Sprintf("Tests for %s with status = %s.%s", TestStrategy_HealthCheck, strconv.FormatBool(res), testInfo.ToString())
		if err := tr.sendTestsFinishedEvent(testInfo, testStartedAt, msg, keptnv2.ResultFailed); err != nil {
			logger.WithError(err).Errorf("Error sending test finished event for %s", testInfo.ToString())
		}
		return nil
	}
	logger.Info("Health Check test passed = " + strconv.FormatBool(res) + ". " + testInfo.ToString())
	return nil
}

func (tr *TestRunner) runWorkload(testInfo TestInfo, workload *Workload) (bool, error) {
	// for testStrategy functional we enforce a 0% error policy!
	breakOnFunctionalIssues := workload.TestStrategy == TestStrategy_Functional

	logger.Infof(
		"Running workload testStrategy=%s, vuser=%d, loopcount=%d, thinktime=%d, funcvalidation=%t, acceptederrors=%f, avgrtvalidation=%d, script=%s",
		workload.TestStrategy, workload.VUser, workload.LoopCount, workload.ThinkTime, breakOnFunctionalIssues, workload.AcceptedErrorRate, workload.AvgRtValidation, workload.Script)
	if runlocal {
		logger.Info("LOCALLY: not executing actual tests!")
		return true, nil
	}

	// the resultdirectory is unique as it contains context but also gives some human readable context such as teststrategy and service
	// this will also be used for TSN parameter
	resultDirectory := fmt.Sprintf("%s_%s_%s_%s_%s", testInfo.Project, testInfo.Service, testInfo.Stage, workload.TestStrategy, testInfo.Context)

	// lets first remove all potentially left over result files from previous runs -> we keep them between runs for troubleshooting though
	err := os.RemoveAll(resultDirectory)
	if err != nil {
		return false, err
	}

	err = os.RemoveAll(resultDirectory + "_result.tlf")
	if err != nil {
		return false, err
	}

	err = os.RemoveAll("output.txt")
	if err != nil {
		return false, err
	}
	return executeJMeter(testInfo, workload, resultDirectory, testInfo.ServiceURL, resultDirectory, breakOnFunctionalIssues)
}

func (tr *TestRunner) sendTestsStartedEvent(testInfo TestInfo) error {
	source, _ := url.Parse("jmeter-service")

	testStartedEventData := keptnv2.TestStartedEventData{}
	testStartedEventData.EventData = testInfo.TestTriggeredData.EventData
	testStartedEventData.Status = keptnv2.StatusSucceeded

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetStartedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", testInfo.Context)
	event.SetExtension("triggeredid", testInfo.TriggeredID)
	event.SetData(cloudevents.ApplicationJSON, testStartedEventData)

	if runlocal {
		log.Println("LOCALLY: Sending Event")
		return nil
	}
	return tr.eventSender.SendEvent(event)

}

func (tr *TestRunner) sendTestsFinishedEvent(testInfo TestInfo, startedAt time.Time, msg string, result keptnv2.ResultType) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testFinishedData.EventData = testInfo.TestTriggeredData.EventData
	// fill in timestamps
	testFinishedData.Test.Start = startedAt.Format(time.RFC3339)
	testFinishedData.Test.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = result
	testFinishedData.Status = keptnv2.StatusSucceeded
	testFinishedData.Message = msg

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", testInfo.Context)
	event.SetExtension("triggeredid", testInfo.TriggeredID)
	event.SetData(cloudevents.ApplicationJSON, testFinishedData)

	return tr.eventSender.SendEvent(event)
}

func (tr *TestRunner) sendErroredTestsFinishedEvent(testInfo TestInfo, startedAt time.Time, msg string) error {
	source, _ := url.Parse("jmeter-service")

	testFinishedData := keptnv2.TestFinishedEventData{}
	testFinishedData.EventData = testInfo.TestTriggeredData.EventData
	// fill in timestamps
	testFinishedData.Test.Start = startedAt.Format(time.RFC3339)
	testFinishedData.Test.End = time.Now().Format(time.RFC3339)
	// set test result
	testFinishedData.Result = keptnv2.ResultFailed
	testFinishedData.Status = keptnv2.StatusErrored
	testFinishedData.Message = msg

	event := cloudevents.NewEvent()
	event.SetType(keptnv2.GetFinishedEventType(keptnv2.TestTaskName))
	event.SetSource(source.String())
	event.SetDataContentType(cloudevents.ApplicationJSON)
	event.SetExtension("shkeptncontext", testInfo.Context)
	event.SetExtension("triggeredid", testInfo.TriggeredID)
	event.SetData(cloudevents.ApplicationJSON, testFinishedData)

	return tr.eventSender.SendEvent(event)
}
