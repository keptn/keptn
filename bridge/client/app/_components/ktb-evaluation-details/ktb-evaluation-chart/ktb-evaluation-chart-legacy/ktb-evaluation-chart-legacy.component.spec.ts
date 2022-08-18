import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEvaluationChartLegacyComponent } from './ktb-evaluation-chart-legacy.component';
import { KtbEvaluationDetailsModule } from '../../ktb-evaluation-details.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EvaluationChartItemMock } from '../../../../_services/_mockData/evaluation-chart-item.mock';
import { SliInfoMock } from '../../../../_services/_mockData/sli-info.mock';

describe(KtbEvaluationChartLegacyComponent.name, () => {
  let component: KtbEvaluationChartLegacyComponent;
  let fixture: ComponentFixture<KtbEvaluationChartLegacyComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEvaluationDetailsModule, HttpClientTestingModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbEvaluationChartLegacyComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should aggregate SLI results', () => {
    const result = Reflect.get(component, 'getSliResultInfos')(EvaluationChartItemMock);
    expect(result).toEqual(SliInfoMock);
  });
});
