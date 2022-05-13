import { Component, OnInit } from '@angular/core';
import {
  EvaluationResultType,
  EvaluationResultTypeExtension,
  IDataPoint,
  IHeatmapTooltipType,
} from '../../../_interfaces/heatmap';
import moment from 'moment';
import { ResultTypes } from '../../../../../shared/models/result-types';

@Component({
  selector: 'ktb-test-heatmap',
  templateUrl: './ktb-test-heatmap.component.html',
  styleUrls: ['./ktb-test-heatmap.component.scss'],
})
export class KtbTestHeatmapComponent implements OnInit {
  public dataPoints: IDataPoint[] = [];
  private sliCount = 12;
  private evaluationCount = 50;

  public ngOnInit(): void {
    this.setDataPoints();
  }

  public setDataPoints(): void {
    this.dataPoints = this.generateTestData(this.sliCount, this.evaluationCount);
  }

  private generateTestData(sliCounter: number, counter: number): IDataPoint[] {
    const categories = [];
    for (let i = 0; i < sliCounter - 1; ++i) {
      categories.push(`response time p${i}`);
    }
    categories.push(`response time p${sliCounter - 1} very long SLI name here`);
    const data: IDataPoint[] = [];
    const dateMillis = new Date().getTime();

    // adding one duplicate (two evaluations have the same time)
    for (const category of [...categories, 'score']) {
      data.push({
        xElement: moment(new Date(dateMillis)).format('YYYY-MM-DD HH:mm'),
        yElement: category,
        color: this.getColor(Math.floor(Math.random() * 4)),
        tooltip: {
          type: IHeatmapTooltipType.SLI,
          value: Math.random(),
          keySli: Math.floor(Math.random() * 2) === 1,
          score: Math.floor(Math.random() * 100),
          passTargets: [
            {
              targetValue: 0,
              criteria: '<=1',
              violated: true,
            },
          ],
          warningTargets: [
            {
              targetValue: 0,
              violated: false,
              criteria: '<=10',
            },
          ],
        },
        identifier: `keptnContext_${-1}`,
        comparedIdentifier: [],
      });
    }

    // fill SLIs with random data (-1 to have an evaluation with "missing" data)
    for (const category of categories) {
      for (let i = 0; i < counter - 1; ++i) {
        data.push({
          xElement: moment(new Date(dateMillis + i * 1000 * 60)).format('YYYY-MM-DD HH:mm'),
          yElement: category,
          color: this.getColor(Math.floor(Math.random() * 4)),
          tooltip: {
            type: IHeatmapTooltipType.SLI,
            value: Math.random(),
            keySli: Math.floor(Math.random() * 2) === 1,
            score: Math.floor(Math.random() * 100),
            passTargets: [
              {
                targetValue: 0,
                criteria: '<=1',
                violated: true,
              },
            ],
            warningTargets: [
              {
                targetValue: 0,
                violated: false,
                criteria: '<=10',
              },
            ],
          },
          identifier: `keptnContext_${i}`,
          comparedIdentifier: [`keptnContext_${i - 1}`, `keptnContext_${i - 2}`],
        });
      }
    }
    categories.push('score');
    for (let i = 0; i < counter; ++i) {
      data.push({
        xElement: moment(new Date(dateMillis + i * 1000 * 60)).format('YYYY-MM-DD HH:mm'),
        yElement: 'score',
        color: this.getColor(Math.floor(Math.random() * 4)),
        tooltip: {
          type: IHeatmapTooltipType.SCORE,
          value: Math.random(),
          fail: Math.floor(Math.random() * 2) === 1,
          failedCount: Math.random(),
          warn: Math.floor(Math.random() * 2) === 1,
          passCount: Math.random(),
          thresholdPass: Math.random(),
          thresholdWarn: Math.random(),
          warningCount: Math.random(),
        },
        identifier: `keptnContext_${i}`,
        comparedIdentifier: [`keptnContext_${i - 1}`, `keptnContext_${i - 2}`],
      });
    }
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
}
