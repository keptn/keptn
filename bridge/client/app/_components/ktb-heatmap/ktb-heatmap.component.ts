/* eslint-disable @typescript-eslint/no-this-alias */
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  Output,
  ViewChild,
} from '@angular/core';
import * as d3 from 'd3';
import { BaseType, ScaleBand, Selection, ValueFn } from 'd3';
import { ResultTypes } from '../../../../shared/models/result-types';
import { DtButton } from '@dynatrace/barista-components/button';
import { v4 as uuid } from 'uuid';
import { KtbHeatmapTooltipComponent } from '../ktb-heatmap-tooltip/ktb-heatmap-tooltip.component';
import {
  EvaluationResultType,
  EvaluationResultTypeExtension,
  IDataPoint,
  IHeatmapTooltipType,
} from '../../_interfaces/heatmap';

type SVGGSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type HighlightSelection = Selection<SVGRectElement, unknown, HTMLElement, unknown>;
type SecondaryHighlightSelections = Selection<SVGRectElement, unknown, SVGGElement, unknown>;
type HeatmapTiles = Selection<SVGRectElement | null, IDataPoint, SVGGElement, unknown>;
type GroupedDataPoints = { [sli: string]: IDataPoint[] };

@Component({
  selector: 'ktb-heatmap',
  templateUrl: './ktb-heatmap.component.html',
  styleUrls: ['./ktb-heatmap.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHeatmapComponent implements OnDestroy, AfterViewInit {
  public readonly uniqueId = `heatmap-${uuid()}`;
  private readonly chartSelector = `#${this.uniqueId}`;
  private readonly heatmapSelector = `${this.chartSelector} .heatmap-container`;
  private readonly svgSelector = `${this.chartSelector}>svg`;
  private readonly firstSliPadding = 6; // "score" will then be 6px smaller than the rest.
  private readonly yAxisLabelWidth = 150;
  private readonly xAxisLabelWidth = 150;
  private readonly heightPerSli = 30;
  private readonly legendHeight = 50;
  private readonly limitSliCount = 10;
  private readonly legendItems: EvaluationResultType[] = [
    ResultTypes.PASSED,
    ResultTypes.WARNING,
    ResultTypes.FAILED,
    EvaluationResultTypeExtension.INFO,
  ];
  private readonly legendDisabledStatus: { [category: string]: boolean } = {
    [ResultTypes.PASSED]: false,
    [ResultTypes.WARNING]: false,
    [ResultTypes.FAILED]: false,
    [EvaluationResultTypeExtension.INFO]: false,
  };
  private readonly showMoreButtonPadding = 6;
  private readonly showMoreButtonHeight = 32 + this.showMoreButtonPadding;
  private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  private readonly minWidthPerXAxisElement = 25;
  private xAxis?: ScaleBand<string>;
  private yAxis?: ScaleBand<string>;
  private dataPointContentWidth = 0;
  private height = 0;
  private _selectedDataPoint?: IDataPoint;
  private mouseCoordinates = { x: 0, y: 0 };
  private groupedData: GroupedDataPoints = {};
  public showMoreVisible = true;
  public showMoreExpanded = false;

  @ViewChild('showMoreButton', { static: false }) showMoreButton!: DtButton;
  @ViewChild('tooltip', { static: false }) tooltip!: KtbHeatmapTooltipComponent;
  @Output() selectedDataPointChange = new EventEmitter<IDataPoint>();

  //TODO:
  // - Remove testing data afterwards

  @Input()
  set dataPoints(data: IDataPoint[]) {
    this.removeHeatmap();
    this.setUniqueHeaders(data, 'date', 'sli');
    this.setUniqueHeaders(data, 'sli', 'date');
    this.groupedData = data.reduce((groupedData: GroupedDataPoints, dataPoint) => {
      groupedData[dataPoint.sli] ||= [];
      groupedData[dataPoint.sli].push(dataPoint);
      return groupedData;
    }, {});
    this.createHeatmap(this.groupedData);
    this.onResize(); // generating the heatmap may introduce a scrollbar
    this.click(this.selectedDataPoint, true); // restore previously selected dataPoint
  }

  @Input()
  public set selectedDataPoint(dataPoint: IDataPoint | undefined) {
    this.click(dataPoint, true);
  }
  public get selectedDataPoint(): IDataPoint | undefined {
    return this._selectedDataPoint;
  }

  private get heatmapInstance(): SVGGSelection {
    return d3.select(this.heatmapSelector);
  }

  private get dataPointContainer(): SVGGSelection {
    return this.heatmapInstance.select('.data-point-container');
  }

  private get yAxisContainer(): SVGGSelection {
    return this.heatmapInstance.select('.y-axis-container');
  }

  private get xAxisContainer(): SVGGSelection {
    return this.heatmapInstance.select('.x-axis-container');
  }

  private get legendContainer(): SVGGSelection {
    return this.heatmapInstance.select('.legend-container');
  }

  private get dataPointElements(): HeatmapTiles {
    return this.dataPointContainer.selectAll('.data-point');
  }

  private get highlight(): HighlightSelection {
    return this.heatmapInstance.select('.highlight-primary');
  }

  private get secondaryHighlights(): SecondaryHighlightSelections {
    return this.heatmapInstance.selectAll('.highlight-secondary');
  }

  private get dataPointContainerRect(): SVGGSelection {
    return this.heatmapInstance.select('.data-point-container-rect');
  }

  constructor(private elementRef: ElementRef, private _changeDetectorRef: ChangeDetectorRef) {
    // has to be globally instead of component bound, else scrolling into it will not have any mouse coordinates
    this.mouseMoveListener = (event: MouseEvent): void => this.onMouseMove(event);
    document.addEventListener('mousemove', this.mouseMoveListener);
  }

  private onMouseMove(event: MouseEvent): void {
    this.mouseCoordinates = {
      x: event.x * window.devicePixelRatio, // coordinates may stay and zoom-level could change. Normalize the coordinates.
      y: event.y * window.devicePixelRatio,
    };
  }

  @HostListener('window:resize', ['$event'])
  public onResize(): void {
    const { width, height } = this.setAndGetAvailableSpace();

    this.resizeSvg(width, height);
    this.resizeXAxis();
    this.resizeDataPoints();
    this.resizeHighlights();
    this.resizeShowMoreButton();
    this.resizeLegend();
    this.resizeDataPointContainerRect();
    this.onScroll(); // on zoom, the tooltip has to be adjusted
  }

  @HostListener('window:scroll')
  private onScroll(): void {
    const x = this.mouseCoordinates.x / window.devicePixelRatio;
    const y = this.mouseCoordinates.y / window.devicePixelRatio;
    const element = this.getDataPointElement(x, y);
    const dt = this.getDataPointThroughCoordinates(x, y);

    if (!element || !dt) {
      this.setTooltipVisibility(false);
      return;
    }

    const mouseEvent = new MouseEvent('move', {
      clientY: y,
      clientX: x,
    });
    this.mouseOver(element);
    this.mouseMove(mouseEvent, dt);
  }

  private getDataPointElement(x: number, y: number): SVGRectElement | undefined {
    const element = document.elementFromPoint(x, y);

    if (!element || !(element instanceof SVGRectElement)) {
      return undefined;
    }
    return element;
  }

  private getDataPointThroughCoordinates(x: number, y: number): IDataPoint | undefined {
    const element = this.getDataPointElement(x, y);

    if (!element || !this.heatmapInstance.node()?.contains(element)) {
      return;
    }

    return d3.select(element)?.datum() as IDataPoint | undefined;
  }

  private removeHeatmap(): void {
    d3.select(this.svgSelector).remove();
  }

  private setTooltipVisibility(visible: boolean): void {
    const element: HTMLElement = this.tooltip._elementRef.nativeElement;
    const classList = element.classList;
    if (visible) {
      return classList.remove('hidden');
    }
    classList.add('hidden');
  }

  private createHeatmap(data: GroupedDataPoints): void {
    const { slis, dates } = this.getAxisElements(data);

    this.setHeight(slis.length);
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = d3.select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    svg.append('g').classed('heatmap-container', true).attr('transform', `translate(${this.yAxisLabelWidth}, 0)`);

    this.setData(data, dates, slis);
    this.createLegend();
    this.resizeDataPointContainerRect();
  }

  private getAxisElements(data: GroupedDataPoints): { slis: string[]; dates: string[] } {
    let slis = Object.keys(data);
    this.showMoreVisible = slis.length > this.limitSliCount;
    if (this.showMoreVisible) {
      slis = this.getLimitedSLIs(slis);
    }
    const allDates = slis.reduce((dates: string[], sli: string) => {
      return [...dates, ...data[sli].map((dataPoint) => dataPoint.date)];
    }, []);
    return {
      slis,
      dates: Array.from(new Set(allDates)),
    };
  }

  private setAndGetAvailableSpace(): { height: number; width: number } {
    const parentElement: HTMLElement = this.elementRef.nativeElement.parentNode;
    const availableSpace = parentElement.getBoundingClientRect();
    const width = availableSpace.width * window.devicePixelRatio; // adjust to zoom-level
    const height =
      this.height + this.xAxisLabelWidth + this.legendHeight + (this.showMoreVisible ? this.showMoreButtonHeight : 0);
    this.dataPointContentWidth = width - this.yAxisLabelWidth;

    return {
      height,
      width,
    };
  }

  private resizeSvg(width: number, height: number): void {
    d3.select(this.svgSelector).attr('viewBox', `0 0 ${width} ${height}`).attr('width', width).attr('height', height);
  }

  private resizeXAxis(): void {
    if (!this.xAxis) {
      return;
    }
    this.xAxis = this.xAxis.range([0, this.dataPointContentWidth]);
    const xAxisContainer = this.xAxisContainer;

    this.setXAxisCoordinates(xAxisContainer);
    this.attachXAxis(xAxisContainer, this.xAxis);
  }

  private resizeDataPoints(): void {
    if (this.xAxis && this.yAxis) {
      this.setDataPointCoordinates(this.dataPointElements, this.xAxis, this.yAxis);
    }
  }

  private resizeHighlights(): void {
    if (!this.selectedDataPoint) {
      return;
    }

    this.setHighlightCoordinates(this.selectedDataPoint.date);
    this.setSecondaryHighlightCoordinates(this.selectedDataPoint.comparedIdentifier);
  }

  private resizeShowMoreButton(): void {
    if (this.showMoreVisible) {
      const htmlElement: HTMLElement = this.showMoreButton._elementRef.nativeElement;

      htmlElement.style.top = `${this.height + this.showMoreButtonPadding / 2}px`;
      htmlElement.style.left = `${this.yAxisLabelWidth}px`;
      htmlElement.style.width = `${this.dataPointContentWidth}px`;
    }
  }

  private resizeLegend(): void {
    const legend = this.legendContainer;
    const fullLength = legend.node()?.getBoundingClientRect().width ?? 0;
    const centerXPosition = (this.dataPointContentWidth - fullLength) / 2;
    const yPosition = this.height + this.xAxisLabelWidth + 10 + (this.showMoreVisible ? this.showMoreButtonHeight : 0);
    legend.attr('transform', `translate(${centerXPosition}, ${yPosition})`);
  }

  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'date', compare: 'sli'): void;
  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'sli', compare: 'date'): void;
  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'sli' | 'date', compare: 'date' | 'sli'): void {
    // TODO: change the first found duplicate to (1)?
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

  private setData(data: GroupedDataPoints, xAxisElements: string[], yAxisElements: string[]): void {
    const heatmap = this.heatmapInstance;
    this.xAxis = this.addXAxis(heatmap, xAxisElements);
    this.yAxis = this.addYAxis(heatmap, yAxisElements);
    this.generateHeatmapTiles(data, this.xAxis, this.yAxis);
  }

  private addXAxis(heatmap: SVGGSelection, dates: string[]): ScaleBand<string> {
    const x = d3.scaleBand().range([0, this.dataPointContentWidth]).domain(dates);
    const xAxisContainer = heatmap.append('g').attr('class', 'x-axis-container');

    this.setXAxisCoordinates(xAxisContainer);
    this.attachXAxis(xAxisContainer, x);
    return x;
  }

  private setXAxisCoordinates(xAxisContainer: SVGGSelection): void {
    xAxisContainer.attr(
      'transform',
      `translate(0, ${this.height + (this.showMoreVisible ? this.showMoreButtonHeight : 0)})`
    );
  }

  /**
   * Attaches the xAxis to the heatmap and sets the correct style.
   * If there are too many elements, every second element is shown in the axis
   * @param xAxisContainer
   * @param x
   * @private
   */
  private attachXAxis(xAxisContainer: SVGGSelection, x: ScaleBand<string>): void {
    xAxisContainer
      .call(d3.axisBottom(x).tickSize(5).tickValues(this.getXAxisReducedElements(x.domain())))
      .selectAll('text')
      .attr('class', 'x-axis-identifier')
      .attr('dx', '-.8em')
      .attr('dy', '.15em');
  }

  private getXAxisReducedElements(elements: string[]): string[] {
    const widthPerDatePoint = this.dataPointContentWidth / elements.length;

    if (widthPerDatePoint < this.minWidthPerXAxisElement) {
      const rest = elements.length % 2;
      // the latest one is the most important one. If even then the latest element is even, if not even then the latest element is not even
      // index starts with 0, that's why we use !== rest instead of === rest
      return elements.filter((_date, index) => index % 2 !== rest);
    }
    return elements;
  }

  private addYAxis(heatmap: SVGGSelection, slis: string[]): ScaleBand<string> {
    const y = d3.scaleBand().range([this.height, 0]).domain(slis);
    const yAxisContainer = heatmap.append('g').attr('class', 'y-axis-container');
    this.attachYAxis(yAxisContainer, y);
    return y;
  }

  /**
   * Attaches the yAxis to the heatmap, sets the correct style and adds ellipsis style to the text if needed
   * @param yAxisContainer
   * @param y
   * @private
   */
  private attachYAxis(yAxisContainer: SVGGSelection, y: ScaleBand<string>): void {
    const yAxis = yAxisContainer.call(d3.axisLeft(y).tickSize(0));
    yAxis.selectAll('.tick').each(this.setEllipsisStyle(this.yAxisLabelWidth));
    yAxis.select('.domain').remove();
  }

  private setEllipsisStyle(labelWidth: number): ValueFn<BaseType, unknown, void> {
    return function (this: BaseType): void {
      const self = d3.select(this as SVGGElement);
      if (self.empty()) {
        return;
      }
      const textElement = self.select('text');
      const originalText = textElement.text();
      let textLength = self.node()?.getBoundingClientRect().width ?? 0;
      let text = originalText;

      while (textLength > labelWidth && text.length > 0) {
        text = text.slice(0, -1);
        textElement.text(text + '...');
        textLength = self.node()?.getBoundingClientRect().width ?? 0;
      }
      self.append('title').text(originalText);
    };
  }

  private mouseOver(element: SVGGElement): void {
    // don't show tooltip if dataPoint is disabled
    if (element.classList.contains('disabled')) {
      return;
    }
    this.setTooltipVisibility(true);
  }

  private mouseLeave(): void {
    this.setTooltipVisibility(false);
  }

  private mouseMove(event: MouseEvent, dataPoint: IDataPoint): void {
    this.tooltip.tooltip = dataPoint.tooltip;
    const htmlElement: HTMLElement = this.tooltip._elementRef.nativeElement;
    const tooltipWidth = htmlElement.getBoundingClientRect().width;
    const scrollbarWidth = this.isScrollbarVisible() ? 18 : 0; // just assume a default scrollbar-width of 18px
    const { top, left } = this.calculateTooltipPosition(tooltipWidth, scrollbarWidth, event.x, event.y);

    htmlElement.style.top = `${top}px`;
    htmlElement.style.left = `${left}px`;
  }

  private calculateTooltipPosition(
    tooltipWidth: number,
    scrollbarWidth: number,
    x: number,
    y: number
  ): { top: number; left: number } {
    const offset = 5; // tooltip should not exactly appear next to the cursor
    const endOfWidth = (x + tooltipWidth) * window.devicePixelRatio + scrollbarWidth + offset;
    let left;
    if (endOfWidth > window.outerWidth) {
      left = x - tooltipWidth - offset;
    } else {
      left = x + offset;
    }
    return {
      top: y + offset,
      left,
    };
  }

  private isScrollbarVisible(): boolean {
    // we only care about the scrollbar of the window (root element) not of a component
    const element = document.querySelector('body')?.firstElementChild;
    if (!element) {
      return false;
    }

    return element.scrollHeight > element.clientHeight;
  }

  private removeHighlights(): void {
    this.highlight.remove();
    this.secondaryHighlights.remove();
  }

  /**
   * Selects the given dataPoint and sets primary and secondary highlights accordingly.
   * @param dataPoint the dataPoint that should be selected
   * @param preSelectDataPoint if true the selected dataPoint is not emitted
   *  and there is a check beforehand if the dataPoint exists
   * @private
   */
  private click(dataPoint?: IDataPoint, preSelectDataPoint = false): void {
    this.removeHighlights();
    const heatmap = this.heatmapInstance;

    if (!this.xAxis || !dataPoint) {
      this._selectedDataPoint = undefined;
      return;
    }
    this._selectedDataPoint = dataPoint;

    if (preSelectDataPoint && !this.findDateThroughIdentifier(dataPoint.identifier)) {
      this._selectedDataPoint = undefined;
      return;
    }

    heatmap.append('rect').attr('class', 'highlight-primary');
    this.setHighlightCoordinates(dataPoint.date);

    const foundIdentifiers = this.getAvailableIdentifiers(dataPoint.comparedIdentifier);
    heatmap.selectAll().data(foundIdentifiers).join('rect').attr('class', 'highlight-secondary');
    this.setSecondaryHighlightCoordinates(foundIdentifiers);

    if (!preSelectDataPoint) {
      this.selectedDataPointChange.emit(dataPoint);
    }
  }

  /**
   * For the special case that the user clicks on an SLI dataPoint that does not exist (another SLI dataPoint in the column exists)
   * @param event$
   * @param element
   */
  public contentClick(event$: MouseEvent, element: SVGRectElement): void {
    const containerY = element.getBoundingClientRect().top;
    const dataPoint = this.getDataPointThroughCoordinates(event$.x, containerY + 5); // offset to make sure to click on the tile
    if (!dataPoint) {
      return;
    }
    this.click(dataPoint);
  }

  /**
   * Returns a subset of the given identifiers that are available in the given dataPoints
   * @param identifiers
   * @private
   */
  private getAvailableIdentifiers(identifiers: string[]): string[] {
    return identifiers.filter((identifier) => !!this.findDateThroughIdentifier(identifier));
  }

  private getHighlightWidth(): number {
    if (!this.xAxis) {
      return 0;
    }
    const xAxisElements = this.xAxis.domain();
    return this.dataPointContentWidth / xAxisElements.length;
  }

  private setHighlightCoordinates(identifier: string): void {
    if (!this.xAxis) {
      return;
    }

    this.highlight
      .attr('x', this.xAxis(identifier) ?? null)
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private setSecondaryHighlightCoordinates(identifiers: string[]): void {
    if (!this.xAxis) {
      return;
    }
    const xAxis = this.xAxis;

    this.secondaryHighlights
      .attr('x', (_dt, index) => {
        const date = this.findDateThroughIdentifier(identifiers[index]);
        if (!date) {
          return null;
        }
        return xAxis(date) ?? null;
      })
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private findDateThroughIdentifier(identifier: string): string | undefined {
    for (const key of Object.keys(this.groupedData)) {
      const dataPoint = this.groupedData[key].find((dt) => dt.identifier === identifier);
      if (dataPoint) {
        return dataPoint.date;
      }
    }
    return undefined;
  }

  private generateHeatmapTiles(
    data: GroupedDataPoints,
    x: ScaleBand<string>,
    y: ScaleBand<string>,
    yAxisElements = y.domain()
  ): void {
    const _this = this;
    let container = this.dataPointContainer;

    if (container.empty()) {
      container = this.heatmapInstance.append('g').classed('data-point-container', true);
      container
        .append('rect')
        .classed('data-point-container-rect', true)
        .on('click', function (this: SVGRectElement, event: PointerEvent) {
          _this.contentClick(event, this);
        });
    }
    const dataPoints = container
      .selectAll()
      .data(yAxisElements)
      .enter()
      .append('g')
      .attr('uitestid', (sli) => sli.replace(/ /g, '-')) // TODO: do we need this?
      .selectAll()
      .data((key) => data[key])
      .join('rect')
      .attr('class', (dataPoint) => dataPoint.color)
      .classed('data-point', true)
      // set all new dataPoints (show all SLIs) to disabled if needed
      .classed('disabled', (dataPoint: IDataPoint) => this.legendDisabledStatus[dataPoint.color])
      .attr('uitestid', (dataPoint) => `ktb-heatmap-tile-${dataPoint.date.replace(/ /g, '-')}`) // TODO: do we need this?
      .on('click', (_event: PointerEvent, dataPoint: IDataPoint) => this.click(dataPoint))
      .on('mouseover', function (this: SVGGElement | null) {
        if (!this) {
          return;
        }
        _this.mouseOver(this);
      })
      .on('mousemove', (event: MouseEvent, dataPoint: IDataPoint) => this.mouseMove(event, dataPoint))
      .on('mouseleave', () => this.mouseLeave());

    this.setDataPointCoordinates(dataPoints, x, y);
  }

  private resizeDataPointContainerRect(): void {
    const coordinates = this.dataPointContainer.node()?.getBoundingClientRect();
    if (!coordinates) {
      return;
    }
    this.dataPointContainerRect.attr('width', this.dataPointContentWidth).attr('height', coordinates.height);
  }

  private setDataPointCoordinates(dataPoints: HeatmapTiles, x: ScaleBand<string>, y: ScaleBand<string>): void {
    const yAxisElements = y.domain();
    const firstSli = yAxisElements[yAxisElements.length - 1];
    dataPoints
      .attr('x', (dataPoint) => x(dataPoint.date) ?? null)
      .attr('y', (dataPoint) => {
        const yCoordinate = y(dataPoint.sli);
        if (yCoordinate !== undefined && dataPoint.sli === firstSli) {
          return yCoordinate + this.firstSliPadding / 2;
        }
        return yCoordinate ?? null;
      })
      .attr('width', x.bandwidth())
      .attr('height', (dataPoint) => {
        const height = y.bandwidth();
        if (dataPoint.sli === firstSli) {
          return height - this.firstSliPadding;
        }
        return height;
      });
  }

  private createLegend(): void {
    const spaceBetweenLegendItems = 30;
    const legend = this.heatmapInstance.append('g').attr('class', 'legend-container');
    let xCoordinate = 0;
    for (const category of this.legendItems) {
      const legendItem = legend
        .append('g')
        .classed('legend-item', true)
        .on('click', () => {
          this.legendDisabledStatus[category] = !this.legendDisabledStatus[category];
          const isDisabled = this.legendDisabledStatus[category];
          this.setLegendDisabled(legendItem, isDisabled);
          this.setDataPointsDisabled(category, isDisabled);
        });
      legendItem
        .append('circle')
        .attr('cx', xCoordinate)
        .attr('r', 6)
        .classed('legend-circle', true)
        .classed(category, true);
      xCoordinate += 10; // space between circle and text
      const text = legendItem.append('text').attr('x', xCoordinate).text(category).classed('legend-text', true);
      const textWidth = text.node()?.getComputedTextLength() ?? 0;
      xCoordinate += textWidth + spaceBetweenLegendItems;
    }
    this.resizeLegend();
  }

  private setLegendDisabled(legendItem: SVGGSelection, status: boolean): void {
    legendItem.select('circle').classed('disabled', status);
  }

  private setDataPointsDisabled(category: EvaluationResultType, isDisabled: boolean): void {
    this.dataPointElements.each(function (this: SVGGElement | null, dataPoint: IDataPoint) {
      if (this && dataPoint.color === category) {
        d3.select(this).classed('disabled', isDisabled);
      }
    });
  }

  public showMoreToggle(): void {
    this.showMoreExpanded = !this.showMoreExpanded;

    if (this.showMoreExpanded) {
      this.expandHeatmap();
    } else {
      this.collapseHeatmap();
    }
    this.onResize();
  }

  private expandHeatmap(): void {
    if (!this.xAxis || !this.yAxis) {
      return;
    }
    const slis = Object.keys(this.groupedData);
    this.setHeight(slis.length);
    this.updateYAxis(slis);

    this.generateHeatmapTiles(this.groupedData, this.xAxis, this.yAxis, this.getHiddenSLIs(slis));
  }

  private getHiddenSLIs(slis: string[]): string[] {
    return slis.slice(0, slis.length - this.limitSliCount);
  }

  private getLimitedSLIs(slis: string[]): string[] {
    return slis.slice(slis.length - this.limitSliCount, slis.length);
  }

  private updateYAxis(slis: string[]): void {
    if (!this.yAxis) {
      return;
    }
    this.yAxis = this.yAxis.range([this.height, 0]).domain(slis);
    this.attachYAxis(this.yAxisContainer, this.yAxis);
  }

  private collapseHeatmap(): void {
    this.setHeight(this.limitSliCount);
    this.dataPointContainer
      .selectAll('g')
      .filter((_element, index) => index >= this.limitSliCount)
      .remove();

    const slis = Object.keys(this.groupedData);
    this.updateYAxis(this.getLimitedSLIs(slis));
  }

  private setHeight(elementCount: number): void {
    this.height = elementCount * this.heightPerSli;
  }

  public ngOnDestroy(): void {
    document.removeEventListener('mousemove', this.mouseMoveListener);
  }

  //<editor-fold desc="test data">
  //TODO: remove
  public ngAfterViewInit(): void {
    this.dataPoints = this.generateTestData(12, 50);
    this.click(this.groupedData.score[1]);
  }

  private generateTestData(sliCounter: number, counter: number): IDataPoint[] {
    const categories = [];
    for (let i = 0; i < sliCounter; ++i) {
      categories.push(`response time p${i}`);
    }
    const data: IDataPoint[] = [];
    const dateMillis = new Date().getTime();

    // adding one duplicate (two evaluations have the same time)
    for (const category of [...categories, 'score']) {
      data.push({
        date: new Date(dateMillis).toISOString(),
        sli: category,
        color: this.getColor(Math.floor(Math.random() * 4)),
        tooltip: {
          type: IHeatmapTooltipType.SLI,
          value: Math.random(),
          keySli: Math.floor(Math.random() * 2) === 1,
          score: Math.floor(Math.random() * 100),
          passTargets: [
            {
              targetValue: 0,
              criteria: '<=1',
              violated: true,
            },
          ],
          warningTargets: [
            {
              targetValue: 0,
              violated: false,
              criteria: '<=10',
            },
          ],
        },
        identifier: `keptnContext_${-1}`,
        comparedIdentifier: [],
      });
    }

    // fill SLIs with random data (-1 to have an evaluation with "missing" data)
    for (const category of categories) {
      for (let i = 0; i < counter - 1; ++i) {
        data.push({
          date: new Date(dateMillis + i).toISOString(),
          sli: category,
          color: this.getColor(Math.floor(Math.random() * 4)),
          tooltip: {
            type: IHeatmapTooltipType.SLI,
            value: Math.random(),
            keySli: Math.floor(Math.random() * 2) === 1,
            score: Math.floor(Math.random() * 100),
            passTargets: [
              {
                targetValue: 0,
                criteria: '<=1',
                violated: true,
              },
            ],
            warningTargets: [
              {
                targetValue: 0,
                violated: false,
                criteria: '<=10',
              },
            ],
          },
          identifier: `keptnContext_${i}`,
          comparedIdentifier: [`keptnContext_${i - 1}`, `keptnContext_${i - 2}`],
        });
      }
    }
    categories.push('score');
    for (let i = 0; i < counter; ++i) {
      data.push({
        date: new Date(dateMillis + i).toISOString(),
        sli: 'score',
        color: this.getColor(Math.floor(Math.random() * 4)),
        tooltip: {
          type: IHeatmapTooltipType.SCORE,
          value: Math.random(),
          fail: Math.floor(Math.random() * 2) === 1,
          failedCount: Math.random(),
          warn: Math.floor(Math.random() * 2) === 1,
          passCount: Math.random(),
          thresholdPass: Math.random(),
          thresholdWarn: Math.random(),
          warningCount: Math.random(),
        },
        identifier: `keptnContext_${i}`,
        comparedIdentifier: [`keptnContext_${i - 1}`, `keptnContext_${i - 2}`],
      });
    }
    return data;
  }

  private getColor(value: number): EvaluationResultType {
    if (value === 0) {
      return ResultTypes.FAILED;
    }
    if (value === 1) {
      return ResultTypes.WARNING;
    }
    if (value === 2) {
      return ResultTypes.PASSED;
    }
    return EvaluationResultTypeExtension.INFO;
  }

  //</editor-fold>
}
/* eslint-enable @typescript-eslint/no-this-alias */
