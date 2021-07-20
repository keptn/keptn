import { ResultTypes } from './result-types';

export type Target = {criteria: string, targetValue: number, violated: boolean};

export interface IndicatorResult {
  value: {
    value: number,
    metric: string,
    success: boolean,
    message: string
  };
  score: number;
  displayName?: string;
  status: ResultTypes;
  passTargets?: Target[];
  warningTargets?: Target[];
  targets?: Target[];
  keySli: boolean;
}
