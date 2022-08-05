import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { filter, map, takeUntil, tap } from 'rxjs/operators';
import { IClientFeatureFlags } from '../../../../../shared/interfaces/feature-flags';
import { FeatureFlagsService } from '../../../_services/feature-flags.service';
import { createDataPoints, parseSloOfEvaluations } from '../ktb-evaluation-details-utils';
import { DataService } from '../../../_services/data.service';
import { Trace } from '../../../_models/trace';
import { IDataPoint } from '../../../_interfaces/heatmap';
import { EvaluationHistory } from '../../../_interfaces/evaluation-history';
import { ChartItem } from '../../../_interfaces/chart';
import { DateUtil } from '../../../_utils/date.utils';
import {
  createChartPoints,
  createChartTooltipLabels,
  createChartXLabels,
} from '../ktb-evaluation-details-line-chart-utils';
import { DateFormatPipe } from 'ngx-moment';
import { Subject } from 'rxjs';
import { IndicatorResult } from '../../../../../shared/interfaces/indicator-result';

export interface IEvaluationSelectionData {
  shouldSelect: boolean;
  evaluation?: Trace;
}
export type TChartType = 'heatmap' | 'chart';
type FuncEventIdToEvaluation = (eventId: string) => Trace | undefined;
type FuncEventIdExists = (eventId: string) => boolean;

@Component({
  selector: 'ktb-evaluation-chart[evaluationData]',
  templateUrl: './ktb-evaluation-chart.component.html',
})
export class KtbEvaluationChartComponent implements OnInit, OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private _evaluationData: IEvaluationSelectionData = { shouldSelect: false };
  public d3Enabled$ = this.featureFlagService.featureFlags$.pipe(
    map((featureFlags: IClientFeatureFlags) => featureFlags.D3_ENABLED)
  );
  public chartType: TChartType | null = 'heatmap';
  public selectedIdentifier = '';
  public dataPoints?: IDataPoint[];
  public evaluationHistoryUpdates?: EvaluationHistory;
  public chartPoints?: ChartItem[];
  public chartXLabels: Record<number, string> = {};
  public chartTooltipLabels: Record<number, string> = {};
  public numberOfMissingEvaluationComparisons = 0;

  @Output() selectedEvaluationChange = new EventEmitter<Trace | undefined>();
  @Output() comparedIndicatorResultsChange = new EventEmitter<IndicatorResult[][]>();

  @Input()
  set evaluationData(evaluationData: IEvaluationSelectionData) {
    this.setEvaluation(evaluationData);
  }
  get evaluationData(): IEvaluationSelectionData {
    return this._evaluationData;
  }

  constructor(
    private dataService: DataService,
    private featureFlagService: FeatureFlagsService,
    private dateFormatPipe: DateFormatPipe,
    private dateUtil: DateUtil
  ) {}

  public ngOnInit(): void {
    this.dataService.evaluationResults
      .pipe(
        filter((results): results is Required<EvaluationHistory> => !!results.traces?.length),
        tap((results) => {
          parseSloOfEvaluations(results.traces);
        }),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((results) => {
        // check if there already are evaluations and the incoming data is an update. If it is an update postpone it and show a refresh button
        if (this.evaluationData.evaluation?.data.evaluationHistory?.length) {
          this.evaluationHistoryUpdates = results;
        } else {
          this.refreshEvaluationBoard(results);
        }
      });
  }

  private setEvaluation(evaluationData: IEvaluationSelectionData): void {
    if (this._evaluationData.evaluation?.id !== evaluationData.evaluation?.id) {
      this._evaluationData = evaluationData;
      this.evaluationHistoryUpdates = undefined;
      if (evaluationData.shouldSelect) {
        this.selectedIdentifier = evaluationData.evaluation?.id ?? '';
      }
    }

    if (this._evaluationData.evaluation) {
      // update or initially load up to 50 evaluations
      this.dataService.loadEvaluationResults(this._evaluationData.evaluation);
    }
  }

  public dataPointChanged(identifier: string): void {
    const mapComparedEventsToEvaluations =
      (history: Trace[]): FuncEventIdToEvaluation =>
      (eventId: string) =>
        history.find((e) => e.id === eventId);
    const filterOutUndefinedEvaluations = (evaluation?: Trace): evaluation is Trace => !!evaluation;
    const mapEvaluationToIndicatorResult = (evaluation: Trace): IndicatorResult[] =>
      evaluation.data.evaluation?.indicatorResults ?? [];

    const evaluationHistory = this.evaluationData.evaluation?.data.evaluationHistory ?? [];
    const changedEvaluation = evaluationHistory.find((value) => value.id === identifier);
    const secondaryHighlightIndices = changedEvaluation?.data.evaluation?.comparedEvents
      ?.map(mapComparedEventsToEvaluations(evaluationHistory))
      .filter(filterOutUndefinedEvaluations)
      .map(mapEvaluationToIndicatorResult);

    this.selectedIdentifier = identifier;
    this.selectedEvaluationChange.emit(changedEvaluation);
    this.comparedIndicatorResultsChange.emit(secondaryHighlightIndices);
    this.setMissingComparedEvaluations(changedEvaluation, evaluationHistory);
    // this.selectEvaluationData(changedEvaluation, true);
  }

  private setMissingComparedEvaluations(evaluation: Trace | undefined, evaluationHistory: Trace[]): void {
    if (!evaluation?.data.evaluation?.comparedEvents) {
      this.numberOfMissingEvaluationComparisons = 0;
      return;
    }
    const filterOutNotFoundEvaluations =
      (history: Trace[]): FuncEventIdExists =>
      (eventId: string) =>
        history.some((historyEvaluation) => historyEvaluation.id === eventId);

    const comparedEventsCount = evaluation.data.evaluation.comparedEvents.length;
    const foundIdentifiers = evaluation.data.evaluation.comparedEvents.filter(
      filterOutNotFoundEvaluations(evaluationHistory)
    );

    // TODO: leave it here or emit it and show it where it was shown before (inside the dt-consumption thing)?
    this.numberOfMissingEvaluationComparisons =
      comparedEventsCount === 0 ? 0 : comparedEventsCount - foundIdentifiers.length;
  }

  public refreshEvaluationBoard(results: EvaluationHistory): void {
    if (this.evaluationData.evaluation) {
      if (results.type === 'evaluationHistory' && results.triggerEvent === this.evaluationData.evaluation) {
        this.evaluationData.evaluation.data.evaluationHistory = [
          ...(results.traces || []),
          ...(this.evaluationData.evaluation.data.evaluationHistory || []),
        ].sort((a, b) => DateUtil.compareTraceTimesDesc(a, b));
      } else if (
        results.type === 'invalidateEvaluation' &&
        this.evaluationData.evaluation.data.project === results.triggerEvent.data.project &&
        this.evaluationData.evaluation.data.service === results.triggerEvent.data.service &&
        this.evaluationData.evaluation.data.stage === results.triggerEvent.data.stage
      ) {
        this.evaluationData.evaluation.data.evaluationHistory =
          this.evaluationData.evaluation.data.evaluationHistory?.filter((e) => e.id !== results.triggerEvent.id);
      }
      const isSelectedIdentifierInHistory =
        this.selectedIdentifier &&
        this.evaluationData.evaluation.data.evaluationHistory?.some((h) => h.id === this.selectedIdentifier);

      const selectedIdentifier = isSelectedIdentifierInHistory ? this.selectedIdentifier : '';
      this.updateChartData(this.evaluationData.evaluation.data.evaluationHistory ?? [], selectedIdentifier);
    }
    this.evaluationHistoryUpdates = undefined;
  }

  private updateChartData(evaluationHistory: Trace[], selectedIdentifier: string): void {
    this.dataPointChanged(selectedIdentifier);
    this.dataPoints = createDataPoints(evaluationHistory);
    this.chartPoints = createChartPoints(evaluationHistory);
    this.chartXLabels = createChartXLabels(evaluationHistory);
    this.chartTooltipLabels = createChartTooltipLabels(evaluationHistory, (time: string) =>
      this.dateFormatPipe.transform(time, this.dateUtil.getDateTimeFormat())
    );
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
