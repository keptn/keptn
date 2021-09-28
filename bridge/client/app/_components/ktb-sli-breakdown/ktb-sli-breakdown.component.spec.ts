import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import { KtbEvaluationDetailsComponent } from '../ktb-evaluation-details/ktb-evaluation-details.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { Evaluations } from '../../_services/_mockData/evaluations.mock';

enum Column {
  DETAILS = 0,
  NAME = 1,
  VALUE = 2,
  WEIGHT = 3,
  PASS_CRITERIA = 4,
  WARNING_CRITERIA = 5,
  RESULT = 6,
  SCORE = 7,
}

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbSliBreakdownComponent;
  let fixture: ComponentFixture<KtbSliBreakdownComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSliBreakdownComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should have expandable row', () => {
    // given
    initEvaluation(1, 0);

    // then
    fixture.detectChanges();
    const rows = fixture.nativeElement.querySelectorAll('dt-row');
    const rowBefore = rows[0].textContent;
    const cells = rows[0].querySelectorAll('dt-cell');
    const firstCell = cells[Column.NAME];

    expect(rows.length).toBe(1);
    expect(fixture.nativeElement.querySelector('dt-table')).toBeTruthy();
    expect(cells[0].querySelector('button')).toBeTruthy();
    expect(firstCell.textContent).not.toContain('compared with');
    expect(firstCell.textContent).toContain('response_time_p95');

    rows[0].click();
    fixture.detectChanges();
    expect(firstCell.textContent).toContain('compared with');

    rows[0].click();
    fixture.detectChanges();
    expect(firstCell.textContent).not.toContain('compared with');
    expect(rowBefore).toBe(rows[0].textContent);
  });

  it('should not have expandable row', () => {
    // given
    initEvaluation(1);

    // then
    fixture.detectChanges();
    const rows = fixture.nativeElement.querySelectorAll('dt-row');
    const rowBefore = rows[0].textContent;
    const cells = rows[0].querySelectorAll('dt-cell');

    expect(cells[Column.DETAILS].querySelector('button')).toBeFalsy();

    rows[0].click();
    fixture.detectChanges();
    expect(rowBefore).toBe(rows[0].textContent);
  });

  it('should have success values', () => {
    // given
    initEvaluation(7, 5);

    // when
    fixture.detectChanges();
    const firstRow = fixture.nativeElement.querySelectorAll('dt-row')[0];
    firstRow.click();
    fixture.detectChanges();

    // then
    const cells = firstRow.querySelectorAll('dt-cell');
    validateIndicatorResult(
      cells,
      true,
      '370.2',
      '334.5',
      '1',
      '+35.65',
      '+10.65%',
      '<=+10% and <600',
      '<=800',
      'passed',
      '100'
    );

    expect(firstRow.querySelector('.error, .error-line')).toBeFalsy();
  });

  it('should have error values', () => {
    // given
    initEvaluation(6, 5);

    // when
    fixture.detectChanges();
    const firstRow = fixture.nativeElement.querySelectorAll('dt-row')[0];
    firstRow.click();
    fixture.detectChanges();

    // then
    const cells = firstRow.querySelectorAll('dt-cell');
    validateIndicatorResult(
      cells,
      false,
      '370.2',
      '1082',
      '0',
      '-712.42',
      '-65.805%',
      '<=+10% and <600',
      '<=800',
      'failed',
      '0'
    );

    expect(cells[Column.PASS_CRITERIA].querySelectorAll('.error.error-line').length).toBe(2);
    expect(cells[Column.WARNING_CRITERIA].querySelectorAll('.error.error-line').length).toBe(1);
    expect(cells[Column.SCORE].querySelector('.error')).toBeTruthy();
    expect(firstRow.querySelector('.success')).toBeFalsy();
  });

  it('should sort by weight asc', () => {
    validateOrder(0, Column.WEIGHT, true, 0, 2, 1);
  });

  it('should sort by weight desc', () => {
    validateOrder(0, Column.WEIGHT, false, 1, 2, 0);
  });

  it('should sort by name asc', () => {
    validateOrder(0, Column.NAME, true, 2, 1, 0);
  });

  it('should sort by name desc', () => {
    validateOrder(0, Column.NAME, false, 0, 1, 2);
  });

  it('should sort by score asc', () => {
    validateOrder(0, Column.SCORE, true, 1, 2, 0);
  });

  it('should sort by score desc', () => {
    validateOrder(0, Column.SCORE, false, 0, 2, 1);
  });

  function validateOrder(selectedEvaluationIndex: number, column: Column, isAsc: boolean, ...indices: number[]) {
    // given
    initEvaluation(selectedEvaluationIndex);
    fixture.detectChanges();

    // when
    for (let i = isAsc ? 1 : 0; i < 2; ++i) {
      fixture.nativeElement.querySelectorAll('dt-header-cell')[column].click();
      fixture.detectChanges();
    }
    // then
    // @ts-ignore
    const selectedEvaluation = Evaluations.data.evaluationHistory[selectedEvaluationIndex];
    // @ts-ignore
    const indicatorNames = fixture.nativeElement.querySelectorAll(`dt-row > dt-cell:nth-child(${Column.NAME + 1})`);
    for (let i = 0; i < indices.length; ++i) {
      // @ts-ignore
      expect(indicatorNames[i].textContent).toEqual(
        selectedEvaluation.data.evaluation.indicatorResults[indices[i]].value.metric
      );
    }
  }

  function initEvaluation(selectedEvaluationIndex: number, comparedEvaluationIndex: number = -1) {
    // @ts-ignore
    const selectedEvaluation = Evaluations.data.evaluationHistory[selectedEvaluationIndex];
    // @ts-ignore
    component.indicatorResults = selectedEvaluation.data.evaluation.indicatorResults;
    // @ts-ignore
    component.score = selectedEvaluation.data.evaluation.score;
    // @ts-ignore
    component.comparedIndicatorResults =
      comparedEvaluationIndex === -1
        ? []
        : Evaluations.data.evaluationHistory[comparedEvaluationIndex].data.evaluation.indicatorResults;
  }

  function validateIndicatorResult(
    cells: HTMLElement[],
    isSuccess: boolean,
    firstValue: string,
    secondValue: string,
    weight: string,
    comparedValueAbsolute: string,
    comparedValueRelative: string,
    passCriteria: string,
    warningCriteria: string,
    result: string,
    score: string
  ) {
    const calculatedValues: NodeListOf<HTMLElement> = cells[Column.VALUE].querySelectorAll(
      `span.${isSuccess ? 'success' : 'error'}`
    );

    expect(calculatedValues.length).toBe(2);
    expect(cells[Column.VALUE].textContent).toContain(firstValue);
    expect(cells[Column.VALUE].textContent).toContain(secondValue);
    expect(cells[Column.WEIGHT].textContent).toContain(weight);
    expect(calculatedValues[0].textContent).toBe(comparedValueAbsolute);
    expect(calculatedValues[1].textContent).toBe(comparedValueRelative);
    expect(cells[Column.PASS_CRITERIA].textContent?.replace(/\s/g, '')).toBe(passCriteria.replace(/\s/g, ''));
    expect(cells[Column.WARNING_CRITERIA].textContent?.replace(/\s/g, '')).toBe(warningCriteria.replace(/\s/g, ''));
    expect(cells[Column.RESULT].textContent).toBe(result);
    expect(cells[Column.SCORE].textContent).toBe(score);
  }
});
