import { Trace } from './trace';
import { IndicatorResult } from './indicator-result';

export interface EvaluationChartDataItem {
  y: number;
  x?: number;
  indicatorResult?: IndicatorResult;
  evaluationData?: Trace;
  label?: string;
  name: string;
  color?: string;
}

export interface EvaluationChartItem {
  metricName: string;
  name: string;
  type: string;
  data: EvaluationChartDataItem [];
  turboThreshold: number;
  visible?: boolean;
  yAxis?: number;
  cursor?: string;
}
