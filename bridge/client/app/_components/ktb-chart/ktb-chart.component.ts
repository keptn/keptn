import {
  AfterViewInit,
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  HostListener,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { axisBottom, axisLeft, axisRight, BaseType, line, ScaleLinear, scaleLinear, select, Selection } from 'd3';
import { v4 as uuid } from 'uuid';
import { BehaviorSubject } from 'rxjs';
import { ChartItem, ChartItemPoint } from '../../_interfaces/chart';
import { getColor, getIconStyle, getTooltipPosition, replaceSpace } from './ktb-chart-utils';

type BarPoint = [number, number, string];
type SVGGSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type SVGPath = Selection<SVGPathElement, [number, number][], HTMLElement, unknown>;
type SVGRect = Selection<SVGRectElement, BarPoint, BaseType, unknown>;
type LinearScale = ScaleLinear<number, number, never>;
type Margin = { top: number; right: number; left: number; bottom: number };
type MetricValue = { label: string; value: number };
type ToolTipState = {
  visible: boolean;
  top: number;
  left: number;
  label: string;
  metricValues: MetricValue[];
};

const _height = 400;
const margin: Margin = { top: 60, right: 60, left: 50, bottom: 100 };
const yTicks = [25, 50, 75, 100];

@Component({
  selector: 'ktb-chart',
  templateUrl: './ktb-chart.component.html',
  styleUrls: ['./ktb-chart.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbChartComponent implements AfterViewInit, OnChanges {
  public readonly uniqueId = `chart-${uuid()}`;
  private readonly chartSelector = `#${this.uniqueId}`;

  @Input()
  public chartItems: ChartItem[] = [];

  @Input()
  public xLabels: Record<number, string> = {};

  @Input()
  public xTooltipLabels: Record<number, string> = {};

  private xScale = scaleLinear();
  private yScaleLeft = scaleLinear();
  private yScaleRight = scaleLinear();

  private xAxisGroup: SVGGSelection | undefined;
  private yAxisGroupLeft: SVGGSelection | undefined;
  private yAxisGroupRight: SVGGSelection | undefined;

  private paths: SVGPath[] = [];
  private rects: SVGRect[] = [];

  @ViewChild('tooltip', { static: false })
  tooltip!: ElementRef;
  private tooltipState = new BehaviorSubject<ToolTipState>({
    visible: false,
    top: 0,
    left: 0,
    label: '',
    metricValues: [],
  });
  public tooltipState$ = this.tooltipState.asObservable();

  getIconStyle = getIconStyle;
  replaceSpace = replaceSpace;

  constructor(private elementRef: ElementRef, private cdr: ChangeDetectorRef) {}

  public ngAfterViewInit(): void {
    this.init();
    this.onResize();
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  public ngOnChanges(_changes: SimpleChanges): void {
    this.draw();
  }

  private init(): void {
    const svg = select(this.chartSelector);
    this.xAxisGroup = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-x');
    this.yAxisGroupLeft = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-y-left');
    this.yAxisGroupRight = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-y-right');
  }

  @HostListener('window:resize', ['$event'])
  private onResize(): void {
    this.draw();
  }

  public onMousemove(event: MouseEvent, data: BarPoint): void {
    const xValue = data[0];
    const elements = document.elementsFromPoint(event.clientX, event.clientY);
    const barArea = elements.filter((e) => e.tagName === 'rect' && e.classList.contains('area'))[0];

    if (!barArea) {
      return;
    }

    const { top, left } = getTooltipPosition(
      { width: window.innerWidth ?? 0, height: window.innerHeight ?? 0 },
      this.tooltip.nativeElement.getBoundingClientRect(),
      barArea.getClientRects()[0]
    );
    const label = this.xTooltipLabels[xValue] ?? this.xLabels[xValue] ?? xValue + '';
    const addMetricValue = (cur: MetricValue[], item: ChartItem): MetricValue[] => {
      const point = item.points.find((p) => p.x === xValue);
      const itemLabel = item.label ?? item.identifier;
      const alreadyInList = cur.find((v) => v.label === label);
      const metricValue = !!point && !alreadyInList ? { label: itemLabel, value: point.y } : undefined;
      return metricValue ? [...cur, metricValue] : cur;
    };
    const metricValues = this.chartItems
      .filter((item) => !(item.invisible === true))
      .reduce(addMetricValue, [] as MetricValue[]);

    this.tooltipState.next({ ...this.tooltipState.getValue(), top, left, label, metricValues });
    this.cdr.detectChanges();
  }

  private draw(): void {
    const { width, height } = this.getAvailableSpace();
    const svg = select(this.chartSelector);
    svg
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('width', width)
      .attr('height', height)
      .attr('preserveAspectRatio', 'xMinYMin meet');

    const xMaxValue = this.getMaxValue((p) => p.x);
    this.xScale.domain([-0.5, xMaxValue + 0.5]).range([margin.left, width - margin.right]);
    const xTicks = this.xScale.ticks().filter((n) => Number.isInteger(n));
    this.yScaleLeft.domain([0, 100]).range([height - margin.bottom, margin.top]);
    this.yScaleRight.domain([0, this.getMaxValue((p) => p.y)]).range([height - margin.bottom, margin.top]);

    if (this.xAxisGroup) {
      this.xAxisGroup.call(
        axisBottom(this.xScale)
          .tickValues(xTicks)
          .tickFormat((d) => {
            const label = this.xLabels[d.valueOf()];
            return label ?? d;
          })
      );
      this.xAxisGroup.attr('transform', `translate(0, ${height - margin.bottom})`);
      this.xAxisGroup.selectAll('text').attr('transform', 'translate(-10,0)rotate(-45)').style('text-anchor', 'end');
    }

    if (this.yAxisGroupLeft) {
      this.yAxisGroupLeft.attr('transform', `translate(${margin.left}, ${0})`);
      this.yAxisGroupLeft.call(
        axisLeft(this.yScaleLeft)
          .tickSize(-(width - margin.left - margin.right))
          .tickValues(yTicks)
      );
    }

    if (this.yAxisGroupRight) {
      this.yAxisGroupRight.attr('transform', `translate(${width - margin.right}, ${0})`);
      const yTickValuesRight = yTicks.map((t) => this.yScaleRight.invert(this.yScaleLeft(t)));
      this.yAxisGroupRight.call(axisRight(this.yScaleRight).tickValues(yTickValuesRight));
    }

    this.paths.forEach((p) => p.remove());
    this.rects.forEach((r) => r.remove());

    this.chartItems.forEach((item, index) => {
      if (item.invisible) {
        return;
      }
      if (item.type === 'score-bar') {
        this.drawBar(item, height);
        return;
      }
      const yScale = item.type === 'score-line' ? this.yScaleLeft : this.yScaleRight;
      this.drawLine(item, yScale, getColor(index));
    });
    this.drawArea(xMaxValue + 1, height);
  }

  private drawLine(item: ChartItem, yScale: LinearScale, color: string): void {
    const svg = select(this.chartSelector);
    const points: [number, number][] = item.points.map((d) => [d.x, d.y]);
    const path = svg
      .append('path')
      .datum(points)
      .attr('fill', 'none')
      .attr('stroke', color)
      .attr('stroke-width', 1.5)
      .attr('class', 'line red')
      .attr('uitestid', `line-${replaceSpace(item.identifier)}`)
      .attr(
        'd',
        line()
          .x((d) => this.xScale(d[0]))
          .y((d) => yScale(d[1]))
      );
    this.paths.push(path);
  }

  private drawBar(item: ChartItem, height: number): void {
    const svg = select(this.chartSelector);
    const points: BarPoint[] = item.points.map((d) => [d.x, d.y, d.color ?? '#7e7e7e']);
    const selection = svg
      .selectAll()
      .data(points)
      .enter()
      .append('g')
      .attr('uitestid', `bar-${replaceSpace(item.identifier)}`);

    const rect = selection
      .append('rect')
      .attr('x', (d) => this.xScale(d[0]) - 1)
      .attr('y', (d) => this.yScaleLeft(d[1]))
      .attr('width', 3)
      .attr('height', (d) => height - margin.bottom - this.yScaleLeft(d[1]))
      .attr('fill', (d) => d[2])
      .attr('uitestid', (d) => `bar-${replaceSpace(item.identifier)}-${d[0]}`);
    this.rects.push(rect);
  }

  private drawArea(xSize: number, height: number): void {
    const svg = select(this.chartSelector);
    const points: BarPoint[] = [...Array(xSize).keys()].map((x) => [x, 100, 'transparent']);
    const selection = svg.selectAll().data(points).enter().append('g').attr('uitestid', `area`);

    const rect = selection
      .append('rect')
      .on('mouseenter', () => this.tooltipState.next({ ...this.tooltipState.getValue(), visible: true }))
      .on('mousemove', (event: MouseEvent, data) => this.onMousemove(event, data))
      .on('mouseleave', () => this.tooltipState.next({ ...this.tooltipState.getValue(), visible: false }))
      .attr('x', (d) => this.xScale(d[0] - 0.5))
      .attr('y', (d) => this.yScaleLeft(d[1]))
      .attr('width', this.xScale(2) - this.xScale(1))
      .attr('height', (d) => height - margin.bottom - this.yScaleLeft(d[1]))
      .attr('fill', (d) => d[2])
      .attr('class', 'area')
      .attr('uitestid', (d) => `area-${d[0]}`);

    this.rects.push(rect);
  }

  public hideChartItem(item: ChartItem): void {
    if (this.isHidingAllowed(item)) {
      item.invisible = !item.invisible;
      this.draw();
    }
  }

  private isHidingAllowed(item: ChartItem): boolean {
    if (item.invisible) {
      return true;
    }
    const itemsVisible = this.chartItems.filter((i) => !(i.invisible === true)).length;
    return itemsVisible > 1;
  }

  private getAvailableSpace(): { width: number; height: number } {
    const parentElement: HTMLElement = this.elementRef.nativeElement;
    const availableSpace = parentElement.getBoundingClientRect();
    return { width: availableSpace.width, height: _height };
  }

  private getMaxValue(fn: (p: ChartItemPoint) => number): number {
    return this.chartItems
      .filter((i) => !(i.invisible === true))
      .reduce((prev, cur) => Math.max(prev, ...cur.points.map(fn)), 0);
  }
}
