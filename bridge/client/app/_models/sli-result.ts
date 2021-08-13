import {ResultTypes} from '../../../shared/models/result-types';
import { Target } from '../../../shared/interfaces/indicator-result';

export class SliResult {
  comparedValue?: number;
  name!: string;
  value!: string | number;
  result!: ResultTypes;
  score!: number;
  passTargets?: Target[] | null;
  warningTargets?: Target[] | null;
  targets?: Target[];
  keySli!: boolean;
  success!: boolean;
  expanded!: boolean;
  calculatedChanges?: {
    absolute: number,
    relative: number
  };
  weight!: number;
}
