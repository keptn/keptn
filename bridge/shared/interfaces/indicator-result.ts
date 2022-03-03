import { ResultTypes } from '../models/result-types';

export type Target = { criteria: string; targetValue: number; violated: boolean };

export interface IndicatorResult {
  value: {
    value: number;
    comparedValue?: number;
    metric: string;
    success: boolean;
    message?: string;
  };
  score: number;
  displayName?: string;
  status: ResultTypes;
  passTargets?: Target[] | null;
  warningTargets?: Target[] | null;
  targets?: Target[];
  keySli: boolean;
}
