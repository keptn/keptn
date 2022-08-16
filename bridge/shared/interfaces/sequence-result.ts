import { ISequenceState } from './sequence';

export interface SequenceResult {
  nextPageKey?: number;
  pageSize?: number;
  totalCount?: number;
  states: ISequenceState[];
}
