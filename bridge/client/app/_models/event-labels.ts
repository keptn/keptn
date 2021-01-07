import {EventTypes} from "./event-types";
import {ApprovalStates} from "./approval-states";

export const EVENT_LABELS = {
  [EventTypes.SERVICE_CREATE]: "Service created",
  [EventTypes.CONFIGURATION_CHANGE]: "Configuration changed",
  [EventTypes.CONFIGURE_MONITORING]: "Configure monitoring",
  [EventTypes.DEPLOYMENT_FINISHED]: "Deployment finished",
  [EventTypes.TESTS_FINISHED]: "Tests finished",
  [EventTypes.START_EVALUATION]: "Evaluation started",
  [EventTypes.EVALUATION_FINISHED]: "Evaluation finished",
  [EventTypes.EVALUATION_INVALIDATED]: "Evaluation invalidated",
  [EventTypes.START_SLI_RETRIEVAL]: "SLI retrieval started",
  [EventTypes.SLI_RETRIEVAL_DONE]: "SLI retrieval done",
  [EventTypes.DONE]: "Done",
  [EventTypes.PROBLEM_OPEN]: "Problem opened",
  [EventTypes.PROBLEM_DETECTED]: "Problem detected",
  [EventTypes.PROBLEM_RESOLVED]: "Problem resolved",
  [EventTypes.PROBLEM_CLOSED]: "Problem closed",
  [EventTypes.APPROVAL_TRIGGERED]: "Approval triggered",
  [EventTypes.APPROVAL_FINISHED]: {
    [ApprovalStates.APPROVED]: "Approval finished",
    [ApprovalStates.DECLINED]: "Approval declined"
  },
  [EventTypes.REMEDIATION_TRIGGERED]: 'Remediation triggered',
  [EventTypes.REMEDIATION_STATUS_CHANGED]: 'Remediation status changed',
  [EventTypes.REMEDIATION_FINISHED]: 'Remediation finished',
  [EventTypes.ACTION_TRIGGERED]: 'Action triggered',
  [EventTypes.ACTION_STARTED]: 'Action started',
  [EventTypes.ACTION_FINISHED]: 'Action finished',
};
