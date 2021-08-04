import {ResultTypes} from './result-types';
import { Target } from './indicator-result';

export class SliResult {
  comparedValue?: number;
  name!: string;
  value!: string | number;
  result!: ResultTypes;
  score!: number;
  passTargets?: Target[];
  warningTargets?: Target[];
  targets?: Target[];
  keySli!: boolean;
  success!: boolean;
  expanded!: boolean;
  calculatedChanges?: {
    absolute: number,
    relative: number
  };
}
