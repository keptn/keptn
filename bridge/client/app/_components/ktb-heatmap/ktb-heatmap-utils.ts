import { GroupedDataPoints, IDataPoint } from '../../_interfaces/heatmap';

export function createGroupedDataPoints(data: IDataPoint[]): GroupedDataPoints {
  setUniqueHeaders(data);
  return data.reduce((groupedData: GroupedDataPoints, dataPoint) => {
    groupedData[dataPoint.yElement] ||= [];
    groupedData[dataPoint.yElement].push(dataPoint);
    return groupedData;
  }, {});
}

function setUniqueHeaders(data: IDataPoint[]): void {
  setUniqueYAxis(data);
  setUniqueXAxis(data);
}

export const getLimitedYElements = (yElements: string[], limit: number): string[] => {
  if (yElements.length <= limit) {
    return yElements;
  }
  return yElements.slice(yElements.length - limit, yElements.length);
};

export function setUniqueYAxis(dataPoints: IDataPoint[]): void {
  // if it is the y-axis compare if there is the same SLI name with the same identifier
  const duplicatesDict: { [identifier: string]: { [compare: string]: number } } = {};
  const hasDuplicates: { [yElement: string]: boolean | undefined } = {};
  dataPoints.forEach((dataPoint, index) => {
    let displayValue = dataPoint.yElement;
    const identifier = dataPoint.identifier;
    let foundIndex;
    const hasYElementDuplicatesWithinIdentifierBeforeChild =
      (foundIndex = dataPoints.findIndex((dt) => dt.yElement === displayValue && dt.identifier === identifier)) <
        index && foundIndex !== -1;
    if (hasYElementDuplicatesWithinIdentifierBeforeChild) {
      duplicatesDict[identifier] ??= {};
      duplicatesDict[identifier][displayValue] ??= 1;
      ++duplicatesDict[identifier][displayValue];
      hasDuplicates[displayValue] = true;
      // occurrence is 1 but since it is the next element after the first found duplicate it's +1
      displayValue = `${dataPoint.yElement} (${duplicatesDict[identifier][displayValue]})`;
    }
    dataPoint.yElement = displayValue;
  });

  // if there are duplicates set the first occurrence to (1)
  dataPoints.forEach((dataPoint) => {
    const uniqueHeader = dataPoint.yElement;

    // if there are duplicates the counter is at least set to 2
    if (hasDuplicates[uniqueHeader]) {
      dataPoint.yElement = `${dataPoint.yElement} (1)`;
    }
  });
}

export function setUniqueXAxis(dataPoints: IDataPoint[]): void {
  const dataPointToIdentifierXElement = (
    dataPointDict: Record<string, string>,
    dataPoint: IDataPoint
  ): Record<string, string> => {
    dataPointDict[dataPoint.identifier] ??= dataPoint.xElement;
    return dataPointDict;
  };
  const identifierToDuplicateDict =
    (identifiers: string[], identifierToXElement: Record<string, string>) =>
    (duplicateCount: Record<string, number>, identifier: string, index: number): Record<string, number> => {
      const previousDuplicates = identifiers.filter(
        (nextIdentifier, nextIndex) =>
          nextIndex < index && identifierToXElement[nextIdentifier] === identifierToXElement[identifier]
      ).length;
      const nextDuplicates = identifiers.filter(
        (nextIdentifier, nextIndex) =>
          nextIndex > index && identifierToXElement[nextIdentifier] === identifierToXElement[identifier]
      ).length;
      duplicateCount[identifier] = previousDuplicates + nextDuplicates === 0 ? 0 : previousDuplicates + 1;
      return duplicateCount;
    };

  const identifierToXElement = dataPoints.reduce<Record<string, string>>(dataPointToIdentifierXElement, {});
  const identifiers = Object.keys(identifierToXElement);
  const duplicates = identifiers.reduce<Record<string, number>>(
    identifierToDuplicateDict(identifiers, identifierToXElement),
    {}
  );

  dataPoints.forEach((dataPoint) => {
    if (duplicates[dataPoint.identifier]) {
      dataPoint.xElement = `${dataPoint.xElement} (${duplicates[dataPoint.identifier]})`;
    }
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
