import { Component } from '@angular/core';
import { ChartItem } from '../../../_interfaces/chart';
import * as testData from './ktb-chart-test-data';

@Component({
  selector: 'ktb-test-chart',
  templateUrl: './ktb-test-chart.component.html',
  styleUrls: ['./ktb-test-chart.component.scss'],
})
export class KtbTestChartComponent {
  labels = testData.labels;
  tooltipLabels = testData.tooltipLabels;
  data = testData.data;

  private currentIndex = 0;

  public addMetric(): void {
    const name = 'Metric ' + this.getSomeNumber(11);
    const item: ChartItem = {
      type: 'metric-line',
      label: name,
      invisible: false,
      points: [
        {
          x: 0,
          y: this.getSomeNumber(1),
        },
        {
          x: 1,
          y: this.getSomeNumber(2),
        },
        {
          x: 2,
          y: this.getSomeNumber(3),
        },
        {
          x: 3,
          y: this.getSomeNumber(3),
        },
        {
          x: 4,
          y: this.getSomeNumber(4),
        },
      ],
    };
    this.data = [...this.data, item];
  }

  private getSomeNumber(addTo: number): number {
    return this.someValues[this.currentIndex++ % this.someValues.length] + addTo + this.currentIndex;
  }

  private readonly someValues = [33, 63, 6, 9, 8, 101, 34, 64, 21, 66, 54, 11, 18, 21, 34, 4, 2, 1, 32, 145, 59];
}
