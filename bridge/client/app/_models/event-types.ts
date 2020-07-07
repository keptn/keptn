export enum EventTypes {
  SERVICE_CREATE = 'sh.keptn.internal.event.service.create',
  CONFIGURATION_CHANGE = 'sh.keptn.event.configuration.change',
  CONFIGURE_MONITORING = 'sh.keptn.event.monitoring.configure',
  DEPLOYMENT_FINISHED = 'sh.keptn.events.deployment-finished',
  TESTS_FINISHED = 'sh.keptn.events.tests-finished',
  START_EVALUATION = 'sh.keptn.event.start-evaluation',
  EVALUATION_DONE = 'sh.keptn.events.evaluation-done',
  START_SLI_RETRIEVAL = 'sh.keptn.internal.event.get-sli',
  SLI_RETRIEVAL_DONE = 'sh.keptn.internal.event.get-sli.done',
  DONE = 'sh.keptn.events.done',
  PROBLEM_OPEN = 'sh.keptn.event.problem.open',
  PROBLEM_DETECTED = 'sh.keptn.events.problem',
  PROBLEM_RESOLVED = 'sh.keptn.events.problem.resolved',
  PROBLEM_CLOSED = 'sh.keptn.event.problem.close',
  APPROVAL_TRIGGERED = 'sh.keptn.event.approval.triggered',
  APPROVAL_FINISHED = 'sh.keptn.event.approval.finished'
};
