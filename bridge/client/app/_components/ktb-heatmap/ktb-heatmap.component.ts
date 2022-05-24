/* eslint-disable @typescript-eslint/no-this-alias */
import {
  ChangeDetectionStrategy,
  ChangeDetectorRef,
  Component,
  ElementRef,
  EventEmitter,
  HostListener,
  Inject,
  Input,
  OnDestroy,
  Output,
  ViewChild,
} from '@angular/core';
import { axisBottom, axisLeft, BaseType, ScaleBand, scaleBand, select, Selection, ValueFn } from 'd3';
import { ResultTypes } from '../../../../shared/models/result-types';
import { DtButton } from '@dynatrace/barista-components/button';
import { v4 as uuid } from 'uuid';
import { KtbHeatmapTooltipComponent } from './ktb-heatmap-tooltip.component';
import {
  EvaluationResultType,
  EvaluationResultTypeExtension,
  GroupedDataPoints,
  IDataPoint,
} from '../../_interfaces/heatmap';
import { DOCUMENT } from '@angular/common';
import {
  calculateTooltipPosition,
  createGroupedDataPoints,
  findDataPointThroughIdentifier,
  getAvailableIdentifiers,
  getAxisElements,
  getDataPointElement,
  getHiddenYElements,
  getLimitedYElements,
  getXAxisReducedElements,
  getYAxisElements,
  isScrollbarVisible,
} from './ktb-heatmap-utils';

type SVGGSelection = Selection<SVGGElement, unknown, HTMLElement, unknown>;
type HighlightSelection = Selection<SVGRectElement, unknown, HTMLElement, unknown>;
type SecondaryHighlightSelections = Selection<SVGRectElement, unknown, SVGGElement, unknown>;
type HeatmapTiles = Selection<SVGRectElement | null, IDataPoint, SVGGElement, unknown>;

@Component({
  selector: 'ktb-heatmap',
  templateUrl: './ktb-heatmap.component.html',
  styleUrls: ['./ktb-heatmap.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbHeatmapComponent implements OnDestroy {
  public readonly uniqueId = `heatmap-${uuid()}`;
  private readonly chartSelector = `#${this.uniqueId}`;
  private readonly heatmapSelector = `${this.chartSelector} .heatmap-container`;
  private readonly svgSelector = `${this.chartSelector}>svg`;
  private readonly mouseMoveListener: (this: Document, _evt: MouseEvent) => void;
  private readonly firstYElementPadding = 6; // "score" will then be 6px smaller than the rest.
  private readonly legendPadding = 10; // padding between xAxis and legend
  private readonly maxYAxisLabelWidth = 150;
  private readonly heightPerYElement = 30;
  private readonly limitYElementCount = 10;
  private readonly minWidthPerXAxisElement = 25;
  private readonly showMoreButtonPadding = 6;
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
  private xAxis: ScaleBand<string> = scaleBand();
  private yAxis: ScaleBand<string> = scaleBand();
  private dataPointContentWidth = 0;
  private height = 0;
  private _selectedDataPoint?: IDataPoint;
  // selectedIdentifier may be an invalid one, but must still be set because it could be set before the dataSource is set
  private _selectedIdentifier?: string;
  private mouseCoordinates = { x: 0, y: 0 };
  private groupedData: GroupedDataPoints = {};
  private yElements: string[] = [];
  public showMoreVisible = false;
  public showMoreExpanded = false;

  @ViewChild('showMoreButton', { static: false }) showMoreButton!: DtButton;
  @ViewChild('tooltip', { static: false }) tooltip!: KtbHeatmapTooltipComponent;
  @Output() selectedIdentifierChange = new EventEmitter<string>();

  @Input()
  public set dataPoints(data: IDataPoint[]) {
    this.removeHeatmap();
    this.groupedData = createGroupedDataPoints(data);
    this.yElements = getYAxisElements(this.groupedData);
    this.createHeatmap(this.groupedData);
    this.onResize(); // generating the heatmap may introduce a scrollbar
    this.selectedIdentifier = this._selectedIdentifier; // restore previously selected dataPoint
  }

  @Input()
  public set selectedIdentifier(identifier: string | undefined) {
    this._selectedIdentifier = identifier;
    const dataPoint = identifier ? findDataPointThroughIdentifier(identifier, this.groupedData) : undefined;
    this.click(dataPoint, true);
  }
  public get selectedIdentifier(): string | undefined {
    return this._selectedIdentifier;
  }

  private get showMoreButtonHeight(): number {
    const element: HTMLElement = this.showMoreButton._elementRef.nativeElement;
    return element.getBoundingClientRect().height + this.showMoreButtonPadding;
  }

  private get heatmapInstance(): SVGGSelection {
    return select(this.heatmapSelector);
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

  constructor(
    private elementRef: ElementRef,
    @Inject(DOCUMENT) private document: Document,
    private _changeDetectorRef: ChangeDetectorRef
  ) {
    // has to be globally instead of component bound, else scrolling into it will not have any mouse coordinates
    this.mouseMoveListener = (event: MouseEvent): void => this.onMouseMove(event);
    this.document.addEventListener('mousemove', this.mouseMoveListener);
  }

  private onMouseMove(event: MouseEvent): void {
    this.mouseCoordinates = {
      x: event.clientX * window.devicePixelRatio, // coordinates may stay and zoom-level could change. Normalize the coordinates.
      y: event.clientY * window.devicePixelRatio,
    };
  }

  @HostListener('window:resize', ['$event'])
  private onResize(): void {
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
    const element = getDataPointElement(x, y);
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

  private getDataPointThroughCoordinates(x: number, y: number): IDataPoint | undefined {
    const element = getDataPointElement(x, y);

    if (!element || !this.heatmapInstance.node()?.contains(element)) {
      return;
    }

    return select(element)?.datum() as IDataPoint | undefined;
  }

  private removeHeatmap(): void {
    select(this.svgSelector).remove();
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
    const { yElements, xElements, showMoreVisible } = getAxisElements(data, this.limitYElementCount);
    this.showMoreVisible = showMoreVisible;
    this._changeDetectorRef.detectChanges(); // update visibility of button

    this.setHeight(yElements.length);
    const availableSpace = this.setAndGetAvailableSpace();

    const svg = select(this.chartSelector).append('svg').attr('preserveAspectRatio', 'xMinYMin meet');
    this.resizeSvg(availableSpace.width, availableSpace.height);

    svg.append('g').classed('heatmap-container', true).attr('transform', `translate(${this.maxYAxisLabelWidth}, 0)`);

    this.setData(data, xElements, yElements);
    this.createLegend();
    this.resizeDataPointContainerRect();
  }

  private setAndGetAvailableSpace(): { height: number; width: number } {
    const parentElement: HTMLElement = this.elementRef.nativeElement.parentNode;
    const availableSpace = parentElement.getBoundingClientRect();
    const width = availableSpace.width;
    const xAxisHeight = this.xAxisContainer.node()?.getBoundingClientRect().height ?? 0;
    const legendHeight = this.legendContainer.node()?.getBoundingClientRect().height ?? 0;
    const height =
      this.height +
      xAxisHeight +
      legendHeight +
      this.legendPadding +
      (this.showMoreVisible ? this.showMoreButtonHeight : 0) +
      10; //padding-bottom
    this.dataPointContentWidth = width - this.maxYAxisLabelWidth;

    return {
      height,
      width,
    };
  }

  private resizeSvg(width: number, height: number): void {
    select(this.svgSelector).attr('viewBox', `0 0 ${width} ${height}`).attr('width', width).attr('height', height);
  }

  private resizeXAxis(): void {
    this.xAxis = this.xAxis.range([0, this.dataPointContentWidth]);
    const xAxisContainer = this.xAxisContainer;

    this.setXAxisCoordinates(xAxisContainer);
    this.attachXAxis(xAxisContainer, this.xAxis);
  }

  private resizeDataPoints(): void {
    this.setDataPointCoordinates(this.dataPointElements);
  }

  private resizeHighlights(): void {
    if (!this._selectedDataPoint) {
      return;
    }

    this.setHighlightCoordinates(this._selectedDataPoint.xElement);
    this.setSecondaryHighlightCoordinates(this._selectedDataPoint.comparedIdentifier);
  }

  private resizeShowMoreButton(): void {
    if (this.showMoreVisible) {
      const htmlElement: HTMLElement = this.showMoreButton._elementRef.nativeElement;

      htmlElement.style.top = `${this.height + this.showMoreButtonPadding / 2}px`;
      htmlElement.style.left = `${this.maxYAxisLabelWidth}px`;
      htmlElement.style.width = `${this.dataPointContentWidth}px`;
    }
  }

  private resizeLegend(): void {
    const legend = this.legendContainer;
    const fullLength = legend.node()?.getBoundingClientRect().width ?? 0;
    const centerXPosition = (this.dataPointContentWidth - fullLength) / 2;
    const xAxisHeight = this.xAxisContainer.node()?.getBoundingClientRect().height ?? 0;
    const yPosition =
      this.height + this.legendPadding + xAxisHeight + (this.showMoreVisible ? this.showMoreButtonHeight : 0);
    legend.attr('transform', `translate(${centerXPosition}, ${yPosition})`);
  }

  private setData(data: GroupedDataPoints, xAxisElements: string[], yAxisElements: string[]): void {
    const heatmap = this.heatmapInstance;
    this.xAxis = this.addXAxis(heatmap, xAxisElements);
    this.yAxis = this.addYAxis(heatmap, yAxisElements);
    this.generateHeatmapTiles(data);
  }

  private addXAxis(heatmap: SVGGSelection, xElements: string[]): ScaleBand<string> {
    const x = this.xAxis.range([0, this.dataPointContentWidth]).domain(xElements);
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
      .call(
        axisBottom(x)
          .tickSize(5)
          .tickValues(getXAxisReducedElements(x.domain(), this.dataPointContentWidth, this.minWidthPerXAxisElement))
      )
      .selectAll('text')
      .attr('class', 'x-axis-identifier')
      .attr('dx', '-.8em')
      .attr('dy', '.15em');
  }

  private addYAxis(heatmap: SVGGSelection, yElements: string[]): ScaleBand<string> {
    const y = this.yAxis.range([this.height, 0]).domain(yElements);
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
    const yAxis = yAxisContainer.call(axisLeft(y).tickSize(0));
    yAxis.selectAll('.tick').each(this.setEllipsisStyle(this.maxYAxisLabelWidth));
    yAxis.select('.domain').remove();
  }

  private setEllipsisStyle(labelWidth: number): ValueFn<BaseType, unknown, void> {
    return function (this: BaseType): void {
      const self = select(this as SVGGElement);
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
    const scrollbarWidth = isScrollbarVisible() ? 18 : 0; // just assume a default scrollbar-width of 18px
    const { top, left } = calculateTooltipPosition(tooltipWidth, scrollbarWidth, event.x, event.y);

    htmlElement.style.top = `${top}px`;
    htmlElement.style.left = `${left}px`;
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

    if (!dataPoint) {
      this._selectedDataPoint = undefined;
      return;
    }

    this._selectedDataPoint = dataPoint;
    this._selectedIdentifier = dataPoint.identifier;

    heatmap.append('rect').attr('class', 'highlight-primary');
    this.setHighlightCoordinates(dataPoint.xElement);

    const foundIdentifiers = getAvailableIdentifiers(dataPoint.comparedIdentifier, this.groupedData);
    heatmap.selectAll().data(foundIdentifiers).join('rect').attr('class', 'highlight-secondary');
    this.setSecondaryHighlightCoordinates(foundIdentifiers);

    if (!preSelectDataPoint) {
      this.selectedIdentifierChange.emit(dataPoint.identifier);
    }
  }

  /**
   * For the special case that the user clicks on an dataPoint that does not exist (another dataPoint in the column exists)
   * @param event$
   * @param element
   */
  private contentClick(event$: MouseEvent, element: SVGRectElement): void {
    const containerY = element.getBoundingClientRect().top;
    const dataPoint = this.getDataPointThroughCoordinates(event$.x, containerY + 5); // offset to make sure to click on the tile
    if (!dataPoint) {
      return;
    }
    this.click(dataPoint);
  }

  private getHighlightWidth(): number {
    const xAxisElements = this.xAxis.domain();
    return this.dataPointContentWidth / xAxisElements.length;
  }

  private setHighlightCoordinates(identifier: string): void {
    this.highlight
      .attr('x', this.xAxis(identifier) ?? null)
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private setSecondaryHighlightCoordinates(identifiers: string[]): void {
    this.secondaryHighlights
      .attr('x', (_dt, index) => {
        const xElement = findDataPointThroughIdentifier(identifiers[index], this.groupedData)?.xElement;
        if (!xElement) {
          return null;
        }
        return this.xAxis(xElement) ?? null;
      })
      .attr('y', 0)
      .attr('height', this.height)
      .attr('width', this.getHighlightWidth());
  }

  private generateHeatmapTiles(data: GroupedDataPoints, yAxisElements = this.yAxis.domain()): void {
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
      .attr('uitestid', (yElement) => yElement.replace(/ /g, '-'))
      .selectAll()
      .data((key) => data[key])
      .join('rect')
      .attr('class', (dataPoint) => dataPoint.color)
      .classed('data-point', true)
      // set all new dataPoints (show all yElements) to disabled if needed
      .classed('disabled', (dataPoint: IDataPoint) => this.legendDisabledStatus[dataPoint.color])
      .attr('uitestid', (dataPoint) => `ktb-heatmap-tile-${dataPoint.identifier.replace(/ /g, '-')}`)
      .on('click', (_event: PointerEvent, dataPoint: IDataPoint) => this.click(dataPoint))
      .on('mouseover', function (this: SVGGElement | null) {
        if (!this) {
          return;
        }
        _this.mouseOver(this);
      })
      .on('mousemove', (event: MouseEvent, dataPoint: IDataPoint) => this.mouseMove(event, dataPoint))
      .on('mouseleave', () => this.mouseLeave());

    this.setDataPointCoordinates(dataPoints);
  }

  private resizeDataPointContainerRect(): void {
    this.dataPointContainerRect.attr('width', this.dataPointContentWidth).attr('height', this.height);
  }

  private setDataPointCoordinates(dataPoints: HeatmapTiles): void {
    const yAxisElements = this.yAxis.domain();
    const firstYElement = yAxisElements[yAxisElements.length - 1];
    dataPoints
      .attr('x', (dataPoint) => this.xAxis(dataPoint.xElement) ?? null)
      .attr('y', (dataPoint) => {
        const yCoordinate = this.yAxis(dataPoint.yElement);
        if (yCoordinate !== undefined && dataPoint.yElement === firstYElement) {
          return yCoordinate + this.firstYElementPadding / 2;
        }
        return yCoordinate ?? null;
      })
      .attr('width', this.xAxis.bandwidth())
      .attr('height', (dataPoint) => {
        const height = this.yAxis.bandwidth();
        if (dataPoint.yElement === firstYElement) {
          return height - this.firstYElementPadding;
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
        select(this).classed('disabled', isDisabled);
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
    this.setHeight(this.yElements.length);
    this.updateYAxis(this.yElements);

    this.generateHeatmapTiles(this.groupedData, getHiddenYElements(this.yElements, this.limitYElementCount));
  }

  private updateYAxis(yElements: string[]): void {
    this.yAxis = this.yAxis.range([this.height, 0]).domain(yElements);
    this.attachYAxis(this.yAxisContainer, this.yAxis);
  }

  private collapseHeatmap(): void {
    this.setHeight(this.limitYElementCount);
    this.dataPointContainer
      .selectAll('g')
      .filter((_element, index) => index >= this.limitYElementCount)
      .remove();

    this.updateYAxis(getLimitedYElements(this.yElements, this.limitYElementCount));
  }

  private setHeight(elementCount: number): void {
    this.height = elementCount * this.heightPerYElement;
  }

  public ngOnDestroy(): void {
    this.document.removeEventListener('mousemove', this.mouseMoveListener);
  }
}
/* eslint-enable @typescript-eslint/no-this-alias */
