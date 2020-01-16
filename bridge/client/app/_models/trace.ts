export class Trace {
  id: string;
  shkeptncontext: string;
  source: string;
  time: Date;
  type: string;
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
      indicatorResults: string;
      result: string;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;
    };
  };

  isFaulty(): boolean {
    let result: boolean = false;
    if(this.data) {
      if(this.data.result == "fail") {
        result = true;
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

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }
}
