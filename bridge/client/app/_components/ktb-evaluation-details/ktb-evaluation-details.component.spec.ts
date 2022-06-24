import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EvaluationsTop10Mock } from '../../_services/_mockData/evaluations-top10.mock';
import { EvaluationsMock } from '../../_services/_mockData/evaluations.mock';
import { Trace } from '../../_models/trace';
import { EvaluationChartItemMock } from '../../_services/_mockData/evaluation-chart-item.mock';
import { SliInfoMock } from '../../_services/_mockData/sli-info.mock';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbEvaluationDetailsModule } from './ktb-evaluation-details.module';

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationDetailsComponent;
  let fixture: ComponentFixture<KtbEvaluationDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbEvaluationDetailsModule, HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEvaluationDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  xit('should have a reduced heatmap size when more than 10 SLOs are configured', () => {
    // given
    component.evaluationData = EvaluationsTop10Mock;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component._heatmapOptions.yAxis[0].categories.length).toEqual(10);
  });

  xit('should have isHeatmapExtendable set to true when more than 10 SLOs are configured ', () => {
    // given
    component.evaluationData = EvaluationsTop10Mock;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component.isHeatmapExtendable).toBeTruthy();
  });

  xit('should have isHeatmapExtendable set to false when less than 10 SLOs are configured', () => {
    // given
    component.evaluationData = EvaluationsMock;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component.isHeatmapExtendable).toBeFalsy();
  });

  xit('should show a Show all SLIs button when more than 10 SLOs are configured', () => {
    // given
    component.evaluationData = EvaluationsTop10Mock;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);
    fixture.detectChanges();
    const evaluationPage = fixture.nativeElement;
    const button = evaluationPage.querySelector('button.button-show-more-slo');

    // then
    expect(button).toBeTruthy();
  });

  xit('should have a full heatmap size when more than 10 SLOs are configured and toggle is triggered', () => {
    // given
    component.evaluationData = EvaluationsTop10Mock;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);
    component.toggleHeatmap();

    // then
    expect(component._heatmapOptions.yAxis[0].categories.length).toEqual(17);
  });

  it('should aggregate SLI results', () => {
    // eslint-disable-next-line @typescript-eslint/ban-ts-comment
    // @ts-ignore
    const result = component.getSliResultInfos(EvaluationChartItemMock);
    expect(result).toEqual(SliInfoMock);
  });
});
