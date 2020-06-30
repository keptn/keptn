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
  [EVENT_TYPES.SERVICE_CREATE]: "Service create",
  [EVENT_TYPES.CONFIGURATION_CHANGE]: "Configuration change",
  [EVENT_TYPES.CONFIGURE_MONITORING]: "Configure monitoring",
  [EVENT_TYPES.DEPLOYMENT_FINISHED]: "Deployment finished",
  [EVENT_TYPES.TESTS_FINISHED]: "Tests finished",
  [EVENT_TYPES.START_EVALUATION]: "Start evaluation",
  [EVENT_TYPES.EVALUATION_DONE]: "Evaluation done",
  [EVENT_TYPES.START_SLI_RETRIEVAL]: "Start SLI retrieval",
  [EVENT_TYPES.SLI_RETRIEVAL_DONE]: "SLI retrieval done",
  [EVENT_TYPES.DONE]: "Done",
  [EVENT_TYPES.PROBLEM_OPEN]: "Problem open",
  [EVENT_TYPES.PROBLEM_DETECTED]: "Problem detected",
  [EVENT_TYPES.PROBLEM_RESOLVED]: "Problem resolved",
  [EVENT_TYPES.PROBLEM_CLOSED]: "Problem closed",
  [EVENT_TYPES.APPROVAL_TRIGGERED]: "Approval triggered",
  [EVENT_TYPES.APPROVAL_FINISHED]: "Approval finished"
};
const EVENT_ICONS = {
  [EVENT_TYPES.CONFIGURATION_CHANGE]: "duplicate",
  [EVENT_TYPES.DEPLOYMENT_FINISHED]: "deploy",
  [EVENT_TYPES.TESTS_FINISHED]: "perfromance-health",
  [EVENT_TYPES.START_EVALUATION]: "traffic-light",
  [EVENT_TYPES.EVALUATION_DONE]: "traffic-light",
  [EVENT_TYPES.START_SLI_RETRIEVAL]: "collector",
  [EVENT_TYPES.SLI_RETRIEVAL_DONE]: "collector",
  [EVENT_TYPES.PROBLEM_OPEN]: "criticalevent",
  [EVENT_TYPES.PROBLEM_DETECTED]: "criticalevent",
  [EVENT_TYPES.PROBLEM_CLOSED]: "applicationhealth",
  [EVENT_TYPES.APPROVAL_TRIGGERED]: "unknown",
  [EVENT_TYPES.APPROVAL_FINISHED]: "checkmark"
};

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
    return this.data.result == 'fail' || this.type === EVENT_TYPES.APPROVAL_FINISHED && this.data.approval.result == 'failed';
  }

  isProblem(): boolean {
    return this.type.indexOf('problem') != -1;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == 'pass' || this.type === EVENT_TYPES.APPROVAL_FINISHED && this.data.approval.result == 'pass') {
        result = true;
      }
    }
    return !this.isFaulty() && result;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      if(this.type === EVENT_TYPES.PROBLEM_DETECTED && this.data.State === "RESOLVED") {
        this.label = EVENT_LABELS[EVENT_TYPES.PROBLEM_RESOLVED];
      } else {
        this.label = EVENT_LABELS[this.type] || this.type;
      }
    }

    return this.label;
  }

  getIcon() {
    if(!this.icon) {
      this.icon = EVENT_ICONS[this.type] || "information";
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

export {Trace, EVENT_LABELS, EVENT_TYPES}
