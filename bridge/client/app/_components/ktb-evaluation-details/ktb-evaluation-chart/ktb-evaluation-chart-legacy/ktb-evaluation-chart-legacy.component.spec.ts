import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationChartLegacyComponent } from './ktb-evaluation-chart-legacy.component';

describe('KtbEvaluationChartLegacyComponent', () => {
  let component: KtbEvaluationChartLegacyComponent;
  let fixture: ComponentFixture<KtbEvaluationChartLegacyComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbEvaluationChartLegacyComponent ]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEvaluationChartLegacyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
