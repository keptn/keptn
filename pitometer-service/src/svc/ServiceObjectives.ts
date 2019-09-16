export interface ServiceObjectives {
  pass: number;
  warning: number;
  objectives: Objective[];
}

export interface Objective {
  name: string;
  threshold: number;
  timeframe: string;
  score: number;
}
