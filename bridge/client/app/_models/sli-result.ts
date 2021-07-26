import {ResultTypes} from './result-types';
import { Target } from './indicator-result';

export interface SliResult {
  comparedValue?: number;
  absoluteChange?: number;
  relativeChange?: number;
  name: string;
  value: string | number;
  result: ResultTypes;
  score: number;
  passTargets?: Target[];
  warningTargets?: Target[];
  targets?: Target[];
  keySli: boolean;
  success: boolean;
  expanded: boolean;
}
