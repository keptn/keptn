import { Component } from '@angular/core';
import { ChartItem } from '../../../_interfaces/chart';

@Component({
  selector: 'ktb-test-chart',
  templateUrl: './ktb-test-chart.component.html',
  styleUrls: ['./ktb-test-chart.component.scss'],
})
export class KtbTestChartComponent {
  public labels = { 1: '2022-02-22 09:22', 2: '2022-02-22 12:03' };
  public tooltipLabels = {
    0: 'SLO evaluation of test from 2022-02-22 07:01',
    1: '2022-02-22 09:22',
    2: '2022-02-22 12:03',
  };
  public data: ChartItem[] = [
    {
      type: 'score-bar',
      identifier: 'Score',
      points: [
        {
          x: 0,
          y: 33,
          color: '#e30505',
          identifier: 'e0',
        },
        {
          x: 1,
          y: 66,
          color: '#ffaa00',
          identifier: 'e1',
        },
        {
          x: 2,
          y: 88,
          color: '#518637',
          identifier: 'e2',
        },
        {
          x: 3,
          y: 50,
          color: '#a97223',
          identifier: 'e3',
        },
        {
          x: 4,
          y: 2,
          color: '#f30010',
          identifier: 'e4',
        },
      ],
    },
    {
      type: 'score-line',
      identifier: 'Score',
      points: [
        {
          x: 0,
          y: 33,
          identifier: 'e0',
        },
        {
          x: 1,
          y: 66,
          identifier: 'e1',
        },
        {
          x: 2,
          y: 88,
          identifier: 'e2',
        },
        {
          x: 3,
          y: 50,
          identifier: 'e3',
        },
        {
          x: 4,
          y: 2,
          identifier: 'e4',
        },
      ],
    },
    {
      type: 'metric-line',
      identifier: 'Metric 1',
      label: 'My custom metric 1 label',
      points: [
        {
          x: 1,
          y: 30,
          identifier: 'e1',
        },
        {
          x: 2,
          y: 40,
          identifier: 'e2',
        },
        {
          x: 3,
          y: 144.5,
          identifier: 'e3',
        },
        {
          x: 4,
          y: 10,
          identifier: 'e4',
        },
      ],
    },
    {
      type: 'metric-line',
      identifier: 'Metric 2',
      invisible: true,
      points: [
        {
          x: 1,
          y: 4,
          identifier: 'e1',
        },
        {
          x: 2,
          y: 5,
          identifier: 'e2',
        },
        {
          x: 3,
          y: 12,
          identifier: 'e3',
        },
        {
          x: 4,
          y: 12,
          identifier: 'e4',
        },
      ],
    },
  ];

  public addMetric(): void {
    const name = 'Metric ' + this.randomNumber(100, 999);
    const item: ChartItem = {
      type: 'metric-line',
      identifier: name,
      invisible: false,
      points: [
        {
          x: 0,
          y: this.randomNumber(0, 255),
          identifier: 'e0',
        },
        {
          x: 1,
          y: this.randomNumber(0, 255),
          identifier: 'e1',
        },
        {
          x: 2,
          y: this.randomNumber(0, 255),
          identifier: 'e2',
        },
        {
          x: 3,
          y: this.randomNumber(0, 255),
          identifier: 'e3',
        },
        {
          x: 4,
          y: this.randomNumber(0, 255),
          identifier: 'e4',
        },
      ],
    };
    this.data = [...this.data, item];
  }

  private randomNumber(min: number, max: number): number {
    return Math.floor(Math.random() * (max - min)) + min + 1;
  }
}
