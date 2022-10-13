import moment from 'moment';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ResultTypes } from '../../../shared/models/result-types';
import { ApprovalStates } from '../../../shared/models/approval-states';
import { EVENT_ICONS } from './event-icons';
import { Trace as ts, TraceData } from '../../../shared/models/trace';
import { DtIconType } from '@dynatrace/barista-icons';

class Trace extends ts {
  traces: Trace[] = [];
  triggeredid?: string;
  type!: EventTypes | string;
  data!: Omit<TraceData, 'evaluationHistory'> & {
    evaluationHistory?: Trace[];
  };
  started?: boolean;
  finished?: boolean;
  source?: string;
  label?: string;
  icon?: DtIconType;
  image?: string;
  plainEvent?: string;
  time?: string;
  labelMap?: Map<string, string>;

  static fromJSON(data: unknown): Trace {
    if (data instanceof Trace) {
      return data;
    }

    const plainEvent = JSON.parse(JSON.stringify(data));
    const trace: Trace = Object.assign(new this(), data, { plainEvent });

    if (trace.data?.evaluationHistory?.length) {
      trace.data.evaluationHistory = trace.data.evaluationHistory.map((t) => Trace.fromJSON(t));
    }

    return trace;
  }

  static traceMapper(traces: Trace[]): Trace[] {
    traces = traces.map((trace) => Trace.fromJSON(trace));
    return ts.traceMapperGlobal(traces);
  }

  static get defaultTrace(): Partial<Trace> {
    return {
      data: {
        project: undefined,
        service: undefined,
        stage: undefined,
      },
      id: undefined,
      type: undefined,
      time: undefined,
      shkeptncontext: undefined,
    };
  }

  get project(): string | undefined {
    return this.data.project;
  }

  get service(): string | undefined {
    return this.data.service;
  }

  get stage(): string | undefined {
    return this.data.stage;
  }

  get labels(): Map<string, string> | undefined {
    if (!this.labelMap) {
      let map: Map<string, string> | undefined;
      if (this.data.labels) {
        map = new Map<string, string>();
        for (const key of Object.keys(this.data.labels)) {
          map.set(key, this.data.labels[key]);
        }
      }
      this.labelMap = map;
    }
    return this.labelMap;
  }

  isSuccessful(stageName?: string): boolean {
    let result = false;
    if (
      (this.isFinished() && this.getFinishedEvent()?.data.result === ResultTypes.PASSED) ||
      (this.isApprovalFinished() && this.isApproved()) ||
      (this.isProblem() && this.isProblemResolvedOrClosed()) ||
      this.isSuccessfulRemediation()
    ) {
      result = stageName ? this.data.stage === stageName : true;
    }
    return !this.isFaulty() && result;
  }

  public isRemediationAction(): boolean {
    return this.type === EventTypes.ACTION_TRIGGERED;
  }

  public getRemediationActionDescription(): string | undefined {
    return this.data.action?.description;
  }

  public getRemediationActionName(): string | undefined {
    return this.data.action?.name;
  }

  private isApproved(): boolean {
    return this.data.approval?.result === ApprovalStates.APPROVED;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels || {}).length > 0;
  }

  getIcon(): DtIconType {
    if (!this.icon) {
      this.icon = EVENT_ICONS[this.getShortType()] || EVENT_ICONS.default;
    }
    return this.icon;
  }

  getChartLabel(): string {
    return this.data.labels?.buildId ?? moment(this.time).format('YYYY-MM-DD HH:mm');
  }

  isStarted(): boolean {
    if (!this.started && this.traces) {
      this.started = this.traces.some((t) => t.isStartedEvent() || t.isStarted());
    }

    return !!this.started;
  }

  isTriggered(): boolean {
    return this.type.endsWith('.triggered');
  }

  isLoading(): boolean {
    return this.isStarted() && !this.isFinished();
  }

  isInvalidated(): boolean {
    return !!this.traces.find((e) => e.isEvaluationInvalidation() && e.triggeredid === this.id);
  }

  getDeploymentUrl(): string | undefined {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }

  getProblemDetails(): string | undefined {
    return this.data.problem?.ImpactedEntity || this.data.problem?.ProblemTitle;
  }
}

export { Trace };
