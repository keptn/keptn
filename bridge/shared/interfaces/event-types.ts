export enum EventTypes {
  SERVICE_CREATE = 'sh.keptn.internal.event.service.create',
  CONFIGURATION_CHANGE = 'sh.keptn.event.configuration.change',
  CONFIGURE_MONITORING = 'sh.keptn.event.monitoring.configure',
  DEPLOYMENT_TRIGGERED = 'sh.keptn.event.deployment.triggered',
  DEPLOYMENT_STARTED = 'sh.keptn.event.deployment.started',
  DEPLOYMENT_FINISHED = 'sh.keptn.event.deployment.finished',
  TESTS_FINISHED = 'sh.keptn.events.tests-finished',
  START_EVALUATION = 'sh.keptn.event.start-evaluation',
  EVALUATION_TRIGGERED = 'sh.keptn.event.evaluation.triggered',
  EVALUATION_TRIGGERED_SUFFIX = '.evaluation.triggered',
  EVALUATION_STARTED = 'sh.keptn.event.evaluation.started',
  EVALUATION_FINISHED = 'sh.keptn.event.evaluation.finished',
  EVALUATION_INVALIDATED = 'sh.keptn.event.evaluation.invalidated',
  START_SLI_RETRIEVAL = 'sh.keptn.internal.event.get-sli',
  SLI_RETRIEVAL_DONE = 'sh.keptn.internal.event.get-sli.done',
  DONE = 'sh.keptn.events.done',
  PROBLEM_OPEN = 'sh.keptn.event.problem.open',
  PROBLEM_DETECTED = 'sh.keptn.events.problem',
  PROBLEM_RESOLVED = 'sh.keptn.events.problem.resolved',
  PROBLEM_CLOSED = 'sh.keptn.event.problem.close',
  APPROVAL_TRIGGERED = 'sh.keptn.event.approval.triggered',
  APPROVAL_STARTED = 'sh.keptn.event.approval.started',
  APPROVAL_FINISHED = 'sh.keptn.event.approval.finished',
  REMEDIATION_TRIGGERED_SUFFIX = 'remediation.triggered',
  REMEDIATION_TRIGGERED = 'sh.keptn.event.remediation.triggered',
  REMEDIATION_STATUS_CHANGED = 'sh.keptn.event.remediation.status.changed',
  REMEDIATION_FINISHED_SUFFIX = 'remediation.finished',
  REMEDIATION_FINISHED = 'sh.keptn.event.remediation.finished',
  ACTION_TRIGGERED = 'sh.keptn.event.action.triggered',
  ACTION_STARTED = 'sh.keptn.event.action.started',
  ACTION_FINISHED = 'sh.keptn.event.action.finished',
  PREFIX = 'sh.keptn.event.',
  PREFIX_SHORT = 'sh.keptn',
}
