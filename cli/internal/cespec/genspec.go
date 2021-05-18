package cespec

import (
	"encoding/json"
	"fmt"
	"github.com/alecthomas/jsonschema"
	keptn "github.com/keptn/go-utils/pkg/lib"
	keptnv2 "github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"os"
)

func createDataJSONSchemaSection(md *MarkDown, eventType string, data interface{}) {
	md.Title("Data Json Schema", 5)
	md.WriteLineBreak()
	md.Writeln("<details><summary>Json Schema of " + eventType + "</summary>")
	md.Writeln("<p>")
	md.WriteLineBreak()
	md.CodeBlock(toJSONSchema(data), "json")
	md.Writeln("</p>")
	md.Writeln("</details>")
	md.WriteLineBreak()
}

func createExampleSection(md *MarkDown, eventType string, data interface{}) {
	md.Title("Example Cloud Event", 5)
	md.WriteLineBreak()
	md.CodeBlock(toJSON(ce(eventType, data)), "json")
	md.WriteLineBreak()
}
func createUpLink(md *MarkDown) {
	md.UpLink()
}

func createSection(md *MarkDown, title string, eventType string, data interface{}) {
	md.Title(title, 4)
	md.Title("Type", 5)
	md.Writeln(eventType)
	createDataJSONSchemaSection(md, eventType, data)
	createExampleSection(md, eventType, data)
	createUpLink(md)
}

func createSectionTitle(md *MarkDown, title string) {
	md.Title(title, 3)
}

// Generate produces a cloudevents.md file in the provided outputDir path
func Generate(outputDir string) {

	md := NewMarkDown()

	md.Title("Keptn CloudEvents", 1)

	md.Bullet().Link("Project", "#project")
	md.Bullet().Link("Service", "#service")
	md.Bullet().Link("Approval", "#approval")
	md.Bullet().Link("Deployment", "#deployment")
	md.Bullet().Link("Test", "#test")
	md.Bullet().Link("Evaluation", "#evaluation")
	md.Bullet().Link("Release", "#release")
	md.Bullet().Link("Get-Action", "#get-action")
	md.Bullet().Link("Action", "#action")
	md.Bullet().Link("Get-SLI", "#get-sli")
	md.Bullet().Link("Monitoring", "#monitoring")
	md.Bullet().Link("Problem", "#problem")

	md.Writeln("---")

	md.Writeln("All Keptn events conform to the CloudEvents spec in [version 1.0](https://github.com/cloudevents/spec/blob/v1.0/spec.md). The CloudEvents specification is a vendor-neutral specification for defining the format of event data.")
	md.WriteLineBreak()
	md.Writeln("In Keptn, events have a payload structure as follows (*Note:* The `triggeredid` is not contained in events of type `triggered` mentioned below):")

	createCEStructureCodeBlock(md)

	md.Title("Type", 2)
	md.Writeln("In Keptn, events follow two different formats of event types. One is related to the overall status of a **task sequence execution**, while the other format is related to the execution of a **certain task within a sequence**.")

	md.Title("Task sequence events", 3)
	md.Writeln("The event type of a Keptn CloudEvent concerning the overall state of a task sequence has the following format:")
	md.WriteLineBreak()
	md.Bullet().Writeln("`sh.keptn.event.[stage].[task sequence].[event status]` - For events concerning the execution of a task sequence")
	md.WriteLineBreak()
	md.Writeln("As indicated by the brackets, the event type is defined by a **stage**, **task sequence** and the **event status**.")
	md.Bullet().Writeln("The task sequence is declared in the [Shipyard](https://github.com/keptn/spec/blob/master/shipyard.md) of a project.")
	md.Bullet().Writeln("The kinds of event states are defined with: `triggered` and `finished`")
	md.WriteLineBreak()
	md.Writeln("For example, if a task sequence with the name `delivery` in the stage `hardening` should be executed, it has to be triggered by sending an event with the type `sh.keptn.event.hardening.delivery.triggered`. Once the `delivery` sequence is completed, a `sh.keptn.event.hardening.delivery.finished` event will be sent to indicate the completion of the task sequence.")
	md.WriteLineBreak()
	md.WriteLineBreak()

	md.Title("Task events", 3)
	md.Writeln("The event type of a Keptn CloudEvent concerning the execution of a certain task within a task sequence has the following format:")
	md.WriteLineBreak()
	md.Bullet().Writeln("`sh.keptn.event.[task].[event status]` - For events concerning the execution of a certain task within a task sequence")
	md.WriteLineBreak()
	md.Writeln("As indicated by the brackets, the event type is defined by a **task** and the **event status**.")
	md.Bullet().Writeln("The task is declared in the [Shipyard](https://github.com/keptn/spec/blob/master/shipyard.md) of a project. For example, a Shipyard can contain tasks like: `deployment`, `test`, or `evaluation`. Consequently, the event type for a `deployment` task would be `sh.keptn.event.deployment.[event status]`")
	md.Bullet().Writeln("The kinds of event states are defined with: `triggered`, `started`, `status.changed`, and `finished` (`status.changed` is optional)")
	md.WriteLineBreak()

	md.Writeln("By combining the *task* and *event status* for the `deployment` task, the event types are:")
	md.WriteLineBreak()

	md.Bullet().Writeln("`sh.keptn.event.deployment.triggered`")
	md.Bullet().Writeln("`sh.keptn.event.deployment.started`")
	md.Bullet().Writeln("`sh.keptn.event.deployment.status.changed`")
	md.Bullet().Writeln("`sh.keptn.event.deployment.finished`")

	md.Title("Data", 2)
	md.Writeln("The data block of a Keptn CloudEvent carries the Keptn Payload of a specific event and contains the properties:")
	md.Bullet().Writeln("labels")
	md.Bullet().Writeln("message")
	md.Bullet().Writeln("project")
	md.Bullet().Writeln("service")
	md.Bullet().Writeln("stage")
	md.Bullet().Writeln(fmt.Sprintf("status: indicates whether the service executing the task was able to perform the task without any unexpected errors. Possible values are `%s`, `%s`,or `%s`", keptnv2.StatusSucceeded, keptnv2.StatusErrored, keptnv2.StatusUnknown))
	md.Bullet().Writeln(fmt.Sprintf("result: indicates the result of a successful task execution without unexpected problems (i.e. status = `%s`), such as the result of an evaluation, or a test execution. Possible values are `%s`, `%s`, or `%s`", keptnv2.StatusSucceeded, keptnv2.ResultPass, keptnv2.ResultWarning, keptnv2.ResultFailed))
	md.Bullet().Writeln("*[task]*")

	md.WriteLineBreak()
	md.Writeln("Like the task property in the event type, the task property in the data block depends on the task declaration in the Shipyard. Based on the example of a `deployment` task, the data block contains a `deployment` property of type object. Hence, any payload can be added to this `deployment` property")
	md.WriteLineBreak()

	md.Writeln("In the following each data block is described and an example of a CloudEvent containing the data block is given.")
	md.MultiBr(2)

	createSectionTitle(md, "Project")
	createSection(md, "Project Create Triggered", keptnv2.GetTriggeredEventType(keptnv2.ProjectCreateTaskName), projectCreateData)
	createSection(md, "Project Create Started", keptnv2.GetStartedEventType(keptnv2.ProjectCreateTaskName), projectCreateStartedEventData)
	createSection(md, "Project Create Finished", keptnv2.GetFinishedEventType(keptnv2.ProjectCreateTaskName), projectCreateFinishedEventData)

	createSectionTitle(md, "Service")
	createSection(md, "Service Create Started", keptnv2.GetStartedEventType(keptnv2.ServiceCreateTaskName), serviceCreateStartedEventData)
	createSection(md, "Service Create Status Changed", keptnv2.GetStatusChangedEventType(keptnv2.ServiceCreateTaskName), serviceCreateStatusChangesData)
	createSection(md, "Service Create Finished", keptnv2.GetFinishedEventType(keptnv2.ServiceCreateTaskName), serviceCreateFinishedEventData)

	createSectionTitle(md, "Approval")
	createSection(md, "Approval Triggered", keptnv2.GetTriggeredEventType(keptnv2.ApprovalTaskName), approvalTriggeredEventData)
	createSection(md, "Approval Started", keptnv2.GetStartedEventType(keptnv2.ApprovalTaskName), approvalStartedEventData)
	createSection(md, "Approval Status Changed", keptnv2.GetStatusChangedEventType(keptnv2.ApprovalTaskName), approvalStatusChangedEventData)
	createSection(md, "Approval Finished", keptnv2.GetFinishedEventType(keptnv2.ApprovalTaskName), approvalFinishedEventData)

	createSectionTitle(md, "Deployment")
	createSection(md, "Deployment Triggered", keptnv2.GetTriggeredEventType(keptnv2.DeploymentTaskName), deploymentTriggeredEventData)
	createSection(md, "Deployment Started", keptnv2.GetStartedEventType(keptnv2.DeploymentTaskName), deploymentStartedEventData)
	createSection(md, "Deployment Status Changed", keptnv2.GetStatusChangedEventType(keptnv2.DeploymentTaskName), deploymentStatusChangedEventData)
	createSection(md, "Deployment Finished", keptnv2.GetFinishedEventType(keptnv2.DeploymentTaskName), deploymentFinishedEventData)

	createSectionTitle(md, "Test")
	createSection(md, "Test Triggered", keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), testTriggeredEventData)
	createSection(md, "Test Started", keptnv2.GetStartedEventType(keptnv2.TestTaskName), testStartedEventData)
	createSection(md, "Test Status Changed", keptnv2.GetTriggeredEventType(keptnv2.TestTaskName), testStatusChangedEventData)
	createSection(md, "Test Finished", keptnv2.GetFinishedEventType(keptnv2.TestTaskName), testTestFinishedEventData)

	createSectionTitle(md, "Evaluation")
	createSection(md, "Evaluation Triggered", keptnv2.GetTriggeredEventType(keptnv2.EvaluationTaskName), evaluationTriggeredEventData)
	createSection(md, "Evaluation Started", keptnv2.GetStartedEventType(keptnv2.EvaluationTaskName), evaluationStartedEventData)
	createSection(md, "Evaluation Status Changed", keptnv2.GetStatusChangedEventType(keptnv2.EvaluationTaskName), evaluationStatusChangedEventData)
	createSection(md, "Evaluation Finished", keptnv2.GetFinishedEventType(keptnv2.EvaluationTaskName), evaluationFinishedEventData)
	createSection(md, "Evaluation Invalidated", keptnv2.GetInvalidatedEventType(keptnv2.EvaluationTaskName), evaluationInvalidatedEventData)

	createSectionTitle(md, "Release")
	createSection(md, "Release Triggered", keptnv2.GetTriggeredEventType(keptnv2.ReleaseTaskName), releaseTriggeredEventData)
	createSection(md, "Release Started", keptnv2.GetStartedEventType(keptnv2.ReleaseTaskName), releaseStartedEventData)
	createSection(md, "Release Status Changed", keptnv2.GetStatusChangedEventType(keptnv2.ReleaseTaskName), releaseStatusChangedEventData)
	createSection(md, "Release Finished", keptnv2.GetFinishedEventType(keptnv2.ReleaseTaskName), releaseFinishedEventData)

	createSectionTitle(md, "Get Action")
	createSection(md, "Get Action Triggered", keptnv2.GetTriggeredEventType(keptnv2.GetActionTaskName), getActionTriggeredEventData)
	createSection(md, "Get Action Started", keptnv2.GetStartedEventType(keptnv2.GetActionTaskName), getActionStartedEventData)
	createSection(md, "Get Action Finished", keptnv2.GetFinishedEventType(keptnv2.GetActionTaskName), getActionFinishedEventData)

	createSectionTitle(md, "Action")
	createSection(md, "Action Triggered", keptnv2.GetTriggeredEventType(keptnv2.ActionTaskName), actionTriggeredEventData)
	createSection(md, "Action Started", keptnv2.GetStartedEventType(keptnv2.ActionTaskName), actionStartedEventData)
	createSection(md, "Action Finished", keptnv2.GetFinishedEventType(keptnv2.ActionTaskName), actionFinishedEventData)

	createSectionTitle(md, "Get SLI")
	createSection(md, "Get SLI Triggered", keptnv2.GetTriggeredEventType(keptnv2.GetSLITaskName), getSLITriggeredEventData)
	createSection(md, "Get SLI Started", keptnv2.GetStartedEventType(keptnv2.GetSLITaskName), getSLIStartedEventData)
	createSection(md, "Get SLI Finished", keptnv2.GetFinishedEventType(keptnv2.GetSLITaskName), getSLIFinishedEventData)

	createSectionTitle(md, "Monitoring")
	createSection(md, "Configure Monitoring", keptn.ConfigureMonitoringEventType, configureMonitoringEventData)

	createSectionTitle(md, "Problem")
	createSection(md, "Problem", keptn.ProblemEventType, problemOpenEventData)
	fmt.Println(md.String())

	file, err := os.Create(outputDir + "/" + "cloudevents.md")
	check(err)
	defer file.Close()
	file.WriteString(md.String())
	file.Sync()

}

func createDataStructureCodeBlock(md *MarkDown) {
	md.CodeBlock(`"data": {
  "required": [
    "labels",
    "message",
    "project",
    "result",
    "service",
    "stage",
    "status",
    "[task]"
  ],
  "properties": {
    "labels": {
      "patternProperties": {
        ".*": {
          "type": "string"
        }
      },
      "type": "object"
    },
    "message": {
      "type": "string",
      "description": "A message from the last task"
    },
    "project": {
      "type": "string",
      "description": "The name of the project"
    },
    "result": {
      "type": "string",
      "description": "The result of the last task"
    },
    "service": {
      "type": "string",
      "description": "The name of the service"
    },
    "stage": {
      "type": "string",
      "description": "The name of the stage"
    },
    "status": {
      "type": "string",
      "description": "The status of the last task"
    },
    "[task]": {
      "type": "object"
    }
  },
  "additionalProperties": false,
  "type": "object"
}`, "json")
}

func createCEStructureCodeBlock(md *MarkDown) {
	md.CodeBlock(`"sh.keptn.event": {
  "required": [
    "data",
    "id",
    "shkeptncontext",
    "source",
    "specversion",
    "time",
    "triggeredid",
    "type"
  ],
  "properties": {
    "data": {
      "type": ["object", "string"],
      "description": "The Keptn event payload depending on the type."
    },
    "id": {
      "type": "string",
      "minLength": 1,
      "description": "Unique UUID of the Keptn event"
    },
    "shkeptncontext": {
      "format": "uuid",
      "type": "string",
      "description": "Unique UUID value that connects various events together"
    },
    "source": {
      "format": "uri-reference",
      "type": "string",
      "minLength": 1,
      "description": "URL to service implementation in Keptn code repo"
    },
    "specversion": {
      "type": "string",
      "minLength": 1,
      "description": "The version of the CloudEvents specification",
      "value": "1.0"
    },
    "shkeptnspecversion": {
      "type": "string",
      "minLength": 1,
      "description": "The version of the Keptn specification",
      "value": "0.2.0"
    },
    "time": {
      "format": "date-time",
      "type": "string",
      "description": "Timestamp of when the event happened"
    },
    "triggeredid": {
      "format": "uuid",
      "type": "string",
      "description": "The event ID that has triggered the step"
    },
    "type": {
      "type": "string",
      "minLength": 1,
      "description": "Type of the Keptn event"
    }
  },
  "additionalProperties": false,
  "type": "object"
}`, "json")
}

func toJSONSchema(j interface{}) string {
	schema := jsonschema.Reflect(j)
	schemaStr, _ := json.MarshalIndent(schema, "", "  ")
	return string(schemaStr)

}

func toJSON(j interface{}) string {
	jsonStr, _ := json.MarshalIndent(j, "", "  ")
	return string(jsonStr)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
