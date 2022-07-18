import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbChartComponent } from './ktb-chart.component';
import { KtbChartModule } from './ktb-chart.module';
import * as testData from './testing/ktb-chart-test-data';
import { DOCUMENT } from '@angular/common';
import { ElementRef } from '@angular/core';
import { firstValueFrom } from 'rxjs';
import Mock = jest.Mock;

describe('KtbChartComponent', () => {
  let component: KtbChartComponent;
  let fixture: ComponentFixture<KtbChartComponent>;
  const parentNodeBoundingClientRectSpy: Mock<DOMRect, [void]> = jest.fn();
  const elementsFromPointSpy: Mock<Element[], [number, number]> = jest.fn();

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbChartModule, HttpClientTestingModule],
      providers: [],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbChartComponent);
    component = fixture.componentInstance;
    mockUIElements();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize the chart', () => {
    // given, when
    fixture.detectChanges();

    // then
    expect(Reflect.get(component, 'xScale')).toBeTruthy();
    expect(Reflect.get(component, 'xScale').domain()).toEqual([-0.5, 0.5]);
    expect(Reflect.get(component, 'xScale').range()).toEqual([50, 940]);
    expect(Reflect.get(component, 'yScaleLeft')).toBeTruthy();
    expect(Reflect.get(component, 'yScaleLeft').domain()).toEqual([0, 100]);
    expect(Reflect.get(component, 'yScaleLeft').range()).toEqual([300, 60]);
    expect(Reflect.get(component, 'yScaleRight')).toBeTruthy();
    expect(Reflect.get(component, 'yScaleRight').domain()).toEqual([0, 0]);
    expect(Reflect.get(component, 'yScaleRight').range()).toEqual([300, 60]);
    expect(Reflect.get(component, 'xAxisGroup')).toBeTruthy();
    expect(Reflect.get(component, 'yAxisGroupLeft')).toBeTruthy();
    expect(Reflect.get(component, 'yAxisGroupRight')).toBeTruthy();
    expect(Reflect.get(component, 'paths').length).toBe(0);
    expect(Reflect.get(component, 'rects').length).toBe(1);
  });

  it('should draw lines and bars', () => {
    // given
    setTestData();

    // when
    fixture.detectChanges();

    // then
    expect(Reflect.get(component, 'xScale').domain()).toEqual([-0.5, 4.5]);
    expect(Reflect.get(component, 'yScaleLeft').domain()).toEqual([0, 100]);
    expect(Reflect.get(component, 'yScaleRight').domain()).toEqual([0, 144.5]);
    expect(Reflect.get(component, 'paths').length).toBe(2);
    expect(Reflect.get(component, 'rects').length).toBe(2);
  });

  it('should fetch the tooltip data', async () => {
    // given
    setTestData();
    fixture.detectChanges();
    component.tooltip = <ElementRef>{ nativeElement: { getBoundingClientRect: () => getDomRect(150, 170) } };

    // when
    component.onMousemove()(new MouseEvent('move', { clientX: 1, clientY: 1 }), [1, 1, 'transparent']);

    // then
    const tooltipState = await firstValueFrom(component.tooltipState$);
    expect(tooltipState).toEqual({
      label: '2022-02-22 09:22',
      left: 40,
      metricValues: [
        {
          label: 'Score',
          value: 66,
        },
        {
          label: 'My custom metric 1 label',
          value: 30,
        },
      ],
      top: 120,
      visible: false,
    });
  });

  it('should show the tooltip', async () => {
    // given
    fixture.detectChanges();

    // when
    component.showTooltip(true)();

    // then
    const tooltipState = await firstValueFrom(component.tooltipState$);
    expect(tooltipState.visible).toBe(true);
  });

  it('should hide the tooltip', async () => {
    // given
    fixture.detectChanges();

    // when
    component.showTooltip(false)();

    // then
    const tooltipState = await firstValueFrom(component.tooltipState$);
    expect(tooltipState.visible).toBe(false);
  });

  it('should hide a chart item', () => {
    // given
    setTestData();
    fixture.detectChanges();

    // when
    component.hideChartItem(component.chartItems[2]);

    // then
    expect(Reflect.get(component, 'paths').length).toBe(1);
  });

  function setTestData(): void {
    component.chartItems = testData.data;
    component.xLabels = testData.labels;
    component.xTooltipLabels = testData.tooltipLabels;
  }

  function mockUIElements(): void {
    const document = TestBed.inject(DOCUMENT);
    const aDiv = document.querySelector('div');
    const classList = aDiv?.classList;
    classList?.add('area');

    const element = <Element>(
      (<unknown>{ tagName: 'rect', classList: classList, getClientRects: () => [getDomRect(50, 80, 20, 40)] })
    );
    elementsFromPointSpy.mockReturnValue([element]);
    document.elementsFromPoint = elementsFromPointSpy;

    parentNodeBoundingClientRectSpy.mockReturnValue(getDomRect(1000, 800));
    const htmlElement: HTMLElement = fixture.nativeElement;
    jest.spyOn(htmlElement, 'getBoundingClientRect').mockImplementation(parentNodeBoundingClientRectSpy);
  }

  function getDomRect(width: number, height: number, top = 0, left = 0): DOMRect {
    return {
      width,
      y: 0,
      right: 0,
      bottom: 0,
      left,
      x: 0,
      height,
      top,
      toJSON(): object {
        return {};
      },
    };
  }
});
