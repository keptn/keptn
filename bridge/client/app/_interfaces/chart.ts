export interface ChartItemPoint {
  x: number;
  y: number;
  identifier: string;
  color?: string;
}

export interface ChartItem {
  identifier: string;
  label?: string;
  type: 'metric-line' | 'score-bar' | 'score-line';
  invisible?: boolean;
  points: ChartItemPoint[];
}
