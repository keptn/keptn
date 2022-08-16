import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationChartComponent } from './ktb-evaluation-chart.component';
import { EvaluationsMock } from '../../../_services/_mockData/evaluations.mock';
import { AppUtils } from '../../../_utils/app.utils';
import { Trace } from '../../../_models/trace';
import { DataService } from '../../../_services/data.service';
import { of } from 'rxjs';
import { ResultTypes } from '../../../../../shared/models/result-types';
import { IEvaluationData } from '../../../../../shared/models/trace';
import { EvaluationHistory } from '../../../_interfaces/evaluation-history';
import { EventTypes } from '../../../../../shared/interfaces/event-types';
import { Component, Input } from '@angular/core';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { DtChartModule } from '@dynatrace/barista-components/chart';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { FlexModule } from '@angular/flex-layout';
import { KtbHeatmapModule } from '../../ktb-heatmap/ktb-heatmap.module';
import { KtbPipeModule } from '../../../_pipes/ktb-pipe.module';
import { DateFormatPipe, MomentModule } from 'ngx-moment';
import { KtbChartModule } from '../../ktb-chart/ktb-chart.module';
import { CommonModule } from '@angular/common';
import { TChartType } from '../ktb-evaluation-details-utils';
import { IndicatorResult } from '../../../../../shared/interfaces/indicator-result';
import { HttpClientTestingModule } from '@angular/common/http/testing';

@Component({
  selector: 'ktb-evaluation-chart-legacy',
  template: '',
})
class FakeKtbEvaluationChartLegacyComponent {
  @Input() evaluationData: unknown;
  @Input() chartType: TChartType = 'heatmap';
  @Input() evaluationHistory: unknown;
}

describe(KtbEvaluationChartComponent.name, () => {
  let component: KtbEvaluationChartComponent;
  let fixture: ComponentFixture<KtbEvaluationChartComponent>;
  let mockEvaluation: Trace;
  let dataService: DataService;

  beforeEach(async () => {
    // disable legacy chart (prevent "animate" error during tests)
    await TestBed.configureTestingModule({
      declarations: [FakeKtbEvaluationChartLegacyComponent, KtbEvaluationChartComponent],
      imports: [
        CommonModule,
        DtButtonGroupModule,
        DtButtonModule,
        DtChartModule,
        DtIconModule.forRoot({
          svgIconLocation: `assets/icons/{{name}}.svg`,
        }),
        DtKeyValueListModule,
        FlexModule,
        KtbChartModule,
        KtbHeatmapModule,
        KtbPipeModule,
        MomentModule,
        HttpClientTestingModule,
      ],
      providers: [DateFormatPipe],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEvaluationChartComponent);
    component = fixture.componentInstance;

    const data = AppUtils.copyObject(EvaluationsMock);
    data.data.evaluationHistory = undefined;
    mockEvaluation = Trace.fromJSON(data);
    dataService = TestBed.inject(DataService);
  });

  describe('input, ngOnInit', () => {
    it('should create', () => {
      fixture.detectChanges();
      expect(component).toBeTruthy();
    });

    it('should load evaluation history', () => {
      // given
      const loadHistorySpy = jest.spyOn(dataService, 'loadEvaluationResults');

      // when
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      // then
      expect(loadHistorySpy).toHaveBeenCalled();
    });

    it('should postpone evaluation history update', () => {
      // given
      const loadEvaluationResultsSpy = jest.spyOn(dataService, 'loadEvaluationResults');
      const refreshEvaluations = jest.spyOn(component, 'refreshEvaluationBoard');
      jest
        .spyOn(dataService, 'evaluationResults', 'get')
        .mockReturnValue(of(getEvaluationResults(), getEvaluationResults()));
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      expect(component.evaluationHistoryUpdates).toBeUndefined();

      // when, then
      fixture.detectChanges();
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      // then
      expect(component.evaluationHistoryUpdates).toEqual(getEvaluationResults());
      expect(refreshEvaluations).toHaveBeenCalledTimes(1);
      expect(loadEvaluationResultsSpy).toHaveBeenCalledTimes(2);
    });

    it('should set selectedIdentifier if shouldSelect is true', () => {
      // given, when
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      // then
      expect(component.selectedIdentifier).toBe('01b1eff1-5bd9-4955-b2ef-30fac990b761');
    });

    it('should not set selectedIdentifier if shouldSelect is false', () => {
      // given, when
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: false,
      };

      // then
      expect(component.selectedIdentifier).toBe('');
    });

    it('should revert postponed history update if evaluation changed', () => {
      // given
      jest
        .spyOn(dataService, 'evaluationResults', 'get')
        .mockReturnValue(of(getEvaluationResults(), getEvaluationResults()));
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      // when
      component.evaluationData = {
        evaluation: Trace.fromJSON({ ...mockEvaluation, id: 'myOtherId' }),
        shouldSelect: true,
      };

      // then
      expect(component.evaluationHistoryUpdates).toBeUndefined();
    });
  });

  describe('refreshEvaluationBoard', () => {
    it('should correctly map data', () => {
      // given
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };

      // when
      component.refreshEvaluationBoard(getEvaluationResults());

      // then
      expect(component.evaluationData.evaluation?.data.evaluationHistory).toEqual(getEvaluationResults().traces);
      expect(component.evaluationData.evaluation?.data.evaluationHistory?.length).toBe(1);
    });

    it('should remove invalidated evaluation', () => {
      // given
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };
      component.refreshEvaluationBoard(getEvaluationResults());

      // when
      component.refreshEvaluationBoard({
        triggerEvent: {
          data: {
            project: 'sockshop',
            stage: 'staging',
            service: 'carts',
          },
          id: 'myId0',
        } as Trace,
        type: 'invalidateEvaluation',
      });

      // then
      expect(component.evaluationData.evaluation?.data.evaluationHistory?.length).toEqual(0);
    });
  });

  describe('dataPointChanged', () => {
    it('should correctly selected evaluation', () => {
      // given
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };
      const result = getEvaluationResults(1, 5);
      component.refreshEvaluationBoard(result);

      // when
      const evaluationChangeSpy = jest.spyOn(component.selectedEvaluationChange, 'emit');
      component.dataPointChanged('myId0');

      // then
      expect(evaluationChangeSpy).toHaveBeenCalledWith(result.traces?.[0]);
    });

    it('should emit indicator results of compared events', () => {
      // given
      component.evaluationData = {
        evaluation: mockEvaluation,
        shouldSelect: true,
      };
      const result = getEvaluationResults(5, 2, 2);
      component.refreshEvaluationBoard(result);

      // when
      const indicatorResultChangeSpy = jest.spyOn(component.comparedIndicatorResultsChange, 'emit');
      component.dataPointChanged('myId4');

      // then
      expect(indicatorResultChangeSpy).toHaveBeenCalledWith([getIndicatorResults(2), []]);
    });

    describe('number of missing comparisons', () => {
      it('should emit number of missing comparisons', () => {
        // given
        component.evaluationData = {
          evaluation: mockEvaluation,
          shouldSelect: true,
        };
        component.refreshEvaluationBoard(getEvaluationResults(1, 5));

        // when
        const missingComparisonSpy = jest.spyOn(component.numberOfMissingEvaluationComparisonsChange, 'emit');
        component.dataPointChanged('myId0');

        // then
        expect(missingComparisonSpy).toHaveBeenCalledWith(4);
      });

      it('should emit 0 missing comparisons if there are not any compared events', () => {
        // given
        component.evaluationData = {
          evaluation: mockEvaluation,
          shouldSelect: true,
        };
        component.refreshEvaluationBoard(getEvaluationResults(1, 0));

        // when
        const missingComparisonSpy = jest.spyOn(component.numberOfMissingEvaluationComparisonsChange, 'emit');
        component.dataPointChanged('myId0');

        // then
        expect(missingComparisonSpy).toHaveBeenCalledWith(0);
      });

      it('should emit 0 missing comparisons if all compared events are loaded', () => {
        // given
        component.evaluationData = {
          evaluation: mockEvaluation,
          shouldSelect: true,
        };
        component.refreshEvaluationBoard(getEvaluationResults(4, 3));

        // when
        const missingComparisonSpy = jest.spyOn(component.numberOfMissingEvaluationComparisonsChange, 'emit');
        component.dataPointChanged('myId3');

        // then
        expect(missingComparisonSpy).toHaveBeenCalledWith(0);
      });
    });
  });

  function getEvaluationResults(traceCount = 1, comparedEventCount = 0, indicatorResultCount = 0): EvaluationHistory {
    return {
      type: 'evaluationHistory',
      triggerEvent: mockEvaluation,
      traces: [...Array(traceCount).keys()].map((idCount) =>
        Trace.fromJSON({
          data: {
            evaluation: {
              // fill every second evaluation with indicatorResults
              indicatorResults: idCount % 2 === 0 ? getIndicatorResults(indicatorResultCount) : [],
              result: ResultTypes.PASSED,
              sloFileContent: '',
              score: 1,
              comparedEvents: [...Array(comparedEventCount).keys()].map((counter) => `myId${counter}`),
            } as unknown as IEvaluationData,
            project: 'myProject',
            stage: 'myStage',
            service: 'myService',
          },
          id: `myId${idCount}`,
          type: EventTypes.EVALUATION_FINISHED,
        })
      ),
    };
  }

  function getIndicatorResults(counter: number): IndicatorResult[] {
    return [...Array(counter).keys()].map((metricCounter) => ({
      value: {
        value: 0,
        metric: `metric${metricCounter}`,
        success: true,
      },
      score: 5,
      status: ResultTypes.PASSED,
      passTargets: [],
      warningTargets: [],
      keySli: false,
    }));
  }
});
