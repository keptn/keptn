import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationChartComponent } from './ktb-evaluation-chart.component';

describe('KtbEvaluationChartComponent', () => {
  let component: KtbEvaluationChartComponent;
  let fixture: ComponentFixture<KtbEvaluationChartComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbEvaluationChartComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEvaluationChartComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
