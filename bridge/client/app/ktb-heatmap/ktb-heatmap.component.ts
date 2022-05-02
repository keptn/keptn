/* eslint-disable @typescript-eslint/no-this-alias */
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import * as d3 from 'd3';
import { BaseType, ScaleBand, Selection, ValueFn } from 'd3';
import { ResultTypes } from '../../../shared/models/result-types';
import { DtButton } from '@dynatrace/barista-components/button';
import { v4 } from 'uuid';

export interface DataPoint {
  date: string;
  sli: string;
  value: number; // or tooltip: { allNeededValues }
  color: EvaluationResultType;
  comparedIndices: number[];
  /**
   * Unique identifier like keptnContext that can be used on tileSelected
   */
  identifier: string;
}

type SVGGSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type HeatmapTooltip = Selection<HTMLDivElement, unknown, HTMLElement, unknown>;
type HeatmapTiles = Selection<SVGRectElement | null, DataPoint, SVGGElement, unknown>;
enum EvaluationResultTypeExtension {
  INFO = 'info',
}
type EvaluationResultType = ResultTypes | EvaluationResultTypeExtension;
type GroupedDataPoints = { [sli: string]: DataPoint[] };

@Component({
  selector: 'ktb-heatmap',
  templateUrl: './ktb-heatmap.component.html',
  styleUrls: ['./ktb-heatmap.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHeatmapComponent implements OnInit, OnDestroy, AfterViewInit {
  public readonly uniqueId = `heatmap-${v4()}`;
  private readonly chartSelector = `#${this.uniqueId}`;
  private readonly heatmapSelector = `${this.chartSelector} .heatmap-container`;
  private readonly tooltipSelector = `${this.chartSelector} .tooltip`;
  private readonly firstSliPadding = 6; // "score" will then be 6px smaller than the rest.
  private readonly yAxisLabelWidth = 100;
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
  // has to be globally instead of component bound, else scrolling into it will not have any mouse coordinates
  private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  private xAxis?: ScaleBand<string>;
  private yAxis?: ScaleBand<string>;
  private dataPointContentWidth = 1720; // - margin-left - margin-right. Margin-left will be the space for xAxis labels
  private height = 0;
  private highlight?: Selection<SVGRectElement, unknown, HTMLElement, unknown>;
  private secondaryHighlights: Selection<SVGRectElement, unknown, HTMLElement, unknown>[] = [];
  private _selectedDataPoint?: DataPoint;
  private mouseCoordinates = { x: 0, y: 0 };
  public showMoreVisible = true;
  public showMoreExpanded = false;
  private groupedData: GroupedDataPoints = {};

  @ViewChild('showMoreButton', { static: false }) showMoreButton!: DtButton;
  @Output() selectedDataPointChange = new EventEmitter<DataPoint>();
  // unsure about:
  // should tileSelected emit the datapoint or just the identifier?
  // Re-positioning of tooltip only on hover-item-change?
  //

  // TODO:
  //  - Create <ktb-heatmap-tooltip #myTooltip>, get it via ViewChild and trigger show/hide with correct x and y coordinates and dataPoint.
  //    Check if myComponentRef.attr.transform(x,y) can be used
  //    repositioning if too far on the left/top,
  //  - Only show every xth date if there are too many dataPoints?
  //  - disable Tooltip if item in legend is disabled
  //  - Remove testing data afterwards

  @Input() set dataPoints(data: DataPoint[]) {
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

  @Input() set selectedDataPoint(dataPoint: DataPoint | undefined) {
    this.click(dataPoint);
  }
  get selectedDataPoint(): DataPoint | undefined {
    return this._selectedDataPoint;
  }

  get dataPointContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.data-point-container') as SVGGSelection;
  }

  get showMoreButtonTopOffset(): number {
    const heatmapHeight = this.dataPointContainer?.node()?.getBoundingClientRect().height ?? 0;
    return heatmapHeight + this.showMoreButtonPadding;
  }

  get yAxisContainer(): SVGGSelection {
    return d3.select(this.heatmapSelector).select('.y-axis-container') as SVGGSelection;
  }

  get showMoreButtonLeftOffset(): number {
    const yAxisWidth = this.yAxisContainer?.node()?.getBoundingClientRect().width ?? 0;
    return yAxisWidth - 2;
  }

  public get tooltip(): HeatmapTooltip {
    return d3.select(this.tooltipSelector);
  }

  constructor(private elementRef: ElementRef) {
    this.mouseMoveListener = (event: MouseEvent): void => this.onMouseMove(event);
    document.addEventListener('mousemove', this.mouseMoveListener);
  }

  public ngOnInit(): void {}

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

    const dt = d3.select(element)?.datum() as DataPoint | undefined;

    if (!dt) {
      this.setTooltipVisibility(false);
      return;
    }

    const mouseEvent = new MouseEvent('move', {
      clientY: y,
      clientX: x,
    });
    this.setTooltipVisibility(true);
    this.mouseMove(mouseEvent, dt);
  }

  private setTooltipVisibility(visible: boolean): void {
    this.tooltip.classed('hidden', !visible);
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
    const dates = Array.from(new Set(allDates));

    this.setHeight(slis.length);
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = d3.select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    svg.append('g').classed('heatmap-container', true).attr('transform', `translate(${this.yAxisLabelWidth}, 0)`);

    this.setData(data, dates, slis, this.showMoreVisible);
    this.createLegend();
  }

  private setAndGetAvailableSpace(): { height: number; width: number } {
    const availableSpace = (this.elementRef.nativeElement.parentNode as HTMLElement).getBoundingClientRect();
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
    if (this.xAxis) {
      this.xAxis = this.xAxis.range([0, this.dataPointContentWidth]);
      const xAxisContainer = d3.select(this.chartSelector).select('.x-axis-container') as SVGGSelection;
      xAxisContainer
        .call(d3.axisBottom(this.xAxis))
        .attr('transform', `translate(0, ${this.height + (this.showMoreVisible ? this.showMoreButtonHeight : 0)})`);
    }
  }

  private resizeDataPoints(): void {
    if (this.xAxis && this.yAxis) {
      this.setDataPointCoordinates(
        d3.select(this.chartSelector).selectAll('.data-point') as unknown as HeatmapTiles,
        this.xAxis,
        this.yAxis
      );
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
      const htmlElement = this.showMoreButton._elementRef.nativeElement as HTMLElement;

      htmlElement.style.top = `${this.showMoreButtonTopOffset}px`;
      htmlElement.style.left = `${this.showMoreButtonLeftOffset}px`;
      htmlElement.style.width = `${this.dataPointContentWidth}px`;
    }
  }

  private resizeLegend(): void {
    const legend = d3.select(this.heatmapSelector).select('.legend-container') as SVGGSelection;
    const fullLength = legend.node()?.getBoundingClientRect().width ?? 0;
    const centerXPosition = (this.dataPointContentWidth - fullLength) / 2;
    const yPosition = this.height + this.xAxisLabelWidth + 10 + (this.showMoreVisible ? this.showMoreButtonHeight : 0);
    legend.attr('transform', `translate(${centerXPosition}, ${yPosition})`);
  }

  private generateTestData(sliCounter: number, counter: number): DataPoint[] {
    const categories = [];
    for (let i = 0; i < sliCounter; ++i) {
      categories.push(`response time p${i}`);
    }
    categories.push('score');
    const data: DataPoint[] = [];
    const dateMillis = new Date().getTime();

    // data.push({
    //   date: new Date(dateMillis).toISOString(),
    //   sli: categories[0],
    //   color: this.getColor(Math.floor(Math.random() * 3)),
    //   value: Math.random(),
    //   identifier: `keptnContext_${0}`,
    //   comparedIndices: [],
    // });
    //
    // data.push({
    //   date: new Date(dateMillis).toISOString(),
    //   sli: categories[1],
    //   color: this.getColor(Math.floor(Math.random() * 3)),
    //   value: Math.random(),
    //   identifier: `keptnContext_${0}`,
    //   comparedIndices: [],
    // });

    for (const category of categories) {
      for (let i = 0; i < counter; ++i) {
        data.push({
          date: new Date(dateMillis + i).toISOString(),
          sli: category,
          color: this.getColor(Math.floor(Math.random() * 4)),
          value: Math.random(),
          identifier: `keptnContext_${i}`,
          comparedIndices: [],
        });
      }
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

  private setUniqueHeaders(dataPoints: DataPoint[], key: 'date', compare: 'sli'): void;
  private setUniqueHeaders(dataPoints: DataPoint[], key: 'sli', compare: 'date'): void;
  private setUniqueHeaders(dataPoints: DataPoint[], key: 'sli' | 'date', compare: 'date' | 'sli'): void {
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
    const heatmap = d3.select(this.heatmapSelector) as SVGGSelection;
    this.xAxis = this.addXAxis(heatmap, xAxisElements, showMoreVisible);
    this.yAxis = this.addYAxis(heatmap, yAxisElements);
    this.buildTooltip();
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
      .attr('dy', '.15em')
      .select('.domain')
      .remove();
    return x;
  }

  private setEllipsisStyle(labelWidth: number): ValueFn<BaseType, unknown, void> {
    return function (this: BaseType): void {
      const self = d3.select(this as SVGTextContentElement);
      if (self) {
        const originalText = self.text();
        let textLength = self.node()?.getComputedTextLength() ?? 0;
        let text = originalText;

        while (textLength > labelWidth && text.length > 0) {
          text = text.slice(0, -1);
          self.text(text + '...');
          textLength = self.node()?.getComputedTextLength() ?? 0;
        }
        self.append('title').text(originalText);
      }
    };
  }

  private addYAxis(heatmap: SVGGSelection, slis: string[]): ScaleBand<string> {
    const y = d3.scaleBand().range([this.height, 0]).domain(slis);
    heatmap
      .append('g')
      .attr('class', 'y-axis-container')
      .call(d3.axisLeft(y).tickSize(0))
      .selectAll('text')
      .each(this.setEllipsisStyle(this.yAxisLabelWidth));
    return y;
  }

  private buildTooltip(): HeatmapTooltip {
    return d3.select(this.chartSelector).append('div').attr('class', 'tooltip');
  }

  private mouseOver(): void {
    this.setTooltipVisibility(true);
  }

  private mouseLeave(): void {
    this.setTooltipVisibility(false);
  }

  private mouseMove(event: MouseEvent, dataPoint: DataPoint): void {
    this.tooltip
      .html('The exact value of<br>this cell is: ' + dataPoint.value)
      .style('left', event.x + 'px')
      .style('top', event.y + 'px');
  }

  private removeHighlights(): void {
    this.highlight?.remove();
    for (const highlight of this.secondaryHighlights) {
      highlight.remove();
    }
  }

  private click(dataPoint?: DataPoint): void {
    this.removeHighlights();
    const heatmap = d3.select(this.heatmapSelector) as SVGGSelection;

    if (!this.xAxis || !dataPoint) {
      this._selectedDataPoint = undefined;
      return;
    }
    this._selectedDataPoint = dataPoint;

    this.highlight = heatmap.append('rect').attr('class', 'highlight-primary');
    this.setHighlightCoordinates(this.highlight, dataPoint.date);

    this.secondaryHighlights = dataPoint.comparedIndices.map((secondary) => {
      const secondaryHighlight = heatmap.append('rect').attr('class', 'highlight-secondary');
      this.setSecondaryHighlightCoordinates(secondaryHighlight, secondary);
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
    yAxisElements = y.domain(),
    createContainer = true
  ): void {
    const heatmap = d3.select(this.heatmapSelector) as SVGGSelection;
    let container: SVGGSelection;
    if (createContainer) {
      container = heatmap.append('g').classed('data-point-container', true);
    } else {
      container = this.dataPointContainer;
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
      .attr('uitestid', (dataPoint) => `ktb-heatmap-tile-${dataPoint.date.replace(/ /g, '-')}`) // TODO: do we need this?
      .on('click', (_event: PointerEvent, dataPoint: DataPoint) => this.click(dataPoint))
      .on('mouseover', () => this.mouseOver())
      .on('mousemove', (event: MouseEvent, dataPoint: DataPoint) => this.mouseMove(event, dataPoint))
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
    const heatmap = d3.select(this.heatmapSelector) as SVGGSelection;
    const legend = heatmap.append('g').attr('class', 'legend-container');
    let xCoordinate = 0;
    for (const category of this.legendItems) {
      const legendItem = legend
        .append('g')
        .classed('legend-item', true)
        .on('click', () => {
          this.disableLegend(heatmap, legendItem, category);
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

  private disableLegend(heatmap: SVGGSelection, legendItem: SVGGSelection, category: EvaluationResultType): void {
    const circle = legendItem.select('circle');
    const isDisabled = circle.classed('disabled');
    circle.classed('disabled', !isDisabled);

    (heatmap.selectAll('.data-point') as HeatmapTiles).each(function (this: SVGGElement | null, dataPoint: DataPoint) {
      if (this && dataPoint.color === category) {
        d3.select(this).classed('disabled', !isDisabled);
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
    const slis = Object.keys(this.groupedData);
    this.setHeight(slis.length);
    if (this.xAxis && this.yAxis) {
      // TODO: also update xAxis? Because we have "score" there won't be any new dates
      this.updateYAxis(slis);

      this.generateHeatmapTiles(this.groupedData, this.xAxis, this.yAxis, this.getHiddenSLIs(slis), false);
    }
  }

  private getHiddenSLIs(slis: string[]): string[] {
    return slis.slice(0, slis.length - this.limitSliCount);
  }

  private getLimitedSLIs(slis: string[]): string[] {
    return slis.slice(slis.length - this.limitSliCount, slis.length);
  }

  private updateYAxis(slis: string[]): void {
    if (this.yAxis) {
      this.yAxis = this.yAxis.range([this.height, 0]).domain(slis);
      this.yAxisContainer.call(d3.axisLeft(this.yAxis).tickSize(0));
    }
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
