import { EventTypes } from '../interfaces/event-types';
import { ResultTypes } from './result-types';
import { IndicatorResult } from '../interfaces/indicator-result';

export interface TraceData {
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
      image: unknown;
    };
  };

  evaluation?: {
    comparedEvents?: string[];
    indicatorResults: IndicatorResult[];
    result: ResultTypes;
    score: number;
    sloFileContent: string;
    timeEnd: Date;
    timeStart: Date;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    score_pass: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    score_warning: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    compare_with: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    include_result_with_score: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    number_of_comparison_results: number;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    number_of_missing_comparison_results: number;
    sloFileContentParsed: string;
  };

  evaluationHistory?: Trace[];

  problem?: {
    // eslint-disable-next-line @typescript-eslint/naming-convention
    ProblemTitle: string;
    // eslint-disable-next-line @typescript-eslint/naming-convention
    ImpactedEntity: string;

    ProblemDetails: {
      tagsOfAffectedEntities: {
        key: string;
        value: string;
      }[];
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
  };
  // eslint-disable-next-line @typescript-eslint/naming-convention
  Tags?: string;
  // eslint-disable-next-line @typescript-eslint/naming-convention
  State?: string;
}

export class Trace {
  id!: string;
  shkeptncontext!: string;
  triggeredid?: string;
  type!: EventTypes | string;
  time?: Date;
  data!: TraceData;

  public getShortImageName(): string | undefined {
    let image;
    if (this.data.image && this.data.tag) {
      image = [this.data.image.split('/').pop(), this.data.tag].join(':');
    } else if (this.data.image) {
      image = this.data.image.split('/').pop();
    } else if (this.data.configurationChange?.values) {
      image = this.getConfigurationChangeImage();
    }
    return image;
  }

  public getConfigurationChangeImage(): string | undefined {
    return typeof this.data.configurationChange?.values.image === 'string'
      ? this.data.configurationChange.values.image.split('/').pop()
      : undefined;
  }

  public getDeploymentUrl(): string | undefined {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }
}
