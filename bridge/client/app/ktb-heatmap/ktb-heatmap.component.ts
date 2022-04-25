/* eslint-disable @typescript-eslint/no-this-alias */
import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Input,
  OnInit,
  Output,
  ViewChild,
} from '@angular/core';
import * as d3 from 'd3';
import { BaseType, ScaleBand, Selection, ValueFn } from 'd3';
import { ResultTypes } from '../../../shared/models/result-types';
import { DtButton } from '@dynatrace/barista-components/button';

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
export class KtbHeatmapComponent implements OnInit, AfterViewInit {
  private xAxis?: ScaleBand<string>;
  private yAxis?: ScaleBand<string>;
  private readonly chartSelector = 'div#myChart';
  private readonly heatmapSelector = `#heatmap-container`;
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
  private dataPointContentWidth = 1720; // - margin-left - margin-right. Margin-left will be the space for xAxis labels
  private height = 0;
  private readonly showMoreButtonHeight = 38;
  private highlight?: Selection<SVGRectElement, unknown, HTMLElement, unknown>;
  private secondaryHighlights: Selection<SVGRectElement, unknown, HTMLElement, unknown>[] = [];
  private _selectedDataPoint?: DataPoint;
  private mouseCoordinates = { x: 0, y: 0 };
  private scrollListener?: (this: Document, _evt: Event) => void;
  // private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  public showMoreVisible = true;
  private showMoreExpanded = false;
  private groupedData: GroupedDataPoints = {};
  private tooltipSelector = '.tooltip';

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
  //  - Consider case for hiding SLIs if there are more than 10
  //  - Remove testing data afterwards

  @Input() set dataPoints(data: DataPoint[]) {
    this.setUniqueHeaders(data, 'date', 'sli');
    this.setUniqueHeaders(data, 'sli', 'date');
    this.groupedData = data.reduce((groupedData: GroupedDataPoints, dataPoint) => {
      groupedData[dataPoint.sli] = groupedData[dataPoint.sli] || [];
      groupedData[dataPoint.sli].push(dataPoint);
      return groupedData;
    }, {});
    this.createHeatmap(this.groupedData);
  }

  @Input() set selectedDataPoint(dataPoint: DataPoint | undefined) {
    this.click(dataPoint);
  }
  get selectedDataPoint(): DataPoint | undefined {
    return this._selectedDataPoint;
  }

  get showMoreButtonTopOffset(): number {
    const heatmapHeight =
      (d3.select(this.heatmapSelector).select('#data-point-container') as SVGGSelection)
        ?.node()
        ?.getBoundingClientRect().height ?? 0;
    const heightOffset = (this.elementRef.nativeElement as HTMLElement).offsetTop;
    return heatmapHeight + heightOffset + 5;
  }

  get showMoreButtonLeftOffset(): number {
    const yAxisWidth =
      (d3.select(this.heatmapSelector).select('.y-axis-container') as SVGGSelection)?.node()?.getBoundingClientRect()
        .width ?? 0;
    const widthOffset = (this.elementRef.nativeElement as HTMLElement).offsetLeft;
    return yAxisWidth + widthOffset - 2;
  }

  public get tooltip(): HeatmapTooltip {
    return d3.select(this.chartSelector).select(this.tooltipSelector);
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
  }

  @HostListener('mousemove', ['$event'])
  private getMouseMoveListener(event: MouseEvent): void {
    this.mouseCoordinates = {
      x: event.x,
      y: event.y,
    };
  }

  @HostListener('window:scroll')
  private getScrollListener(): void {
    const tooltip = this.tooltip;

    const element = document.elementFromPoint(this.mouseCoordinates.x, this.mouseCoordinates.y);
    if (!element || !(element instanceof SVGRectElement)) {
      tooltip.classed('hidden', true);
      return;
    }
    const dt = d3.select(element)?.datum() as DataPoint | undefined;
    if (dt) {
      const mouseEvent = new MouseEvent('move', {
        clientY: this.mouseCoordinates.y,
        clientX: this.mouseCoordinates.x,
      });
      tooltip.classed('hidden', false);
      this.mouseMove(tooltip, mouseEvent, dt);
    } else {
      tooltip.classed('hidden', true);
    }
  }

  constructor(private elementRef: ElementRef) {}

  public ngOnInit(): void {
    this.dataPoints = this.generateTestData(12, 10); // TODO: remove testing data afterwards
    this.click(this.groupedData.score[1]);
  }

  public ngAfterViewInit(): void {
    this.resizeShowMoreButton();
  }

  private createHeatmap(data: GroupedDataPoints): void {
    let slis = Object.keys(data);
    this.showMoreVisible = slis.length > this.limitSliCount;
    if (this.showMoreVisible) {
      slis = slis.slice(slis.length - this.limitSliCount, slis.length);
    }
    const allDates = slis.reduce((dates: string[], sli: string) => {
      return [...dates, ...data[sli].map((dataPoint) => dataPoint.date)];
    }, []);
    const dates = Array.from(new Set(allDates));

    this.setHeight(slis.length);
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = d3.select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    svg.append('g').attr('id', 'heatmap-container').attr('transform', `translate(${this.yAxisLabelWidth}, 0)`);

    this.setData(data, dates, slis, this.showMoreVisible);
    this.createLegend();
  }

  private setAndGetAvailableSpace(): { height: number; width: number } {
    const availableSpace = (this.elementRef.nativeElement.parentNode as HTMLElement).getBoundingClientRect();
    const height =
      this.height + this.xAxisLabelWidth + this.legendHeight + (this.showMoreVisible ? this.showMoreButtonHeight : 0);
    this.dataPointContentWidth = availableSpace.width - this.yAxisLabelWidth;
    const width = availableSpace.width;
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
      const xAxisContainer = d3.select('.x-axis-container') as SVGGSelection;
      xAxisContainer
        .call(d3.axisBottom(this.xAxis))
        .attr('transform', `translate(0, ${this.height + (this.showMoreVisible ? this.showMoreButtonHeight : 0)})`);
    }
  }

  private resizeDataPoints(): void {
    if (this.xAxis && this.yAxis) {
      this.setDataPointCoordinates(d3.selectAll('.data-point') as unknown as HeatmapTiles, this.xAxis, this.yAxis);
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
    const legend = d3.select('#legend-container') as SVGGSelection;
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
    const tooltip = this.buildTooltip();
    this.generateHeatmapTiles(data, this.xAxis, this.yAxis, tooltip);
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

  private mouseOver(tooltip: HeatmapTooltip): void {
    tooltip.classed('hidden', false);
  }

  private mouseLeave(tooltip: HeatmapTooltip): void {
    tooltip.classed('hidden', true);
  }

  private mouseMove(tooltip: HeatmapTooltip, event: MouseEvent, dataPoint: DataPoint): void {
    tooltip
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
    tooltip: HeatmapTooltip,
    yAxisElements = y.domain(),
    createContainer = true
  ): void {
    const heatmap = d3.select(this.heatmapSelector) as SVGGSelection;
    let container: SVGGSelection;
    if (createContainer) {
      container = heatmap.append('g').attr('id', 'data-point-container');
    } else {
      container = heatmap.select('#data-point-container');
    }
    const dataPoints = container
      .selectAll()
      .data(yAxisElements)
      .enter()
      .append('g')
      .attr('id', (sli) => sli.replace(/ /g, '-'))
      .selectAll()
      .data((key) => data[key])
      .join('rect')
      .attr('class', (dataPoint) => dataPoint.color)
      .classed('data-point', true)
      .attr('uitestid', (dataPoint) => `ktb-heatmap-tile-${dataPoint.date.replace(/ /g, '-')}`) // TODO: do we need this?
      .on('click', (_event: PointerEvent, dataPoint: DataPoint) => this.click(dataPoint))
      .on('mouseover', () => this.mouseOver(tooltip))
      .on('mousemove', (event: MouseEvent, dataPoint: DataPoint) => this.mouseMove(tooltip, event, dataPoint))
      .on('mouseleave', () => this.mouseLeave(tooltip));

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
    const legend = heatmap.append('g').attr('id', 'legend-container');
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

    // in this case the yAxis height changes
    // resize xAxis, legend and showMoreButton

    if (this.showMoreExpanded) {
      return this.expandHeatmap();
    }
    this.collapseHeatmap();
  }

  private expandHeatmap(): void {
    const slis = Object.keys(this.groupedData);
    this.setHeight(slis.length);
    if (this.xAxis && this.yAxis) {
      // TODO: also update xAxis? Because we have "score" there won't be any new dates
      const { width, height } = this.setAndGetAvailableSpace();

      this.yAxis = this.yAxis.range([this.height, 0]).domain(slis);
      (d3.select(this.heatmapSelector).select('.y-axis-container') as SVGGSelection).call(
        d3.axisLeft(this.yAxis).tickSize(0)
      );
      const tooltip = d3.select(this.chartSelector).select('.tooltip') as HeatmapTooltip;
      this.generateHeatmapTiles(
        this.groupedData,
        this.xAxis,
        this.yAxis,
        tooltip,
        slis.slice(0, slis.length - this.limitSliCount),
        false
      );

      this.resizeSvg(width, height);
      this.resizeShowMoreButton();
      this.resizeXAxis();
      this.resizeHighlights();
      this.resizeLegend();
    }
  }

  private collapseHeatmap(): void {
    this.setHeight(this.limitSliCount);
    // remove all dataPoints with SLI-index > 10
  }

  private setHeight(elementCount: number): void {
    this.height = elementCount * this.heightPerSli;
  }
}
/* eslint-enable @typescript-eslint/no-this-alias */
