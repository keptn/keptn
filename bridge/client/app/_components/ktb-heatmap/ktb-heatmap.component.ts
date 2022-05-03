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
  private readonly showMoreButtonPadding = 6;
  private readonly showMoreButtonHeight = 32 + this.showMoreButtonPadding;
  private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  private xAxis?: ScaleBand<string>;
  private yAxis?: ScaleBand<string>;
  private dataPointContentWidth = 0;
  private height = 0;
  private highlight?: Selection<SVGRectElement, unknown, HTMLElement, unknown>;
  private secondaryHighlights: Selection<SVGRectElement, unknown, HTMLElement, unknown>[] = [];
  private _selectedDataPoint?: IDataPoint;
  private mouseCoordinates = { x: 0, y: 0 };
  private groupedData: GroupedDataPoints = {};
  public showMoreVisible = true;
  public showMoreExpanded = false;

  @ViewChild('showMoreButton', { static: false }) showMoreButton!: DtButton;
  @ViewChild('tooltip', { static: false }) tooltip!: KtbHeatmapTooltipComponent;
  @Output() selectedDataPointChange = new EventEmitter<IDataPoint>();
  // unsure about:
  // should tileSelected emit the datapoint or just the identifier?
  // Re-positioning of tooltip only on hover-item-change?
  // Should the secondaryHighlight be an index or the identifier of the dataPoint?

  // TODO:
  //  - Only show every xth date if there are too many dataPoints?
  //  - Remove testing data afterwards

  @Input()
  set dataPoints(data: IDataPoint[]) {
    // TODO: remove previous heatmap
    //  What to do with selected datapoint?
    //  Should the heatmap group the dataPoints or the component that provides them?
    this.setUniqueHeaders(data, 'date', 'sli');
    this.setUniqueHeaders(data, 'sli', 'date');
    this.groupedData = data.reduce((groupedData: GroupedDataPoints, dataPoint) => {
      groupedData[dataPoint.sli] = groupedData[dataPoint.sli] || [];
      groupedData[dataPoint.sli].push(dataPoint);
      return groupedData;
    }, {});
    this.createHeatmap(this.groupedData);
    this.onResize(); // generating the heatmap may introduce a scrollbar
  }

  @Input()
  set selectedDataPoint(dataPoint: IDataPoint | undefined) {
    this.click(dataPoint);
  }
  get selectedDataPoint(): IDataPoint | undefined {
    return this._selectedDataPoint;
  }

  private get dataPointContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.data-point-container');
  }

  private get yAxisContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.y-axis-container');
  }

  private get xAxisContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.x-axis-container');
  }

  private get legendContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.legend-container');
  }

  private get dataPointElements(): HeatmapTiles {
    return this.dataPointContainer.selectAll('.data-point');
  }

  constructor(private elementRef: ElementRef, private _changeDetectorRef: ChangeDetectorRef) {
    // has to be globally instead of component bound, else scrolling into it will not have any mouse coordinates
    this.mouseMoveListener = (event: MouseEvent): void => this.onMouseMove(event);
    document.addEventListener('mousemove', this.mouseMoveListener);
  }

  public ngAfterViewInit(): void {
    this.dataPoints = this.generateTestData(12, 50); // TODO: remove testing data afterwards
    this.click(this.groupedData.score[1]);
    this.resizeShowMoreButton();
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
    this.onScroll(); // on zoom, the tooltip has to be adjusted
  }

  @HostListener('window:scroll')
  private onScroll(): void {
    const x = this.mouseCoordinates.x / window.devicePixelRatio;
    const y = this.mouseCoordinates.y / window.devicePixelRatio;
    const element = document.elementFromPoint(x, y);

    if (!element || !(element instanceof SVGRectElement)) {
      this.setTooltipVisibility(false);
      return;
    }

    const isDataPointInHeatmapInstance = (d3.select(this.chartSelector).node() as HTMLElement | undefined)?.contains(
      element
    );
    if (!isDataPointInHeatmapInstance) {
      this.setTooltipVisibility(false);
      return;
    }

    const dt = d3.select(element)?.datum() as IDataPoint | undefined;

    if (!dt) {
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

  private setTooltipVisibility(visible: boolean): void {
    const element: HTMLElement = this.tooltip._elementRef.nativeElement;
    const classList = element.classList;
    if (visible) {
      return classList.remove('hidden');
    }
    classList.add('hidden');
  }

  private createHeatmap(data: GroupedDataPoints): void {
    let slis = Object.keys(data);
    this.showMoreVisible = slis.length > this.limitSliCount;
    if (this.showMoreVisible) {
      slis = this.getLimitedSLIs(slis);
    }
    const allDates = slis.reduce((dates: string[], sli: string) => {
      return [...dates, ...data[sli].map((dataPoint) => dataPoint.date)];
    }, []);
    const distinctDates = Array.from(new Set(allDates));

    this.setHeight(slis.length);
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = d3.select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    svg.append('g').classed('heatmap-container', true).attr('transform', `translate(${this.yAxisLabelWidth}, 0)`);

    this.setData(data, distinctDates, slis, this.showMoreVisible);
    this.createLegend();
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
    d3.select(`${this.chartSelector}>svg`)
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('width', width)
      .attr('height', height);
  }

  private resizeXAxis(): void {
    if (!this.xAxis) {
      return;
    }
    this.xAxis = this.xAxis.range([0, this.dataPointContentWidth]);
    this.xAxisContainer
      .call(d3.axisBottom(this.xAxis))
      .attr('transform', `translate(0, ${this.height + (this.showMoreVisible ? this.showMoreButtonHeight : 0)})`);
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
    if (this.highlight) {
      this.setHighlightCoordinates(this.highlight, this.selectedDataPoint.date);
    }
    for (let i = 0; i < this.secondaryHighlights.length; ++i) {
      this.setSecondaryHighlightCoordinates(this.secondaryHighlights[i], this.selectedDataPoint.comparedIndices[i]);
    }
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

  private generateTestData(sliCounter: number, counter: number): IDataPoint[] {
    const categories = [];
    for (let i = 0; i < sliCounter; ++i) {
      categories.push(`response time p${i}`);
    }
    const data: IDataPoint[] = [];
    const dateMillis = new Date().getTime();

    // adding one duplicate (two evaluations have the same time)
    for (const category of categories) {
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
        comparedIndices: [],
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
          comparedIndices: [],
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
        comparedIndices: [],
      });
    }

    data[5].comparedIndices = [0, 1, 2]; //TODO: has to be set for all SLIs
    data[6].comparedIndices = [0, 1];
    data[7].comparedIndices = [5, 6];
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

  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'date', compare: 'sli'): void;
  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'sli', compare: 'date'): void;
  private setUniqueHeaders(dataPoints: IDataPoint[], key: 'sli' | 'date', compare: 'date' | 'sli'): void {
    // TODO: change the first found duplicate to (1)?
    dataPoints.forEach((dataPoint, index) => {
      let duplicates = 0;
      let uniqueHeader = dataPoint[key];
      let foundIndex;
      while (
        (foundIndex = dataPoints.findIndex(
          // eslint-disable-next-line @typescript-eslint/no-loop-func
          (dt) => dt[key] === uniqueHeader && dt[compare] === dataPoint[compare]
        )) < index &&
        foundIndex !== -1
      ) {
        ++duplicates;
        uniqueHeader = `${dataPoint[key]} (${duplicates})`;
      }
      dataPoint[key] = uniqueHeader;
    });
  }

  private setData(
    data: GroupedDataPoints,
    xAxisElements: string[],
    yAxisElements: string[],
    showMoreVisible: boolean
  ): void {
    const heatmap: SVGGSelection = d3.select(this.heatmapSelector);
    this.xAxis = this.addXAxis(heatmap, xAxisElements, showMoreVisible);
    this.yAxis = this.addYAxis(heatmap, yAxisElements);
    this.generateHeatmapTiles(data, this.xAxis, this.yAxis);
  }

  private addXAxis(heatmap: SVGGSelection, dates: string[], showMoreVisible: boolean): ScaleBand<string> {
    const x = d3.scaleBand().range([0, this.dataPointContentWidth]).domain(dates);
    heatmap
      .append('g')
      .attr('class', 'x-axis-container')
      .attr('transform', `translate(0, ${this.height + (showMoreVisible ? this.showMoreButtonHeight : 0)})`)
      .call(d3.axisBottom(x).tickSize(5))
      .selectAll('text')
      .attr('class', 'x-axis-identifier')
      .attr('dx', '-.8em')
      .attr('dy', '.15em');
    return x;
  }

  private addYAxis(heatmap: SVGGSelection, slis: string[]): ScaleBand<string> {
    const y = d3.scaleBand().range([this.height, 0]).domain(slis);
    const yAxisContainer = heatmap.append('g').attr('class', 'y-axis-container');
    const yAxis = this.callYAxis(yAxisContainer, y);
    this.setYAxisStyle(yAxis);
    return y;
  }

  private callYAxis(yAxisContainer: SVGGSelection, y: ScaleBand<string>): SVGGSelection {
    return yAxisContainer.call(d3.axisLeft(y).tickSize(0));
  }

  private setYAxisStyle(yAxis: SVGGSelection): void {
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
    const scrollbarWidth = this.isScrollbarVisible() ? 18 : 0; // just take a default value
    const offset = 5; // tooltip should not exactly appear on the cursor

    const endOfWidth = (event.x + tooltipWidth) * window.devicePixelRatio + scrollbarWidth + offset;

    htmlElement.style.top = event.y + offset + 'px';

    if (endOfWidth > window.outerWidth) {
      htmlElement.style.left = event.x - tooltipWidth - offset + 'px';
    } else {
      htmlElement.style.left = event.x + offset + 'px';
    }
  }

  private isScrollbarVisible(): boolean {
    const element = document.querySelector('body')?.firstElementChild;
    if (!element) {
      return false;
    }

    return element.scrollHeight > element.clientHeight;
  }

  private removeHighlights(): void {
    this.highlight?.remove();
    for (const highlight of this.secondaryHighlights) {
      highlight.remove();
    }
  }

  private click(dataPoint?: IDataPoint): void {
    this.removeHighlights();
    const heatmap: SVGGSelection = d3.select(this.heatmapSelector);

    if (!this.xAxis || !dataPoint) {
      this._selectedDataPoint = undefined;
      return;
    }
    this._selectedDataPoint = dataPoint;

    this.highlight = heatmap.append('rect').attr('class', 'highlight-primary');
    this.setHighlightCoordinates(this.highlight, dataPoint.date);

    this.secondaryHighlights = dataPoint.comparedIndices.map((index) => {
      const secondaryHighlight = heatmap.append('rect').attr('class', 'highlight-secondary');
      this.setSecondaryHighlightCoordinates(secondaryHighlight, index);
      return secondaryHighlight;
    });

    this.selectedDataPointChange.emit(dataPoint);
  }

  private getHighlightWidth(): number {
    if (!this.xAxis) {
      return 0;
    }
    const xAxisElements = this.xAxis.domain();
    return this.dataPointContentWidth / xAxisElements.length;
  }

  private setHighlightCoordinates(
    highlight: Selection<SVGRectElement, unknown, HTMLElement, unknown>,
    identifier: string
  ): void {
    if (!this.xAxis) {
      return;
    }
    highlight
      .attr('x', this.xAxis(identifier) ?? null)
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private setSecondaryHighlightCoordinates(
    secondaryHighlight: Selection<SVGRectElement, unknown, HTMLElement, unknown>,
    secondaryIndex: number
  ): void {
    if (!this.xAxis) {
      return;
    }
    const xAxisElements = this.xAxis.domain();
    secondaryHighlight
      .attr('x', this.xAxis(xAxisElements[secondaryIndex]) ?? null)
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private generateHeatmapTiles(
    data: GroupedDataPoints,
    x: ScaleBand<string>,
    y: ScaleBand<string>,
    yAxisElements = y.domain()
  ): void {
    const _this = this;
    const categoryDisabledStatus = this.getCategoryDisabledStatus();
    let container: SVGGSelection = this.dataPointContainer;

    if (container.empty()) {
      container = d3.select(this.heatmapSelector).append('g').classed('data-point-container', true);
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
      .classed('disabled', (dataPoint: IDataPoint) => categoryDisabledStatus[dataPoint.color])
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
    const legend = d3.select(this.heatmapSelector).append('g').attr('class', 'legend-container');
    let xCoordinate = 0;
    for (const category of this.legendItems) {
      const legendItem = legend
        .append('g')
        .classed('legend-item', true)
        .on('click', () => {
          const isDisabled = this.disableLegend(legendItem);
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

  private disableLegend(legendItem: SVGGSelection): boolean {
    const circle = legendItem.select('circle');
    const isDisabled = circle.classed('disabled');
    circle.classed('disabled', !isDisabled);
    return !isDisabled;
  }

  private getCategoryDisabledStatus(): { [category: string]: boolean } {
    const legendContainer = this.legendContainer;

    return this.legendItems.reduce(
      (categoryStatus: { [category: string]: boolean }, category: EvaluationResultType) => {
        categoryStatus[category] = !legendContainer.select(`.legend-circle.${category}.disabled`).empty();
        return categoryStatus;
      },
      {}
    );
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
    this.onScroll();
  }

  private expandHeatmap(): void {
    if (!this.xAxis || !this.yAxis) {
      return;
    }
    const slis = Object.keys(this.groupedData);
    this.setHeight(slis.length);
    // TODO: also update xAxis? Because we have "score" there won't be any new dates
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
    const yAxis = this.callYAxis(this.yAxisContainer, this.yAxis);
    this.setYAxisStyle(yAxis);
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
}
/* eslint-enable @typescript-eslint/no-this-alias */
