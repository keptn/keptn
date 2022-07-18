import { Component, OnInit } from '@angular/core';
import {
  EvaluationResultType,
  EvaluationResultTypeExtension,
  IDataPoint,
  IHeatmapTooltipType,
} from '../../../_interfaces/heatmap';
import moment from 'moment';
import { ResultTypes } from '../../../../../shared/models/result-types';

function getColor(value: number): EvaluationResultType {
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

export function generateTestData(sliCounter: number, counter: number): IDataPoint[] {
  const categories = ['score'];
  const slis = [];
  for (let i = 0; i < sliCounter - 1; ++i) {
    slis.push(`response time p${i}`);
  }
  slis.push(`response time p${sliCounter - 1} very long SLI name here`);
  categories.push(...slis);
  const data: IDataPoint[] = [];
  const dateMillis = new Date().getTime();

  data.push({
    xElement: moment(new Date(dateMillis - 1000 * 60)).format('YYYY-MM-DD HH:mm'),
    yElement: 'score',
    color: getColor(-1 % 4),
    tooltip: {
      type: IHeatmapTooltipType.SCORE,
      value: -1,
      fail: -1 % 2 === 1,
      failedCount: -1 + 4,
      warn: -1 % 2 === 0,
      passCount: -1,
      thresholdPass: -1 + 1,
      thresholdWarn: -1 + 2,
      warningCount: -1 + 3,
    },
    identifier: `keptnContext_${-1}`,
    comparedIdentifier: [],
  });

  // adding one duplicate (two evaluations have the same time)
  let y = -1;
  for (const category of categories) {
    data.push({
      xElement: moment(new Date(dateMillis)).format('YYYY-MM-DD HH:mm'),
      yElement: category,
      color: getColor(y % 4),
      tooltip: {
        type: IHeatmapTooltipType.SLI,
        value: y,
        keySli: y % 2 === 1,
        score: y,
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
    ++y;
  }

  // fill SLIs with random data (-1 to have an evaluation with "missing" data)
  let offset = 0;
  for (let i = 0; i < counter - 1; ++i) {
    data.push({
      xElement: moment(new Date(dateMillis + i * 1000 * 60)).format('YYYY-MM-DD HH:mm'),
      yElement: 'score',
      color: getColor(i % 4),
      tooltip: {
        type: IHeatmapTooltipType.SCORE,
        value: i,
        fail: i % 2 === 1,
        failedCount: i + 4,
        warn: i % 2 === 0,
        passCount: i,
        thresholdPass: i + 1,
        thresholdWarn: i + 2,
        warningCount: i + 3,
      },
      identifier: `keptnContext_${i}`,
      comparedIdentifier: [`keptnContext_${i - 1}`, `keptnContext_${i - 2}`],
    });
    offset++;
    for (const category of slis) {
      data.push({
        xElement: moment(new Date(dateMillis + i * 1000 * 60)).format('YYYY-MM-DD HH:mm'),
        yElement: category,
        color: getColor((i + offset) % 4),
        tooltip: {
          type: IHeatmapTooltipType.SLI,
          value: i,
          keySli: i % 2 === 1,
          score: i,
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
      ++offset;
    }
    offset = 0;
  }
  return data;
}

@Component({
  selector: 'ktb-test-heatmap',
  templateUrl: './ktb-test-heatmap.component.html',
  styleUrls: ['./ktb-test-heatmap.component.scss'],
})
export class KtbTestHeatmapComponent implements OnInit {
  public dataPoints: IDataPoint[] = [];
  public selectedIdentifier?: string;
  private sliCount = 12;
  private evaluationCount = 50;

  public ngOnInit(): void {
    this.setDataPoints();
  }

  public setDataPoints(): void {
    this.selectedIdentifier = 'keptnContext_2';
    this.dataPoints = generateTestData(this.sliCount, this.evaluationCount);
  }
}
