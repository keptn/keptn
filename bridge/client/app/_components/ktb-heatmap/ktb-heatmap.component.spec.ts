/* eslint-disable @typescript-eslint/dot-notation */
import Mock = jest.Mock;
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbHeatmapComponent } from './ktb-heatmap.component';
import { EvaluationResultTypeExtension, IDataPoint, IHeatmapTooltipType } from '../../_interfaces/heatmap';
import { DOCUMENT } from '@angular/common';
import { TestUtils } from '../../_utils/test.utils';
import { KtbHeatmapModule } from './ktb-heatmap.module';

describe('KtbHeatmapComponent', () => {
  let component: KtbHeatmapComponent;
  let fixture: ComponentFixture<KtbHeatmapComponent>;
  const elementFromPointSpy: Mock<Element | null, [number, number]> = jest.fn();
  const getComputedTextLengthSpy: Mock<number, [void]> = jest.fn();
  const parentNodeBoundingClientRectSpy: Mock<DOMRect, [void]> = jest.fn();
  const devicePixelRatioSpy: Mock<number, [void]> = jest.fn();
  const outerWidthSpy: Mock<number, [void]> = jest.fn();

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbHeatmapModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbHeatmapComponent);
    component = fixture.componentInstance;
    mockUIElements();
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should group dataPoints', () => {
    // given, when
    component.dataPoints = mockDataPoints(2, 1);

    // then
    const noRefDataPoints = mockDataPoints(2, 1);
    expect(component['groupedData']).toEqual({
      score: [noRefDataPoints[1], noRefDataPoints[3]],
      response_time_p0: [noRefDataPoints[0], noRefDataPoints[2]],
    });
    // verify order of keys
    expect(Object.keys(component['groupedData'])).toEqual(['response_time_p0', 'score']);
  });

  it('should add (#number) to dataPoints if date is not unique', () => {
    // given
    const dataPoints = mockDuplicateDataPoints(1, 2, 2);

    // when
    component.dataPoints = dataPoints;

    // then
    const noRefDataPoints = mockDuplicateDataPoints(1, 2, 2);

    expect(dataPoints).not.toEqual(noRefDataPoints);
    expect(component['groupedData']).toEqual({
      response_time_p0: [
        mockSliDataPoint('myEvaluation0duplicate00', 'myDate0duplicate0', 'response_time_p0'),
        mockSliDataPoint('myEvaluation0duplicate01', 'myDate0duplicate0 (1)', 'response_time_p0'),
        mockSliDataPoint('myEvaluation0duplicate02', 'myDate0duplicate0 (2)', 'response_time_p0'),
        mockSliDataPoint('myEvaluation0duplicate10', 'myDate0duplicate1', 'response_time_p0'),
        mockSliDataPoint('myEvaluation0duplicate11', 'myDate0duplicate1 (1)', 'response_time_p0'),
        mockSliDataPoint('myEvaluation0duplicate12', 'myDate0duplicate1 (2)', 'response_time_p0'),
      ],
      score: [
        mockScoreDataPoint('myEvaluation0duplicate00', 'myDate0duplicate0'),
        mockScoreDataPoint('myEvaluation0duplicate01', 'myDate0duplicate0 (1)'),
        mockScoreDataPoint('myEvaluation0duplicate02', 'myDate0duplicate0 (2)'),
        mockScoreDataPoint('myEvaluation0duplicate10', 'myDate0duplicate1'),
        mockScoreDataPoint('myEvaluation0duplicate11', 'myDate0duplicate1 (1)'),
        mockScoreDataPoint('myEvaluation0duplicate12', 'myDate0duplicate1 (2)'),
      ],
    });
  });

  it('should not add any number to dataPoints if date is unique', () => {
    // given
    const dataPoints = mockDataPoints(2, 1);

    // when
    component.dataPoints = dataPoints;

    // then
    const noRefDataPoints = mockDataPoints(2, 1);

    expect(dataPoints).toEqual(noRefDataPoints);
  });

  it('should correctly save mouse movement', () => {
    // given
    devicePixelRatioSpy.mockReturnValue(2);

    // when
    component['onMouseMove'](
      new MouseEvent('move', {
        clientX: 2,
        clientY: 4,
      })
    );

    // then
    expect(component['mouseCoordinates']).toEqual({
      x: 4,
      y: 8,
    });
  });

  it('should not select dataPoint if it is not in the dataSource', () => {
    component.dataPoints = [];
    component.selectedIdentifier = 'myIdentifier';
    expect(component.selectedIdentifier).toBe('myIdentifier');
    expect(component['_selectedDataPoint']).toBeUndefined();
  });

  it('should select dataPoint and not emit it if it is preselected', () => {
    // given
    const dataPoints = mockDataPoints(1, 0);
    const emitSpy = jest.spyOn(component.selectedIdentifierChange, 'emit');
    component.dataPoints = dataPoints;

    // when
    component.selectedIdentifier = dataPoints[0].identifier;

    // then
    expect(component.selectedIdentifier).toEqual(dataPoints[0].identifier);
    expect(emitSpy).not.toHaveBeenCalled();
  });

  it('should select dataPoint with previously set identifier that was not found', () => {
    // given
    const emitSpy = jest.spyOn(component.selectedIdentifierChange, 'emit');

    // when
    component.selectedIdentifier = 'myEvaluation0';
    const dataPoints = mockDataPoints(1, 0);
    component.dataPoints = dataPoints;

    // then
    expect(component.selectedIdentifier).toEqual(dataPoints[0].identifier);
    expect(component['_selectedDataPoint']).toEqual(dataPoints[0]);
    expect(emitSpy).not.toHaveBeenCalled();
  });

  it('should select dataPoint and emit it if it is preselected', () => {
    // given
    const dataPoints = mockDataPoints(1, 0);
    const emitSpy = jest.spyOn(component.selectedIdentifierChange, 'emit');
    component.dataPoints = dataPoints;

    // when
    component['click'](dataPoints[0], false);

    // then
    expect(component.selectedIdentifier).toEqual(dataPoints[0].identifier);
    expect(emitSpy).toHaveBeenCalled();
  });

  it('should set correct height', () => {
    // when
    component['setHeight'](5);

    // then
    expect(component['height']).toBe(150);
  });

  it('should remove mouseOver listener on destroy', () => {
    const listener = component['mouseMoveListener'];
    const removeEventListenerSpy = jest.spyOn(TestBed.inject(DOCUMENT), 'removeEventListener');

    // when
    component.ngOnDestroy();

    // then
    expect(removeEventListenerSpy).toHaveBeenCalledWith('mousemove', listener);
  });

  it('should correctly set showMoreButton-Style', () => {
    // given
    component.dataPoints = mockDataPoints(1, 10);
    component['height'] = 100;
    component['dataPointContentWidth'] = 200;

    // when
    component['resizeShowMoreButton']();

    // then
    const button: HTMLElement = component.showMoreButton._elementRef.nativeElement;

    expect(button.style.top).toBe('103px');
    expect(button.style.left).toBe('150px');
    expect(button.style.width).toBe('200px');
  });

  it('should not set showMoreButton-Style if button is not visible', () => {
    // given
    component.dataPoints = mockDataPoints(1, 9);
    component['height'] = 100;
    component['dataPointContentWidth'] = 200;

    // when
    component['resizeShowMoreButton']();

    // then
    const button: HTMLElement = component.showMoreButton._elementRef.nativeElement;

    expect(button.style.top).toBe('');
    expect(button.style.left).toBe('');
    expect(button.style.width).toBe('');
  });

  it('should return highlight width', () => {
    // just test 3 different widths to be set correctly
    for (let i = 1; i < 4; ++i) {
      // given
      component.dataPoints = mockDataPoints(i, 9);

      // when
      const width = component['getHighlightWidth']();

      // then
      expect(width).toBe(850 / i);
    }
  });

  /**
   * Mocks and adds a spy to:
   * <br/>- SVGElement.getComputedTextLength()
   * <br/>- document.elementFromPoint()
   * <br/>- Width of parentNode
   * <br/>- Ratio of window
   * <br/>- Width of window
   */
  function mockUIElements(): void {
    const document = TestBed.inject(DOCUMENT);
    elementFromPointSpy.mockReturnValue(null);
    document.elementFromPoint = elementFromPointSpy;

    getComputedTextLengthSpy.mockReturnValue(100);
    Object.defineProperty(SVGElement.prototype, 'getComputedTextLength', {
      value: getComputedTextLengthSpy,
      writable: false,
    });
    parentNodeBoundingClientRectSpy.mockReturnValue(getDomRect(1000));
    const htmlElement: HTMLElement = fixture.nativeElement.parentNode;
    jest.spyOn(htmlElement, 'getBoundingClientRect').mockImplementation(parentNodeBoundingClientRectSpy);

    devicePixelRatioSpy.mockReturnValue(1);

    TestUtils.overridePropertyWithSpy(window, 'devicePixelRatio', devicePixelRatioSpy);
    TestUtils.overridePropertyWithSpy(window, 'outerWidth', outerWidthSpy);
  }

  function getDomRect(width: number): DOMRect {
    return {
      width,
      y: 0,
      right: 0,
      bottom: 0,
      left: 0,
      x: 0,
      height: 0,
      top: 0,
      toJSON(): object {
        return {};
      },
    };
  }
});

export function mockDataPoints(
  counter: number,
  slis: number,
  identifierSuffix = '',
  dateSuffix = '',
  mockIdentifiers = false
): IDataPoint[] {
  const dataPoints: IDataPoint[] = [];
  for (let i = 0; i < counter; ++i) {
    const identifier = `myEvaluation${i}${identifierSuffix}`;
    const identifierBefore = mockIdentifiers ? `myEvaluation${i - 1}${identifierSuffix}` : undefined;
    const xElement = `myDate${i}${dateSuffix}`;
    for (let y = 0; y < slis; ++y) {
      dataPoints.push(mockSliDataPoint(identifier, xElement, `response_time_p${y}`, identifierBefore));
    }
    dataPoints.push(mockScoreDataPoint(identifier, xElement, identifierBefore));
  }
  return dataPoints;
}

function mockDuplicateDataPoints(slis: number, duplicates: number, duplicatesPerDate: number): IDataPoint[] {
  // duplicates: how many different (dates with) duplicates
  const dataPoints: IDataPoint[] = [];
  for (let i = 0; i < duplicates; ++i) {
    for (let y = 0; y < duplicatesPerDate + 1; ++y) {
      const duplicate = mockDataPoints(1, slis, `duplicate${i}${y}`, `duplicate${i}`);
      dataPoints.push(...duplicate);
    }
  }
  return dataPoints;
}

function mockScoreDataPoint(identifier: string, xElement: string, identifierBefore?: string): IDataPoint {
  return {
    comparedIdentifier: identifierBefore ? [identifierBefore] : [],
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
    yElement: 'score',
  };
}

function mockSliDataPoint(
  identifier: string,
  xElement: string,
  yElement: string,
  identifierBefore?: string
): IDataPoint {
  return {
    comparedIdentifier: identifierBefore ? [identifierBefore] : [],
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
/* eslint-enable @typescript-eslint/dot-notation */
