import {EventTypes} from "./event-types";
import {ApprovalStates} from "./approval-states";

export const EVENT_ICONS = {
  [EventTypes.CONFIGURATION_CHANGE]: "duplicate",
  [EventTypes.DEPLOYMENT_FINISHED]: "deploy",
  [EventTypes.TESTS_FINISHED]: "perfromance-health",
  [EventTypes.START_EVALUATION]: "traffic-light",
  [EventTypes.EVALUATION_FINISHED]: "traffic-light",
  [EventTypes.START_SLI_RETRIEVAL]: "collector",
  [EventTypes.SLI_RETRIEVAL_DONE]: "collector",
  [EventTypes.PROBLEM_OPEN]: "criticalevent",
  [EventTypes.PROBLEM_DETECTED]: "criticalevent",
  [EventTypes.PROBLEM_CLOSED]: "applicationhealth",
  [EventTypes.APPROVAL_TRIGGERED]: "unknown",
  [EventTypes.APPROVAL_FINISHED]: {
    [ApprovalStates.APPROVED]: "checkmark",
    [ApprovalStates.DECLINED]: "abort",
  }
};
