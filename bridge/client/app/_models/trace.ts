let labels = {
  "sh.keptn.internal.event.service.create": "Service create",
  "sh.keptn.event.configuration.change": "Configuration change",
  "sh.keptn.event.monitoring.configure": "Configure monitoring",
  "sh.keptn.events.deployment-finished": "Deployment finished",
  "sh.keptn.events.tests-finished": "Tests finished",
  "sh.keptn.event.start-evaluation": "Start evaluation",
  "sh.keptn.events.evaluation-done": "Evaluation done",
  "sh.keptn.internal.event.get-sli": "Start SLI retrieval",
  "sh.keptn.internal.event.get-sli.done": "SLI retrieval done",
  "sh.keptn.events.done": "Done",
  "sh.keptn.event.problem.open": "Problem open",
  "sh.keptn.events.problem": "Problem detected",
  "sh.keptn.events.problem.resolved": "Problem resolved",
  "sh.keptn.event.problem.close": "Problem closed"
};
let icons = {
  "sh.keptn.event.configuration.change": "duplicate",
  "sh.keptn.events.deployment-finished": "deploy",
  "sh.keptn.events.tests-finished": "perfromance-health",
  "sh.keptn.event.start-evaluation": "traffic-light",
  "sh.keptn.events.evaluation-done": "traffic-light",
  "sh.keptn.internal.event.get-sli": "collector",
  "sh.keptn.internal.event.get-sli.done": "collector",
  "sh.keptn.event.problem.open": "criticalevent",
  "sh.keptn.events.problem": "criticalevent",
  "sh.keptn.event.problem.close": "applicationhealth"
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
    return this.data.result == 'fail';
  }

  isProblem(): boolean {
    return this.type.indexOf('problem') != -1;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == 'pass') {
        result = true;
      }
    }
    return !this.isFaulty() && result;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      if(this.type === "sh.keptn.events.problem" && this.data.State === "RESOLVED") {
        this.label = labels["sh.keptn.events.problem.resolved"];
      } else {
        this.label = labels[this.type] || this.type;
      }
    }

    return this.label;
  }

  getIcon() {
    if(!this.icon) {
      this.icon = icons[this.type] || "information";
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

export {Trace, labels}
