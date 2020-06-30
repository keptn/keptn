enum EventTypes {
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
const EVENT_LABELS = {
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
const EVENT_ICONS = {
  [EventTypes.CONFIGURATION_CHANGE]: "duplicate",
  [EventTypes.DEPLOYMENT_FINISHED]: "deploy",
  [EventTypes.TESTS_FINISHED]: "perfromance-health",
  [EventTypes.START_EVALUATION]: "traffic-light",
  [EventTypes.EVALUATION_DONE]: "traffic-light",
  [EventTypes.START_SLI_RETRIEVAL]: "collector",
  [EventTypes.SLI_RETRIEVAL_DONE]: "collector",
  [EventTypes.PROBLEM_OPEN]: "criticalevent",
  [EventTypes.PROBLEM_DETECTED]: "criticalevent",
  [EventTypes.PROBLEM_CLOSED]: "applicationhealth",
  [EventTypes.APPROVAL_TRIGGERED]: "unknown",
  [EventTypes.APPROVAL_FINISHED]: "checkmark"
};

const DEFAULT_ICON = "information";

class Trace {
  id: string;
  shkeptncontext: string;
  source: string;
  time: Date;
  type: string;
  label: string;
  icon: string;
  image: string;
  plainEvent: string;
  data: {
    project: string;
    service: string;
    stage: string;

    image: string;
    tag: string;

    deploymentURILocal: string;
    deploymentURIPublic: string;

    deploymentstrategy: string;
    labels: Map<any, any>;
    result: string;
    teststrategy: string;

    start: Date;
    end: Date;

    canary: {
      action: string;
      value: number;
    };
    eventContext: {
      shkeptncontext: string;
      token: string;
    };
    valuesCanary: {
      image: string;
    };

    evaluationdetails: {
      indicatorResults: any;
      result: string;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;
    };

    evaluationHistory: Trace[];

    ProblemTitle: string;
    ImpactedEntity: string;
    ProblemDetails: {
      tagsOfAffectedEntities: {
        key: string;
        value: string;
      }
    };

    approval: {
      result: string;
      status: string;
    };

    Tags: string;
    State: string;
  };

  isFaulty(): string {
    let result: string = null;
    if(this.data) {
      if(this.isFailed() || this.isProblem()) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isWarning(): string {
    let result: string = null;
    if(this.data) {
      if(this.data.result == 'warning') {
        result = this.data.stage;
      }
    }
    return result;
  }

  isFailed(): boolean {
    return this.data.result == 'fail' || this.type === EventTypes.APPROVAL_FINISHED && this.data.approval.result == 'failed';
  }

  isProblem(): boolean {
    return this.type.indexOf('problem') != -1;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == 'pass' || this.type === EventTypes.APPROVAL_FINISHED && this.data.approval.result == 'pass') {
        result = true;
      }
    }
    return !this.isFaulty() && result;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      if(this.type === EventTypes.PROBLEM_DETECTED && this.data.State === "RESOLVED") {
        this.label = EVENT_LABELS[EventTypes.PROBLEM_RESOLVED];
      } else {
        this.label = EVENT_LABELS[this.type] || this.type;
      }
    }

    return this.label;
  }

  getIcon() {
    if(!this.icon) {
      this.icon = EVENT_ICONS[this.type] || DEFAULT_ICON;
    }
    return this.icon;
  }

  getShortImageName() {
    if(!this.image) {
      if(this.data.image && this.data.tag)
        this.image = [this.data.image.split("/").pop(), this.data.tag].join(":");
      else if(this.data.image)
        this.image = this.data.image.split("/").pop();
      else if(this.data.valuesCanary)
        this.image = this.data.valuesCanary.image.split("/").pop();
    }

    return this.image;
  }

  getProject(): string {
    return this.data.project;
  }

  getService(): string {
    return this.data.service;
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data, { plainEvent: JSON.parse(JSON.stringify(data)) });
  }
}

export {Trace, EVENT_LABELS, EventTypes}
