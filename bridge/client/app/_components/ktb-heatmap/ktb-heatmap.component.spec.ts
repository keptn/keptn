// eslint-disable-next-line @typescript-eslint/ban-ts-comment
// @ts-nocheck : private access
import Mock = jest.Mock;
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHeatmapComponent } from './ktb-heatmap.component';
import { AppModule } from '../../app.module';
import { EvaluationResultTypeExtension, IDataPoint, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { DOCUMENT } from '@angular/common';

describe('KtbHeatmapComponent', () => {
  let component: KtbHeatmapComponent;
  let fixture: ComponentFixture<KtbHeatmapComponent>;
  let elementFromPointSpy: Mock<Element | null, [number, number]>;
  let getComputedTextLengthSpy: Mock<number, [void]>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();
    mockUIElements();
    fixture = TestBed.createComponent(KtbHeatmapComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should group dataPoints', () => {
    const dataPoints = mockDataPoints(2, 1);
    component.dataPoints = dataPoints;
    expect(component.groupedData).toEqual({
      score: [dataPoints[0], dataPoints[2]],
      response_time_p0: [dataPoints[1], dataPoints[3]],
    });
  });

  function mockDataPoints(counter: number, slis: number): IDataPoint[] {
    const dataPoints: IDataPoint[] = [];
    for (let i = 0; i < counter; ++i) {
      const identifier = `myEvaluation${i}`;
      const xElement = `myDate${i}`;
      dataPoints.push(mockScoreDataPoint(identifier, xElement, 'score'));
      for (let y = 0; y < slis; ++y) {
        dataPoints.push(mockSliDataPoint(identifier, xElement, `response_time_p${y}`));
      }
    }
    return dataPoints;
  }

  function mockScoreDataPoint(identifier: string, xElement: string, yElement: string): IDataPoint {
    return {
      comparedIdentifier: [],
      identifier,
      tooltip: {
        warningCount: 0,
        thresholdWarn: 0,
        thresholdPass: 0,
        passCount: 0,
        warn: false,
        value: 0,
        type: IHeatmapTooltipType.SCORE,
        failedCount: 0,
        fail: true,
      },
      color: EvaluationResultTypeExtension.INFO,
      xElement,
      yElement,
    };
  }

  function mockSliDataPoint(identifier: string, xElement: string, yElement: string): IDataPoint {
    return {
      comparedIdentifier: [],
      identifier,
      tooltip: {
        type: IHeatmapTooltipType.SLI,
        keySli: false,
        score: 0,
        warningTargets: [],
        passTargets: [],
        value: 0,
      },
      color: EvaluationResultTypeExtension.INFO,
      xElement,
      yElement,
    };
  }

  /**
   * Mocks and adds a spy to:
   * <br/>- svgElement.getComputedTextLength()
   * <br/>- document.elementFromPoint()
   */
  function mockUIElements(): void {
    const document = TestBed.inject(DOCUMENT);
    elementFromPointSpy = jest.fn();
    elementFromPointSpy.mockReturnValue(null);
    document.elementFromPoint = elementFromPointSpy;

    getComputedTextLengthSpy = jest.fn();
    getComputedTextLengthSpy.mockReturnValue(100);
    Object.defineProperty(SVGElement.prototype, 'getComputedTextLength', {
      value: getComputedTextLengthSpy,
      writable: false,
    });
  }
});
