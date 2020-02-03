import {Stage} from "./stage";

export class Trace {
  id: string;
  shkeptncontext: string;
  source: string;
  time: Date;
  type: string;
  label: string;
  plainEvent: string;
  data: {
    project: string;
    service: string;
    stage: string;

    deploymentstrategy: string;
    labels: string;
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

    ImpactedEntity: string;
    Tags: string;
  };

  isFaulty(): string {
    let result: string = null;
    if(this.data) {
      if(this.data.result == "fail" || this.type.indexOf('problem.open') != -1) {
        result = this.data.stage;
      }
    }
    return result;
  }

  isSuccessful(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == "pass") {
        result = true;
      }
    }
    return !this.isFaulty() && result;
  }

  getLabel(): string {
    // TODO: use translation file
    if(!this.label) {
      switch(this.type) {
        case "sh.keptn.internal.event.service.create": {
          this.label = "Service create";
          break;
        }
        case "sh.keptn.event.configuration.change": {
          this.label = "Configuration change";
          break;
        }
        case "sh.keptn.event.monitoring.configure": {
          this.label = "Configure monitoring";
          break;
        }
        case "sh.keptn.events.deployment-finished": {
          this.label = "Deployment finished";
          break;
        }
        case "sh.keptn.events.tests-finished": {
          this.label = "Tests finished";
          break;
        }
        case "sh.keptn.events.evaluation-done": {
          this.label = "Evaluation done";
          break;
        }
        case "sh.keptn.internal.event.get-sli": {
          this.label = "Start SLI retrieval";
          break;
        }
        case "sh.keptn.internal.event.get-sli.done": {
          this.label = "SLI retrieval done";
          break;
        }
        case "sh.keptn.events.done": {
          this.label = "Done";
          break;
        }
        case "sh.keptn.event.problem.open": {
          this.label = "Problem open";
          break;
        }
        case "sh.keptn.events.problem": {
          this.label = "Problem detected";
          break;
        }
        case "sh.keptn.event.problem.close": {
          this.label = "Problem closed";
          break;
        }
        default: {
          this.label = this.type;
          break;
        }
      }
    }

    return this.label;
  }

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
