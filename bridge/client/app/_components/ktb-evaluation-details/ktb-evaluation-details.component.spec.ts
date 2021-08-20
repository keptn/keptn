import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { EvaluationsTop10 } from '../../_services/_mockData/evaluations-top10.mock';
import { DataServiceMock } from '../../_services/data.service.mock';
import { Evaluations } from '../../_services/_mockData/evaluations.mock';
import { Trace } from '../../_models/trace';

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationDetailsComponent;
  let fixture: ComponentFixture<KtbEvaluationDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        DataServiceMock,
      ],
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
    component.evaluationData = EvaluationsTop10;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component._heatmapOptions.yAxis[0].categories.length).toEqual(10);
  });

  xit('should have isHeatmapExtendable set to true when more than 10 SLOs are configured ', () => {
    // given
    component.evaluationData = EvaluationsTop10;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component.isHeatmapExtendable).toBeTruthy();
  });

  xit('should have isHeatmapExtendable set to false when less than 10 SLOs are configured', () => {
    // given
    component.evaluationData = Evaluations;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);

    // then
    expect(component.isHeatmapExtendable).toBeFalsy();
  });

  xit('should show a Show all SLIs button when more than 10 SLOs are configured', () => {
    // given
    component.evaluationData = EvaluationsTop10;

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
    component.evaluationData = EvaluationsTop10;

    // when
    component.updateChartData(component.evaluationData.data.evaluationHistory as Trace[]);
    component.toggleHeatmap();

    // then
    expect(component._heatmapOptions.yAxis[0].categories.length).toEqual(17);
  });
});
