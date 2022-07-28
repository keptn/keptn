import { ChangeDetectorRef, Component, Input } from '@angular/core';
import { Trace } from '../../_models/trace';
import { EvaluationResult } from '../../../../shared/interfaces/evaluation-result';
import { DataService } from '../../_services/data.service';
import { DateUtil } from '../../_utils/date.utils';
import {
  EvaluationBadgeVariant,
  getEvaluationBadgeState,
  getEvaluationResultBadgeState,
} from '../ktb-evaluation-badge/ktb-evaluation-badge.utils';

export interface EventData {
  project: string;
  stage: string;
  service: string;
}

interface EvaluationInfo {
  trace: Trace | undefined;
  showHistory: boolean;
  data: EventData;
}

@Component({
  selector: 'ktb-evaluation-info',
  templateUrl: './ktb-evaluation-info.component.html',
  styleUrls: ['./ktb-evaluation-info.component.scss'],
})
export class KtbEvaluationInfoComponent {
  private eventData?: EventData;
  private _evaluationHistory?: Trace[];
  private _evaluation?: Trace;
  public readonly evaluationHistoryCount = 5;
  public showHistory = false;
  public evaluationsLoaded = false;
  public EvaluationBadgeFillState = EvaluationBadgeVariant;
  public getEvaluationState = getEvaluationBadgeState;
  public getEvaluationResultState = getEvaluationResultBadgeState;

  @Input() public overlayDisabled = false;
  @Input() public fillState = EvaluationBadgeVariant.FILL;
  @Input() public set evaluation(evaluation: Trace | undefined) {
    this._evaluation = evaluation;
  }
  get evaluation(): Trace | undefined {
    return this._evaluation;
  }

  @Input() evaluationResult?: EvaluationResult;

  @Input()
  public set evaluationInfo(evaluation: EvaluationInfo | undefined) {
    const idBefore = this.evaluation?.id;
    this.evaluation = evaluation?.trace;
    this.evaluationsLoaded = !!evaluation?.trace?.data.evaluationHistory?.length;
    this.showHistory = evaluation?.showHistory ?? false;
    this.eventData = evaluation?.data;

    if (idBefore !== evaluation?.trace?.id || (!evaluation?.trace && evaluation?.data)) {
      this.fetchEvaluationHistory();
    }
  }

  get evaluationHistory(): Trace[] {
    return (
      this._evaluationHistory ||
      this.evaluation?.data?.evaluationHistory
        ?.filter((evaluation) => evaluation.id !== this.evaluation?.id)
        .slice(0, this.evaluationHistoryCount) ||
      []
    );
  }

  constructor(private dataService: DataService, private changeDetectorRef_: ChangeDetectorRef) {}

  private fetchEvaluationHistory(): void {
    const evaluation = this.evaluation;
    let _eventData = this.eventData;
    if (this.evaluation && this.evaluation.data.project && this.evaluation.data.stage && this.evaluation.data.service) {
      _eventData = {
        project: this.evaluation.data.project,
        service: this.evaluation.data.service,
        stage: this.evaluation.data.stage,
      };
    }

    if (this.showHistory && _eventData) {
      // currently the event endpoint does not support skipping entries
      // the other endpoint we have does not support excluding invalidated evaluations
      // we can't use fromTime here if we have a limit. 10 new evaluations and limit to 5 would not pull the new ones
      this.dataService
        .getEvaluationResults(_eventData, this.evaluationHistoryCount + (this.evaluation ? 1 : 0), false)
        .subscribe((traces: Trace[]) => {
          traces.sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
          this.evaluationsLoaded = true;
          // we don't have an evaluation trace if the sequence is currently running or it just doesn't have an evaluation task
          if (evaluation) {
            this._evaluationHistory = undefined;
            evaluation.data.evaluationHistory = traces;
          } else {
            this._evaluationHistory = traces;
          }
          this.changeDetectorRef_.markForCheck();
        });
    }
  }
}
