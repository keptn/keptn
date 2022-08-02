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
