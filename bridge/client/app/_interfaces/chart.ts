import { Trace } from '../_models/trace';
import { IndicatorResult } from '../../../shared/interfaces/indicator-result';

export type DrawType = 'metric-line' | 'score-bar' | 'score-line';

export interface ChartItemPoint {
  x: number;
  y: number;
  color?: string;
}

export interface ChartItem {
  label: string;
  type: DrawType;
  invisible?: boolean;
  points: ChartItemPoint[];
}

export type ChartItemPointInfo = Record<string, { points: ChartItemPoint[]; label?: string } | undefined>;
export type FuncEvaluationToChartItemPoint = (evaluation: Trace, index: number) => ChartItemPoint;
export type FuncMetricToChartItem = (metric: string) => ChartItem;
export type FuncMapIndicatorResult = (indicatorResult: IndicatorResult) => void;
export type FuncDateToString = (date: string, index: number) => string;
export type FuncDateToDict = (date: string) => number | undefined;
