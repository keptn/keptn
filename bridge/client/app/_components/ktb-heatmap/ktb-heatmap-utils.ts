import { GroupedDataPoints, IDataPoint } from '../../_interfaces/heatmap';

export function createGroupedDataPoints(data: IDataPoint[]): GroupedDataPoints {
  setUniqueHeaders(data, 'xElement', 'yElement');
  setUniqueHeaders(data, 'yElement', 'xElement');
  return data.reduce((groupedData: GroupedDataPoints, dataPoint) => {
    groupedData[dataPoint.yElement] ||= [];
    groupedData[dataPoint.yElement].push(dataPoint);
    return groupedData;
  }, {});
}

export const getLimitedYElements = (yElements: string[], limit: number): string[] => {
  if (yElements.length <= limit) {
    return yElements;
  }
  return yElements.slice(yElements.length - limit, yElements.length);
};

export function setUniqueHeaders(dataPoints: IDataPoint[], key: 'xElement', compare: 'yElement'): void;
export function setUniqueHeaders(dataPoints: IDataPoint[], key: 'yElement', compare: 'xElement'): void;
export function setUniqueHeaders(
  dataPoints: IDataPoint[],
  key: 'yElement' | 'xElement',
  compare: 'xElement' | 'yElement'
): void {
  const duplicatesDict: { [key: string]: { [compare: string]: number } } = {};
  dataPoints.forEach((dataPoint, index) => {
    let uniqueHeader = dataPoint[key];
    const compareWith = dataPoint[compare];
    let foundIndex;
    while (
      (foundIndex = dataPoints.findIndex(
        // eslint-disable-next-line @typescript-eslint/no-loop-func
        (dt) => dt[key] === uniqueHeader && dt[compare] === compareWith
      )) < index &&
      foundIndex !== -1
    ) {
      if (duplicatesDict[uniqueHeader] === undefined) {
        duplicatesDict[uniqueHeader] = {};
      }
      if (duplicatesDict[uniqueHeader][compareWith] === undefined) {
        duplicatesDict[uniqueHeader][compareWith] = 0;
      }
      ++duplicatesDict[uniqueHeader][compareWith];
      uniqueHeader = `${dataPoint[key]} (${duplicatesDict[uniqueHeader][compareWith]})`;
    }
    dataPoint[key] = uniqueHeader;
  });
}

export function getXAxisReducedElements(
  elements: string[],
  dataPointContentWidth: number,
  minWidthPerXAxisElement: number
): string[] {
  const widthPerDataPoint = dataPointContentWidth / elements.length;

  if (widthPerDataPoint < minWidthPerXAxisElement) {
    const rest = elements.length % 2;
    // the latest one is the most important one. If even then the latest element is even, if not even then the latest element is not even
    // index starts with 0, that's why we use !== rest instead of === rest
    return elements.filter((_xElement, index) => index % 2 !== rest);
  }
  return elements;
}

export function calculateTooltipPosition(
  tooltipWidth: number,
  scrollbarWidth: number,
  x: number,
  y: number
): { top: number; left: number } {
  const offset = 5; // tooltip should not exactly appear next to the cursor
  const endOfWidth = (x + tooltipWidth) * window.devicePixelRatio + scrollbarWidth + offset;
  let left = x;
  if (endOfWidth > window.outerWidth) {
    left -= tooltipWidth - offset;
  } else {
    left += offset;
  }
  return {
    top: y + offset,
    left,
  };
}

/**
 * Returns a subset of the given identifiers that are available in the given dataPoints
 * @param identifiers
 * @param groupedData
 * @private
 */
export function getAvailableIdentifiers(identifiers: string[], groupedData: GroupedDataPoints): string[] {
  return identifiers.filter((identifier) => !!findDataPointThroughIdentifier(identifier, groupedData));
}

export function findDataPointThroughIdentifier(
  identifier: string,
  groupedData: GroupedDataPoints
): IDataPoint | undefined {
  for (const key of Object.keys(groupedData)) {
    const dataPoint = groupedData[key].find((dt) => dt.identifier === identifier);
    if (dataPoint) {
      return dataPoint;
    }
  }
  return undefined;
}

export function getHiddenYElements(yElements: string[], limitYElementCount: number): string[] {
  return yElements.slice(0, yElements.length - limitYElementCount);
}

export function getDataPointElement(x: number, y: number): SVGRectElement | undefined {
  const element = document.elementFromPoint(x, y);

  if (!element || !(element instanceof SVGRectElement)) {
    return undefined;
  }
  return element;
}

export function getYAxisElements(data: GroupedDataPoints): string[] {
  return Object.keys(data).reverse();
}

export function getAxisElements(
  data: GroupedDataPoints,
  limitYElementCount: number
): { yElements: string[]; xElements: string[]; showMoreVisible: boolean } {
  let yElements = getYAxisElements(data);
  const showMoreVisible = yElements.length > limitYElementCount;
  yElements = getLimitedYElements(yElements, limitYElementCount);
  const allXElements = [...yElements].reverse().reduce((xElements: string[], yElement: string) => {
    return [...xElements, ...data[yElement].map((dataPoint) => dataPoint.xElement)];
  }, []);
  return {
    yElements,
    xElements: Array.from(new Set(allXElements)),
    showMoreVisible,
  };
}

export function isScrollbarVisible(): boolean {
  // we only care about the scrollbar of the window (root element) not of a component
  const element = document.querySelector('body')?.firstElementChild;
  return !!element && element.scrollHeight > element.clientHeight;
}
