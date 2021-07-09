import {ResultTypes} from './result-types';

export class SliResult {
  public comparedValue?: number;
  public absoluteChange?: number;
  public relativeChange?: number;
  public name: string;
  public value: string | number;
  public result: ResultTypes;
  public score: number;
  public passTargets?: {criteria: string, targetValue: number, violated: boolean}[];
  public warningTargets?: {criteria: string, targetValue: number, violated: boolean}[];
  public targets?: {criteria: string, targetValue: number, violated: boolean}[];
  public keySli: boolean;
  public success: boolean;
  public expanded = false;
}
