import {ResultTypes} from './result-types';

export class Sequence {
  name: string;
  project: string;
  service: string;
  shkeptncontext: string;
  stages: [
    {
      image: string,
      latestEvaluation: {
        result: ResultTypes,
        score: number
      },
      latestEvent: {
        id: string,
        time: string,
        type: string
      },
      latestFailedEvent: {
        id: string,
        time: string,
        type: string
      },
      name: string
    }
  ];
  state: string;
  time: string;
  problemTitle?: string;

  static fromJSON(data: any): Sequence {
    return Object.assign(new this(), data);
  }

  public getStage(stageName: string) {
    return this.stages.find(stage => stage.name === stageName);
  }

  public isFaulty(stageName: string) {
    return !!this.getStage(stageName).latestFailedEvent;
  }
}
