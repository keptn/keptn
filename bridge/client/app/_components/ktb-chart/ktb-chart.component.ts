import {
  AfterViewInit,
  ChangeDetectionStrategy,
  Component,
  ElementRef,
  HostListener,
  Input,
  OnChanges,
  SimpleChanges,
} from '@angular/core';
import { axisBottom, axisLeft, axisRight, BaseType, line, ScaleLinear, scaleLinear, select, Selection } from 'd3';
import { v4 as uuid } from 'uuid';

type SVGGSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type SVGPath = Selection<SVGPathElement, [number, number][], HTMLElement, unknown>;
type SVGRect = Selection<SVGRectElement, [number, number, string], BaseType, unknown>;
type LinearScale = ScaleLinear<number, number, never>;
type Margin = { top: number; right: number; left: number; bottom: number };

export interface ChartItemPoint {
  x: number;
  y: number;
  identifier: string;
  color?: string;
}

export interface ChartItem {
  identifier: string;
  label?: string;
  type: 'metric-line' | 'score-bar' | 'score-line';
  invisible?: boolean;
  points: ChartItemPoint[];
}

const _height = 400;
const _margin: Margin = { top: 60, right: 60, left: 50, bottom: 100 };
const yTicks = [25, 50, 75, 100];

const colors = [
  '#9355b7',
  '#7dc540',
  '#14a8f5',
  '#f5d30f',
  '#ef651f',
  '#dc172a',
  '#00b9cc',
  '#522273',
  '#1f7e1e',
  '#004999',
  '#ab8300',
  '#8d380f',
  '#93060e',
  '#006d75',
];

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

  xScale = scaleLinear();
  yScaleLeft = scaleLinear();
  yScaleRight = scaleLinear();

  xAxisGroup: SVGGSelection | undefined;
  yAxisGroupLeft: SVGGSelection | undefined;
  yAxisGroupRight: SVGGSelection | undefined;

  private paths: SVGPath[] = [];
  private rects: SVGRect[] = [];

  constructor(private elementRef: ElementRef) {}

  ngAfterViewInit(): void {
    this.init();
    this.onResize();
  }

  private init(): void {
    const svg = select(this.chartSelector).on('mousemove', (event: MouseEvent) => this.onMousemove(event));
    this.xAxisGroup = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-x');
    this.yAxisGroupLeft = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-y-left');
    this.yAxisGroupRight = svg.append('g').attr('class', 'axis').attr('uitestid', 'axis-y-right');
  }

  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  ngOnChanges(changes: SimpleChanges): void {
    this.draw();
  }

  public onMousemove(event: MouseEvent): void {
    console.log(document.elementsFromPoint(event.clientX, event.clientY));
  }

  @HostListener('window:resize', ['$event'])
  private onResize(): void {
    this.draw();
  }

  private draw(): void {
    const { width, height } = this.getAvailableSpace();
    const margin = this.margin;
    const svg = select(this.chartSelector);
    svg
      .attr('viewBox', `0 0 ${width} ${height}`)
      .attr('width', width)
      .attr('height', height)
      .attr('preserveAspectRatio', 'xMinYMin meet');

    this.xScale.domain([-0.5, this.getMaxValue((p) => p.x) + 0.5]).range([margin.left, width - margin.right]);
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
      this.drawLine(item, yScale, this.getColor(index));
    });
  }

  drawLine(item: ChartItem, yScale: LinearScale, color: string): void {
    const svg = select(this.chartSelector);
    const points: [number, number][] = item.points.map((d) => [d.x, d.y]);
    const path = svg
      .append('path')
      .datum(points)
      .attr('fill', 'none')
      .attr('stroke', color)
      .attr('stroke-width', 1.5)
      .attr('class', 'line red')
      .attr('uitestid', `line-${this.replaceSpace(item.identifier)}`)
      .attr(
        'd',
        line()
          .x((d) => this.xScale(d[0]))
          .y((d) => yScale(d[1]))
      );
    this.paths.push(path);
  }

  drawBar(item: ChartItem, height: number): void {
    const svg = select(this.chartSelector);
    const margin = this.margin;
    const points: [number, number, string][] = item.points.map((d) => [d.x, d.y, d.color ?? '#7e7e7e']);
    const selection = svg
      .selectAll()
      .data(points)
      .enter()
      .append('g')
      .attr('uitestid', `bar-${this.replaceSpace(item.identifier)}`);

    const rect = selection
      .append('rect')
      .attr('x', (d) => this.xScale(d[0]) - 1)
      .attr('y', (d) => this.yScaleLeft(d[1]))
      .attr('width', 3)
      .attr('height', (d) => height - margin.bottom - this.yScaleLeft(d[1]))
      .attr('fill', (d) => d[2])
      .attr('uitestid', (d) => `bar-${this.replaceSpace(item.identifier)}-${d[0]}`);
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

  public getAvailableSpace(): { width: number; height: number } {
    const parentElement: HTMLElement = this.elementRef.nativeElement;
    const availableSpace = parentElement.getBoundingClientRect();
    return { width: availableSpace.width, height: _height };
  }

  private getMaxValue(fn: (p: ChartItemPoint) => number): number {
    return this.chartItems
      .filter((i) => !(i.invisible === true))
      .reduce((prev, cur) => Math.max(prev, ...cur.points.map(fn)), 0);
  }

  public getColor(index: number): string {
    return colors[index % colors.length];
  }

  public replaceSpace(value: string): string {
    return value.replace(/ /g, '-');
  }

  public getIconStyle(index: number, invisible?: boolean): string {
    if (invisible === true) {
      return `--dt-icon-color: #cccccc`;
    }
    return `--dt-icon-color: ${this.getColor(index)}`;
  }

  private get margin(): Margin {
    // const zoom = 1 / window.devicePixelRatio;
    return {
      top: _margin.top, // * zoom,
      right: _margin.right, // * zoom,
      bottom: _margin.bottom, // * zoom,
      left: _margin.left, // * zoom,
    };
  }
}
