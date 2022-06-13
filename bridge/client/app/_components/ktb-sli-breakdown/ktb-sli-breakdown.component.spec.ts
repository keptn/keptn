import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import { Trace } from '../../_models/trace';
import { IndicatorResult } from '../../../../shared/interfaces/indicator-result';
import { SloConfig } from '../../../../shared/interfaces/slo-config';
import { parse as parseYaml } from 'yaml';
import { KtbSliBreakdownModule } from './ktb-sli-breakdown.module';

describe('KtbSliBreakdownComponent', () => {
  let component: KtbSliBreakdownComponent;
  let fixture: ComponentFixture<KtbSliBreakdownComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSliBreakdownModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSliBreakdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should assemble table with all details', () => {
    // given
    initEvaluation(11);

    // then
    expect(component.tableEntries.data.length).toBe(1);
    expect(component.tableEntries.data[0]).toEqual({
      name: 'response_time_p95',
      value: 315.8,
      result: 'pass',
      score: 100,
      passTargets: [
        {
          criteria: '<=+10%',
          targetValue: 392.5944454106811,
          violated: false,
        },
        {
          criteria: '<600',
          targetValue: 600,
          violated: false,
        },
      ],
      warningTargets: [
        {
          criteria: '<=800',
          targetValue: 800,
          violated: false,
        },
      ],
      keySli: false,
      success: true,
      expanded: false,
      weight: 1,
      comparedValue: 356.9,
      calculatedChanges: {
        absolute: -41.03,
        relative: -11.49,
      },
    });
  });

  it('should provide compare value from payload', () => {
    // given
    const spy = jest.spyOn(component, 'calculateComparedValue');
    initEvaluation(11);

    // then
    expect(spy).not.toHaveBeenCalled();
    expect(component.tableEntries.data[0].comparedValue).toBe(356.9);
  });

  it('should retrieve compared value from comparedIndicatorResult if not in payload', () => {
    // given
    const spy = jest.spyOn(component, 'calculateComparedValue');
    initEvaluation(10);

    // then
    expect(spy).toHaveBeenCalled();
    expect(component.tableEntries.data[0].comparedValue).toBe(365.2);
  });

  it('should have weight fallback to 1', () => {
    // given
    initEvaluation(0);

    // then
    expect(component.objectives?.[0].weight).toBeUndefined();
    expect(component.tableEntries.data[0].weight).toBe(1);
  });

  function initEvaluation(selectedEvaluationIndex: number): void {
    const selectedEvaluation = EvaluationsMock.data.evaluationHistory?.[selectedEvaluationIndex] as Trace;
    component.indicatorResults = selectedEvaluation.data.evaluation?.indicatorResults as IndicatorResult[];

    const sloFileContentParsed = parseYaml(atob(selectedEvaluation.data.evaluation?.sloFileContent ?? '')) as SloConfig;
    component.objectives = sloFileContentParsed.objectives;
    component.score = selectedEvaluation.data.evaluation?.score as number;

    component.comparedIndicatorResults =
      EvaluationsMock.data.evaluationHistory
        ?.filter(
          (evaluation: Trace) => selectedEvaluation.data.evaluation?.comparedEvents?.indexOf(evaluation.id) !== -1
        )
        .map((evaluation: Trace) => evaluation.data.evaluation?.indicatorResults as IndicatorResult[]) ?? [];
  }
});
