export interface ChartItemPoint {
  x: number;
  y: number;
  identifier: string; //TODO: why do we need this?
  color?: string;
}

export interface ChartItem {
  identifier: string;
  label?: string;
  type: 'metric-line' | 'score-bar' | 'score-line';
  invisible?: boolean;
  points: ChartItemPoint[];
}

export type IChartItemPointInfo = Record<string, { points: ChartItemPoint[]; label?: string } | undefined>;
