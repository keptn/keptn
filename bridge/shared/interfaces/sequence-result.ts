import { Sequence } from '../models/sequence';

export interface SequenceResult {
  nextPageKey?: number;
  pageSize?: number;
  totalCount?: number;
  states: Sequence[];
}
