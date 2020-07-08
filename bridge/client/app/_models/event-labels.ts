import {EventTypes} from "./event-types";

export const EVENT_LABELS = {
  [EventTypes.SERVICE_CREATE]: "Service create",
  [EventTypes.CONFIGURATION_CHANGE]: "Configuration change",
  [EventTypes.CONFIGURE_MONITORING]: "Configure monitoring",
  [EventTypes.DEPLOYMENT_FINISHED]: "Deployment finished",
  [EventTypes.TESTS_FINISHED]: "Tests finished",
  [EventTypes.START_EVALUATION]: "Start evaluation",
  [EventTypes.EVALUATION_DONE]: "Evaluation done",
  [EventTypes.START_SLI_RETRIEVAL]: "Start SLI retrieval",
  [EventTypes.SLI_RETRIEVAL_DONE]: "SLI retrieval done",
  [EventTypes.DONE]: "Done",
  [EventTypes.PROBLEM_OPEN]: "Problem open",
  [EventTypes.PROBLEM_DETECTED]: "Problem detected",
  [EventTypes.PROBLEM_RESOLVED]: "Problem resolved",
  [EventTypes.PROBLEM_CLOSED]: "Problem closed",
  [EventTypes.APPROVAL_TRIGGERED]: "Approval triggered",
  [EventTypes.APPROVAL_FINISHED]: "Approval finished"
};
