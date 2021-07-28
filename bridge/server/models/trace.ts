import { ResultTypes } from './result-types.js';
import { EventTypes } from './event-types.js';
import { IndicatorResult } from '../interfaces/indicator-result.js';

export class Trace {
  id!: string;
  shkeptncontext!: string;
  triggeredid?: string;
  type!: EventTypes | string;
  time?: Date;
  data!: {
    project?: string;
    service?: string;
    stage?: string;

    image?: string;
    tag?: string;

    deployment?: {
      deploymentNames: string[];
      deploymentURIsLocal: string[];
      deploymentURIsPublic: string[];
      deploymentstrategy: string;
      gitCommit: string;
    };

    deploymentURILocal?: string;
    deploymentURIPublic?: string;

    message?: string;

    labels?: { [key: string]: string };
    result?: ResultTypes;
    teststrategy?: string;

    start?: Date;
    end?: Date;

    canary?: {
      action: string;
      value: number;
    };
    eventContext?: {
      shkeptncontext: string;
      token: string;
    };
    configurationChange?: {
      values: {
        image: unknown
      }
    };

    evaluation?: {
      comparedEvents?: string[];
      indicatorResults: IndicatorResult[];
      result: ResultTypes;
      score: number;
      sloFileContent: string;
      timeEnd: Date;
      timeStart: Date;
      score_pass: string;
      score_warning: string;
      compare_with: string;
      include_result_with_score: string;
      number_of_comparison_results: number;
      number_of_missing_comparison_results: number;
      sloFileContentParsed: string;
    };

    evaluationHistory?: Trace[];

    problem?: {
      ProblemTitle: string;
      ImpactedEntity: string;
      ProblemDetails: {
        tagsOfAffectedEntities: {
          key: string;
          value: string;
        }[]
      };
    };

    approval?: {
      result: string;
      status: string;
    };

    action?: {
      action: string;
      description: string;
      name: string;
    }

    Tags?: string;
    State?: string;
  };

  static fromJSON(data: unknown): Trace {
    if (data instanceof Trace) {
      return data;
    }

    const plainEvent = JSON.parse(JSON.stringify(data));
    const trace: Trace = Object.assign(new this(), data, { plainEvent });

    if (trace.data?.evaluationHistory?.length) {
      trace.data.evaluationHistory = trace.data.evaluationHistory.map(t => Trace.fromJSON(t));
    }

    return trace;
  }

  getShortImageName(): string | undefined {
    let image;
    if (this.data.image && this.data.tag) {
      image = [this.data.image.split('/').pop(), this.data.tag].join(':');
    }
    else if (this.data.image) {
      image = this.data.image.split('/').pop();
    }
    else if (this.data.configurationChange?.values) {
      image = this.getConfigurationChangeImage();
    }
    return image;
  }

  private getConfigurationChangeImage(): string | undefined {
    return typeof this.data.configurationChange?.values.image === 'string'
      ? this.data.configurationChange.values.image.split('/').pop()
      : undefined;
  }

  public getDeploymentUrl() {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }
}
