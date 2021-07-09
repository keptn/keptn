import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';
import { KtbSliBreakdownComponent } from './ktb-sli-breakdown.component';
import { KtbEvaluationDetailsComponent } from '../ktb-evaluation-details/ktb-evaluation-details.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../../app.module';
import { Evaluations } from '../../_services/_mockData/evaluations.mock';

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbSliBreakdownComponent;
  let fixture: ComponentFixture<KtbSliBreakdownComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSliBreakdownComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  it('should have expandable row', () => {
    // given
    const selectedEvaluation = Evaluations.data.evaluationHistory[1];
    component.indicatorResults = selectedEvaluation.data.evaluation.indicatorResults;
    component.score = selectedEvaluation.data.evaluation.score;
    component.comparedIndicatorResults = Evaluations.data.evaluationHistory[0].data.evaluation.indicatorResults;

    // then
    fixture.detectChanges();
    const rows = fixture.nativeElement.querySelectorAll('dt-row');
    const rowBefore = rows[0].innerText;
    const cells = rows[0].querySelectorAll('dt-cell');
    const firstCell = cells[1];

    expect(rows.length).toBe(1);
    expect(fixture.nativeElement.querySelector('dt-table')).toBeTruthy();
    expect(cells[0].querySelector('button')).toBeTruthy();
    expect(firstCell.innerText).not.toContain('compared with');
    expect(firstCell.innerText).toContain('response_time_p95');

    rows[0].click();
    fixture.detectChanges();
    expect(firstCell.innerText).toContain('compared with');

    rows[0].click();
    fixture.detectChanges();
    expect(firstCell.innerText).not.toContain('compared with');
    expect(rowBefore).toBe(rows[0].innerText);
  });

  it('should not have expandable row', () => {
    // given
    const selectedEvaluation = Evaluations.data.evaluationHistory[1];
    component.indicatorResults = selectedEvaluation.data.evaluation.indicatorResults;
    component.score = selectedEvaluation.data.evaluation.score;
    component.comparedIndicatorResults = [];

    // then
    fixture.detectChanges();
    const rows = fixture.nativeElement.querySelectorAll('dt-row');
    const rowBefore = rows[0].innerText;
    const cells = rows[0].querySelectorAll('dt-cell');

    expect(cells[0].querySelector('button')).toBeFalsy();

    rows[0].click();
    fixture.detectChanges();
    expect(rowBefore).toBe(rows[0].innerText);
  });

  it('should have success values', () => {
    // given
    const selectedEvaluation = Evaluations.data.evaluationHistory[7];
    component.indicatorResults = selectedEvaluation.data.evaluation.indicatorResults;
    component.score = selectedEvaluation.data.evaluation.score;
    component.comparedIndicatorResults = Evaluations.data.evaluationHistory[5].data.evaluation.indicatorResults;

    // then
    fixture.detectChanges();
    const firstRow = fixture.nativeElement.querySelectorAll('dt-row')[0];
    firstRow.click();
    fixture.detectChanges();

    const cells = firstRow.querySelectorAll('dt-cell');
    const values = cells[2].querySelectorAll('span.success');

    expect(values.length).toBe(2);
    expect(cells[2].innerText).toContain('370.2');
    expect(cells[2].innerText).toContain('334.5');
    expect(values[0].innerText).toBe('+35.65');
    expect(values[1].innerText).toBe('+10.65%');
    expect(cells[3].innerText).toBe('<=+10% and <600');
    expect(cells[4].innerText).toBe('<=800');
    expect(cells[5].innerText).toBe('passed');
    expect(cells[6].innerText).toBe('100');
    expect(firstRow.querySelector('.error, .error-line')).toBeFalsy();
  });

  it('should have error values', () => {
    // given
    const selectedEvaluation = Evaluations.data.evaluationHistory[6];
    component.indicatorResults = selectedEvaluation.data.evaluation.indicatorResults;
    component.score = selectedEvaluation.data.evaluation.score;
    component.comparedIndicatorResults = Evaluations.data.evaluationHistory[5].data.evaluation.indicatorResults;

    // then
    fixture.detectChanges();
    const firstRow = fixture.nativeElement.querySelectorAll('dt-row')[0];
    firstRow.click();
    fixture.detectChanges();

    const cells = firstRow.querySelectorAll('dt-cell');
    const values = cells[2].querySelectorAll('span.error');

    expect(values.length).toBe(2);
    expect(cells[2].innerText).toContain('370.2');
    expect(cells[2].innerText).toContain('1082');
    expect(values[0].innerText).toBe('-712.42');
    expect(values[1].innerText).toBe('-65.805%');
    expect(cells[3].innerText).toBe('<=+10% and <600');
    expect(cells[4].innerText).toBe('<=800');
    expect(cells[5].innerText).toBe('failed');
    expect(cells[6].innerText).toBe('0');
    expect(cells[3].querySelectorAll('.error.error-line').length).toBe(2);
    expect(cells[4].querySelectorAll('.error.error-line').length).toBe(1);
    expect(cells[6].querySelector('.error')).toBeTruthy();
    expect(firstRow.querySelector('.success')).toBeFalsy();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
