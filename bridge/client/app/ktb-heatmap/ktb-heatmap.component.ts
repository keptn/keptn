/* eslint-disable @typescript-eslint/no-this-alias */
import { Component, ElementRef, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import * as d3 from 'd3';
import { BaseType, ScaleBand, Selection, ValueFn } from 'd3';
import { ResultTypes } from '../../../shared/models/result-types';

export interface DataPoint {
  date: string;
  sli: string;
  value: number; // or tooltip: { allNeededValues }
  color: ResultTypes;
  comparedIndices: number[];
  /**
   * Unique identifier like keptnContext that can be used on tileSelected
   */
  identifier: string;
}

type HeatmapSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type HeatmapTooltip = Selection<HTMLDivElement, unknown, HTMLElement, unknown>;
type HeatMapMouseFunction = (this: BaseType, event: MouseEvent, d: DataPoint) => void;

@Component({
  selector: 'ktb-heatmap',
  templateUrl: './ktb-heatmap.component.html',
  styleUrls: ['./ktb-heatmap.component.scss'],
})
export class KtbHeatmapComponent implements OnInit, OnDestroy {
  private _data: DataPoint[] = [];
  private heatmap?: Selection<SVGGElement, unknown, HTMLElement, unknown>;
  private xAxis?: ScaleBand<string>;
  private readonly chartSelector = 'div#myChart';
  private readonly yAxisLabelWidth = 100;
  private readonly xAxisLabelWidth = 150;
  private width = 1920; // - margin-left - margin-right. Margin-left will be the space for xAxis labels
  private height = 40;
  private heightPerSli = 40; // 40 per SLI
  private highlight?: Selection<SVGRectElement, unknown, HTMLElement, unknown>;
  private secondaryHighlights: Selection<SVGRectElement, unknown, HTMLElement, unknown>[] = [];
  private _selectedDataPoint?: DataPoint;
  private mouseCoordinates = { x: 0, y: 0 };
  private scrollListener?: (this: BaseType, _evt: Event) => void;
  private mouseMoveListener?: (this: Document, _evt: MouseEvent) => void;

  @Output() selectedDataPointChange = new EventEmitter<DataPoint>();
  // unsure about:
  // should tileSelected emit the datapoint or just the identifier?
  // Re-positioning of tooltip only on hover-item-change?

  // TODO:
  //  - Create <ktb-heatmap-tooltip #myTooltip>, get it via ViewChild and trigger show/hide with correct x and y coordinates and dataPoint.
  //    Check if myComponentRef.attr.transform(x,y) can be used
  //  - Only show every xth date if there are too many dataPoints
  //  - Consider case for hide/show SLIs if there are more than 10
  //  - Show "warning", "failed", "succeeded" legend underneath heatmap. Also filter?
  //  - There is also an INFO state. Not only fail, warn, pass.
  //  - Allow filtering by clicking on the legend
  //  - Remove testing data afterwards
  //  - tooltip: repositioning on too much on left, top,

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

  constructor(private elementRef: ElementRef) {
    this.mouseMoveListener = this.getMouseMoveListener();
    document.addEventListener('mousemove', this.mouseMoveListener);
  }

  public ngOnInit(): void {
    this.dataPoints = this.generateTestData(12, 50); // TODO: remove testing data afterwards
    this.click(this.dataPoints[5]);
  }

  private createHeatmap(data: DataPoint[]): void {
    this.setUniqueHeaders(data, 'date', 'sli');
    this.setUniqueHeaders(data, 'sli', 'date');
    const slis = Array.from(new Set(data.map((dataPoint) => dataPoint.sli)));
    const dates = Array.from(new Set(data.map((dataPoint) => dataPoint.date)));

    this.height = slis.length * this.heightPerSli;
    this.heatmap = d3
      .select(this.chartSelector)
      .classed('svg-container', true)
      .append('svg')
      .attr('preserveAspectRatio', 'xMinYMin meet')
      .attr('viewBox', `0 0 ${this.width + this.yAxisLabelWidth} ${this.height + this.xAxisLabelWidth}`)
      .classed('svg-content-responsive', true)
      .append('g')
      .attr('transform', 'translate(' + this.yAxisLabelWidth + ',' + 0 + ')');

    this.setData(this.heatmap, data, dates, slis);
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
          color: this.getColor(Math.floor(Math.random() * 3)),
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

  private setUniqueHeaders(dataPoints: DataPoint[], key: 'sli' | 'date', compare: 'sli' | 'date'): void {
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
    yAxisElements: string[]
  ): void {
    this.xAxis = this.addXAxis(heatmap, xAxisElements);
    const y = this.addYAxis(heatmap, yAxisElements);
    const tooltip = this.buildTooltip();
    this.generateHeatmapTiles(heatmap, data, this.xAxis, y, tooltip);
  }

  private getColor(value: number): ResultTypes {
    if (value === 0) {
      return ResultTypes.FAILED;
    }
    if (value === 1) {
      return ResultTypes.WARNING;
    }
    return ResultTypes.PASSED;
  }

  private addXAxis(heatmap: HeatmapSelection, dates: string[]): ScaleBand<string> {
    const x = d3.scaleBand().range([0, this.width]).domain(dates);
    heatmap
      .append('g')
      .attr('class', 'x-axis-container')
      .attr('transform', `translate(0, ${this.height})`) // TODO: can be increased to add room for a button
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
    return function (
      this: BaseType,
      dataPoint: unknown,
      index: number,
      dataPoints: BaseType[] | ArrayLike<BaseType>
    ): void {
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
    const y = d3.scaleBand().range([this.height, 0]).domain(slis).paddingInner(0.03);
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

  private mouseOver(tooltip: HeatmapTooltip): HeatMapMouseFunction {
    return function (this: BaseType, _event: MouseEvent, _d: DataPoint): void {
      tooltip.classed('show', true);
    };
  }

  private mouseLeave(tooltip: HeatmapTooltip): HeatMapMouseFunction {
    return function (this: BaseType, _event: MouseEvent, _d: DataPoint): void {
      tooltip.classed('show', false);
    };
  }

  private mouseMove(tooltip: HeatmapTooltip): HeatMapMouseFunction {
    return function (event: MouseEvent, d: DataPoint): void {
      tooltip
        .html('The exact value of<br>this cell is: ' + d.value)
        .style('left', event.x + window.scrollX + 'px')
        .style('top', event.y + window.scrollY + 'px');
    };
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
    const xAxis = this.xAxis;
    const xAxisElements = xAxis.domain();
    const width = this.width / xAxisElements.length;

    this.highlight = heatmap
      .append('rect')
      .attr('x', this.xAxis(dataPoint.date) ?? null)
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', width)
      .attr('class', 'highlight-primary');

    this.secondaryHighlights = dataPoint.comparedIndices.map((secondary) =>
      heatmap
        .append('rect')
        .attr('x', xAxis(xAxisElements[secondary]) ?? null)
        .attr('y', 0)
        .attr('height', this.height)
        .attr('width', width)
        .attr('class', 'highlight-secondary')
    );

    this.selectedDataPointChange.emit(dataPoint);
  }

  private generateHeatmapTiles(
    heatmap: HeatmapSelection,
    data: DataPoint[],
    x: ScaleBand<string>,
    y: ScaleBand<string>,
    tooltip: HeatmapTooltip
  ): void {
    const _this = this;
    heatmap
      .selectAll()
      .data(data)
      .join('rect')
      .attr('x', (dataPoint) => x(dataPoint.date) ?? null)
      .attr('y', (dataPoint) => y(dataPoint.sli) ?? null)
      .attr('class', (dataPoint) => `data-point ${dataPoint.color}`)
      .attr('uitestid', (dataPoint) => `ktb-heatmap-tile-${dataPoint.sli.replace(/ /g, '-')}`) // TODO: do we need this?
      .attr('width', x.bandwidth())
      .attr('height', y.bandwidth())
      .on('click', function (this: BaseType, _event: MouseEvent, dataPoint: DataPoint) {
        _this.click(dataPoint);
      })
      .on('mouseover', this.mouseOver(tooltip))
      .on('mousemove', this.mouseMove(tooltip))
      .on('mouseleave', this.mouseLeave(tooltip));

    this.scrollListener = this.getScrollListener(tooltip);
    document.addEventListener('scroll', this.scrollListener, false);
  }

  private getMouseMoveListener(): (this: Document, _evt: MouseEvent) => void {
    const _this = this;
    return function (this: Document, event: MouseEvent) {
      _this.mouseCoordinates = {
        x: event.x,
        y: event.y,
      };
    };
  }

  private getScrollListener(tooltip: HeatmapTooltip): (this: BaseType, _evt: Event) => void {
    const _this = this;
    return function (this: BaseType, _evt: Event): void {
      const element = document.elementFromPoint(_this.mouseCoordinates.x, _this.mouseCoordinates.y);
      if (!element || !(element instanceof SVGRectElement)) {
        tooltip.classed('show', false);
        return;
      }
      const dt = d3.select(element)?.datum() as DataPoint | undefined;
      const mouseEvent = new MouseEvent('move', {
        clientY: _this.mouseCoordinates.y,
        clientX: _this.mouseCoordinates.x,
      });
      if (dt) {
        tooltip.classed('show', true);
        _this.mouseMove(tooltip).bind(this)(mouseEvent, dt);
      } else {
        tooltip.classed('show', false);
      }
    };
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
