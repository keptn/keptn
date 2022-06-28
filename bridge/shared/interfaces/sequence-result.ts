import { ISequence } from './sequence';

export interface SequenceResult {
  nextPageKey?: number;
  pageSize?: number;
  totalCount?: number;
  states: ISequence[];
}
