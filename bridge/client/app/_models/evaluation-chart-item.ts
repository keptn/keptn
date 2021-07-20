import { Trace } from './trace';
import { IndicatorResult } from './indicator-result';

export interface EvaluationChartItem {
  metricName: string;
  name: string;
  type: string;
  data: {
    y: number,
    x?: number,
    indicatorResult?: IndicatorResult,
    evaluationData?: Trace,
    label?: string,
    name: string,
    color?: string
  } [];
  turboThreshold: number;
  visible?: boolean;
  yAxis?: number;
  cursor?: string;
}
