import moment from 'moment';
import { EventTypes } from '../../../shared/interfaces/event-types';
import { ResultTypes } from '../../../shared/models/result-types';
import { ApprovalStates } from './approval-states';
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
  heatmapLabel?: string;
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

  public isApprovalTriggered(): boolean {
    return this.type === EventTypes.APPROVAL_TRIGGERED;
  }

  public isDirectDeployment(): boolean {
    return this.type === EventTypes.DEPLOYMENT_FINISHED && this.data?.deployment?.deploymentstrategy === 'direct';
  }

  private isApproved(): boolean {
    return this.data.approval?.result === ApprovalStates.APPROVED;
  }

  public isDeployment(): string | undefined {
    return this.type === EventTypes.DEPLOYMENT_TRIGGERED ? this.data.stage : undefined;
  }

  public isEvaluationInvalidation(): boolean {
    return this.type === EventTypes.EVALUATION_INVALIDATED;
  }

  hasLabels(): boolean {
    return Object.keys(this.data.labels || {}).length > 0;
  }

  getProblemTitle(): string | undefined {
    return this.data.problem?.ProblemTitle;
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

  getHeatmapLabel(): string {
    if (!this.heatmapLabel) {
      this.heatmapLabel = this.getChartLabel();
    }
    return this.heatmapLabel;
  }

  setHeatmapLabel(label: string): void {
    this.heatmapLabel = label;
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

  getRemediationAction(): Trace | undefined {
    return this.findTrace((t) => t.isRemediationAction());
  }

  getDeploymentUrl(): string | undefined {
    return this.data.deployment?.deploymentURIsPublic?.find(() => true);
  }

  findLastTrace(comp: (args: Trace) => boolean): Trace | undefined {
    if (comp(this)) {
      return this;
    } else {
      return this.traces.reduce((result: Trace | undefined, trace) => trace.findTrace(comp) || result, undefined);
    }
  }

  getProblemDetails(): string | undefined {
    return this.data.problem?.ImpactedEntity || this.data.problem?.ProblemTitle;
  }
}

export { Trace };
