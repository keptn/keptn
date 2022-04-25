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

type HeatmapSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type HeatmapTooltip = Selection<HTMLDivElement, unknown, HTMLElement, unknown>;
type HeatmapTiles = Selection<SVGRectElement | null, DataPoint, SVGGElement, unknown>;
enum EvaluationResultTypeExtension {
  INFO = 'info',
}
type EvaluationResultType = ResultTypes | EvaluationResultTypeExtension;

@Component({
  selector: 'ktb-heatmap',
  templateUrl: './ktb-heatmap.component.html',
  styleUrls: ['./ktb-heatmap.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHeatmapComponent implements OnDestroy, AfterViewInit {
  private _data: DataPoint[] = [];
  private heatmap?: Selection<SVGGElement, unknown, HTMLElement, unknown>;
  private xAxis?: ScaleBand<string>;
  private yAxis?: ScaleBand<string>;
  private readonly chartSelector = 'div#myChart';
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
  private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  public showMoreVisible = true;

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
  //  - Only show every xth date if there are too many dataPoints
  //  - Consider case for hide/show SLIs if there are more than 10. Add button in .html and reposition/display it with d3
  //  - Remove testing data afterwards

  @Input() set dataPoints(data: DataPoint[]) {
    this._data = data;
    this.createHeatmap(data);
  }
  get dataPoints(): DataPoint[] {
    return this._data;
  }

  @Input() set selectedDataPoint(dataPoint: DataPoint | undefined) {
    this.click(dataPoint);
  }
  get selectedDataPoint(): DataPoint | undefined {
    return this._selectedDataPoint;
  }

  get showMoreButtonTopOffset(): number {
    const heatmapHeight =
      (this.heatmap?.select('#data-point-container') as HeatmapSelection)?.node()?.getBoundingClientRect().height ?? 0;
    const heightOffset = (this.elementRef.nativeElement as HTMLElement).offsetTop;
    return heatmapHeight + heightOffset + 5;
  }

  get showMoreButtonLeftOffset(): number {
    const yAxisWidth =
      (this.heatmap?.select('.y-axis-container') as HeatmapSelection)?.node()?.getBoundingClientRect().width ?? 0;
    const widthOffset = (this.elementRef.nativeElement as HTMLElement).offsetLeft;
    return yAxisWidth + widthOffset - 2;
  }

  constructor(private elementRef: ElementRef) {
    this.mouseMoveListener = (event: MouseEvent): void => this.getMouseMoveListener(event);
    document.addEventListener('mousemove', this.mouseMoveListener);
  }

  public ngAfterViewInit(): void {
    this.dataPoints = this.generateTestData(12, 10); // TODO: remove testing data afterwards
    this.click(this.dataPoints[5]);
    this.resizeShowMoreButton();
  }

  private createHeatmap(data: DataPoint[]): void {
    this.setUniqueHeaders(data, 'date', 'sli');
    this.setUniqueHeaders(data, 'sli', 'date');
    const slis = Array.from(new Set(data.map((dataPoint) => dataPoint.sli)));
    const dates = Array.from(new Set(data.map((dataPoint) => dataPoint.date)));

    this.height = slis.length * this.heightPerSli;
    this.showMoreVisible = slis.length > this.limitSliCount;
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = d3.select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    this.heatmap = svg.append('g').attr('transform', `translate(${this.yAxisLabelWidth}, 0)`);

    this.setData(this.heatmap, data, dates, slis, this.showMoreVisible);
    this.createLegend(this.heatmap, this.showMoreVisible);
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

  private resizeSvg(width: number, height: number): void {
    d3.select(`${this.chartSelector}>svg`)
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('width', width)
      .attr('height', height);
  }

  private resizeXAxis(): void {
    if (this.xAxis) {
      this.xAxis = this.xAxis.range([0, this.dataPointContentWidth]);
      (d3.select('.x-axis-container') as HeatmapSelection).call(d3.axisBottom(this.xAxis));
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
    const legend = d3.select('#legend-container') as HeatmapSelection;
    const fullLength = legend.node()?.getBoundingClientRect().width ?? 0;
    const centerXPosition = (this.dataPointContentWidth - fullLength) / 2;
    legend.attr('transform', `translate(${centerXPosition}, ${0})`);
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
    heatmap: HeatmapSelection,
    data: DataPoint[],
    xAxisElements: string[],
    yAxisElements: string[],
    showMoreVisible: boolean
  ): void {
    this.xAxis = this.addXAxis(heatmap, xAxisElements, showMoreVisible);
    this.yAxis = this.addYAxis(heatmap, yAxisElements);
    const tooltip = this.buildTooltip();
    this.generateHeatmapTiles(heatmap, data, this.xAxis, this.yAxis, tooltip);
  }

  private addXAxis(heatmap: HeatmapSelection, dates: string[], showMoreVisible: boolean): ScaleBand<string> {
    const x = d3.scaleBand().range([0, this.dataPointContentWidth]).domain(dates);
    heatmap
      .append('g')
      .attr('class', 'x-axis-container')
      .attr('transform', `translate(0, ${this.height + (showMoreVisible ? this.showMoreButtonHeight : 0)})`) // TODO: can be increased to add room for a button
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

  private addYAxis(heatmap: HeatmapSelection, slis: string[]): ScaleBand<string> {
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

    if (!this.xAxis || !this.heatmap || !dataPoint) {
      this._selectedDataPoint = undefined;
      return;
    }
    const heatmap = this.heatmap;
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
    console.log(this.getHighlightWidth());
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
    heatmap: HeatmapSelection,
    data: DataPoint[],
    x: ScaleBand<string>,
    y: ScaleBand<string>,
    tooltip: HeatmapTooltip
  ): void {
    const dataPoints = heatmap
      .append('g')
      .attr('id', 'data-point-container')
      .selectAll()
      .data(data)
      .join('rect')
      .attr('class', (dataPoint) => dataPoint.color)
      .classed('data-point', true)
      .attr(
        'uitestid',
        (dataPoint) => `ktb-heatmap-tile-${dataPoint.sli.replace(/ /g, '-')}-${dataPoint.date.replace(/ /g, '-')}`
      ) // TODO: do we need this?
      .on('click', (_event: PointerEvent, dataPoint: DataPoint) => this.click(dataPoint))
      .on('mouseover', () => this.mouseOver(tooltip))
      .on('mousemove', (event: MouseEvent, dataPoint: DataPoint) => this.mouseMove(tooltip, event, dataPoint))
      .on('mouseleave', () => this.mouseLeave(tooltip));

    this.setDataPointCoordinates(dataPoints, x, y);
    this.scrollListener = (): void => this.getScrollListener(tooltip);
    document.addEventListener('scroll', this.scrollListener, false);
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

  private getMouseMoveListener(event: MouseEvent): void {
    this.mouseCoordinates = {
      x: event.x,
      y: event.y,
    };
  }

  private getScrollListener(tooltip: HeatmapTooltip): void {
    const element = document.elementFromPoint(this.mouseCoordinates.x, this.mouseCoordinates.y);
    if (!element || !(element instanceof SVGRectElement)) {
      tooltip.classed('hidden', true);
      return;
    }
    const dt = d3.select(element)?.datum() as DataPoint | undefined;
    const mouseEvent = new MouseEvent('move', {
      clientY: this.mouseCoordinates.y,
      clientX: this.mouseCoordinates.x,
    });
    if (dt) {
      tooltip.classed('hidden', false);
      this.mouseMove(tooltip, mouseEvent, dt);
    } else {
      tooltip.classed('hidden', true);
    }
  }

  private createLegend(heatmap: HeatmapSelection, showMoreVisible: boolean): void {
    const legendPadding = 30;
    const legend = heatmap.append('g').attr('id', 'legend-container');
    const yCoordinate = this.height + this.xAxisLabelWidth + 10 + (showMoreVisible ? this.showMoreButtonHeight : 0);
    let xCoordinate = 0;
    for (const category of this.legendItems) {
      const legendContainer = legend
        .append('g')
        .classed('legend-item', true)
        .on('click', () => {
          this.disableLegend(heatmap, legendContainer, category);
        });
      legendContainer
        .append('circle')
        .attr('cx', xCoordinate)
        .attr('cy', yCoordinate)
        .attr('r', 6)
        .classed('legend-circle', true)
        .classed(category, true);
      xCoordinate += 10;
      const text = legendContainer
        .append('text')
        .attr('x', xCoordinate)
        .attr('y', yCoordinate)
        .text(category)
        .classed('legend-text', true);
      const textWidth = text.node()?.getComputedTextLength() ?? 0;
      xCoordinate += textWidth + legendPadding;
    }
    this.resizeLegend();
  }

  private disableLegend(heatmap: HeatmapSelection, legendItem: HeatmapSelection, category: EvaluationResultType): void {
    const circle = legendItem.select('circle');
    const isDisabled = circle.classed('disabled');
    circle.classed('disabled', !isDisabled);

    (heatmap.selectAll('.data-point') as HeatmapTiles).each(function (this: SVGGElement | null, dataPoint: DataPoint) {
      if (this && dataPoint.color === category) {
        d3.select(this).classed('disabled', !isDisabled);
      }
    });
  }

  public ngOnDestroy(): void {
    this.removeScrollListener();
    this.removeMouseMoveListener();
  }

  private removeScrollListener(): void {
    if (!this.scrollListener) {
      return;
    }
    document.removeEventListener('scroll', this.scrollListener);
  }

  private removeMouseMoveListener(): void {
    if (!this.mouseMoveListener) {
      return;
    }
    document.removeEventListener('mousemove', this.mouseMoveListener);
  }
}
/* eslint-enable @typescript-eslint/no-this-alias */
