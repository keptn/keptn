import { mockDataPoints } from './ktb-heatmap.component.spec';
import {
  calculateTooltipPosition,
  createGroupedDataPoints,
  findXElementThroughIdentifier,
  getAvailableIdentifiers,
  getAxisElements,
  getHiddenYElements,
  getLimitedYElements,
  getXAxisReducedElements,
} from './ktb-heatmap-utils';
import { GroupedDataPoints } from '../../_interfaces/heatmap';
import { TestUtils } from '../../_utils/test.utils';
import Mock = jest.Mock;

describe('KtbHeatmapUtils', () => {
  const elementFromPointSpy: Mock<Element | null, [number, number]> = jest.fn();
  const devicePixelRatioSpy: Mock<number, [void]> = jest.fn();
  const outerWidthSpy: Mock<number, [void]> = jest.fn();

  beforeEach(async () => {
    mockUIElements();
  });

  it('should group dataPoints', () => {
    // given, when
    const groupedDataPoints = createGroupedDataPoints(mockDataPoints(2, 1));

    // then
    const noRefDataPoints = mockDataPoints(2, 1);
    expect(groupedDataPoints).toEqual({
      score: [noRefDataPoints[1], noRefDataPoints[3]],
      response_time_p0: [noRefDataPoints[0], noRefDataPoints[2]],
    });
    // verify order of keys
    expect(Object.keys(groupedDataPoints)).toEqual(['response_time_p0', 'score']);
  });

  it('should correctly return axis elements without duplicates', () => {
    const groupedData = mockGroupedData(2, 1);
    expect(getAxisElements(groupedData, 10)).toEqual({
      xElements: ['myDate0', 'myDate1'],
      yElements: ['response_time_p0', 'score'],
      showMoreVisible: false,
    });
  });

  it('should return a reduced set of elements (even index) for xAxis labels if the width is too small', () => {
    // given, when
    const reducedDates = getXAxisReducedElements(generateArray(3), 50, 25);

    // then
    expect(reducedDates).toEqual(['0', '2']);
  });

  it('should return a reduced set of elements (odd index) for xAxis labels if the width is too small', () => {
    // given, when
    const reducedDates = getXAxisReducedElements(generateArray(4), 50, 25);

    // then
    expect(reducedDates).toEqual(['1', '3']);
  });

  it('should not return all elements for xAxis if width is enough', () => {
    // given, when
    const reducedDates = getXAxisReducedElements(generateArray(3), 150, 25);

    // then
    expect(reducedDates).toEqual(['0', '1', '2']);
  });

  it('should return only available identifiers', () => {
    // given
    const groupedData = mockGroupedData(2, 1, '', '', true);

    // when
    const foundIdentifiers = getAvailableIdentifiers(['myEvaluation-1', 'myEvaluation0'], groupedData);

    // then
    expect(foundIdentifiers).toEqual(['myEvaluation0']);
  });

  it('should find xElement through identifier', () => {
    // given
    const groupedData = mockGroupedData(3, 1, '', '', true);

    // when
    const foundXElement = findXElementThroughIdentifier('myEvaluation1', groupedData);

    // then
    expect(foundXElement).toEqual('myDate1');
  });

  it('should not find xElement through identifier', () => {
    // given
    const groupedData = mockGroupedData(3, 1, '', '', true);

    // when
    const foundXElement = findXElementThroughIdentifier('notFound', groupedData);

    // then
    expect(foundXElement).toBeUndefined();
  });

  it('should return hidden yElements', () => {
    // when
    const hiddenElements = getHiddenYElements(generateArray(12), 10);

    // then
    expect(hiddenElements).toEqual(['0', '1']);
  });

  it('should not return hidden yElements', () => {
    // when
    const hiddenElements = getHiddenYElements(generateArray(10), 10);

    // then
    expect(hiddenElements).toEqual([]);
  });

  it('should return all yElements if yElements is <= 10', () => {
    const myArray: string[] = [];
    for (let i = 0; i < 10; ++i) {
      myArray.push(i.toString());

      // when
      const yElements = getLimitedYElements(myArray, 10);

      // then
      expect(yElements).toEqual(myArray);
    }
  });

  it('should return limited set of yElements', () => {
    // when
    const hiddenElements = getLimitedYElements(generateArray(12), 10);

    // then
    expect(hiddenElements).toEqual(generateArray(10, 2));
  });

  it('should show tooltip on the left if scrollbar is visible', () => {
    // given
    devicePixelRatioSpy.mockReturnValue(2);
    outerWidthSpy.mockReturnValue(1000);

    // when
    const coordinates = calculateTooltipPosition(100, 20, 388, 250);

    // then
    expect(coordinates).toEqual({
      top: 255,
      left: 293,
    });
  });

  it('should show tooltip on the right if scrollbar is not visible', () => {
    // given
    devicePixelRatioSpy.mockReturnValue(2);
    outerWidthSpy.mockReturnValue(1000);

    // when
    const coordinates = calculateTooltipPosition(100, 0, 388, 250);

    // then
    expect(coordinates).toEqual({
      top: 255,
      left: 393,
    });
  });

  /**
   * Mocks and adds a spy to:
   * <br/>- document.elementFromPoint()
   * <br/>- Ratio of window
   * <br/>- Width of window
   */
  function mockUIElements(): void {
    elementFromPointSpy.mockReturnValue(null);
    document.elementFromPoint = elementFromPointSpy;

    devicePixelRatioSpy.mockReturnValue(1);
    TestUtils.overridePropertyWithSpy(window, 'devicePixelRatio', devicePixelRatioSpy);
    TestUtils.overridePropertyWithSpy(window, 'outerWidth', outerWidthSpy);
  }

  function generateArray(counter: number, offset = 0): string[] {
    const array: string[] = [];
    for (let i = 0; i < counter; ++i) {
      array.push(`${i + offset}`);
    }
    return array;
  }

  function mockGroupedData(
    counter: number,
    slis: number,
    identifierSuffix = '',
    dateSuffix = '',
    mockIdentifiers = false
  ): GroupedDataPoints {
    return createGroupedDataPoints(mockDataPoints(counter, slis, identifierSuffix, dateSuffix, mockIdentifiers));
  }
});
