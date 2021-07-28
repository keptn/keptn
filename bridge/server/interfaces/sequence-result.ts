import { Sequence } from '../models/sequence.js';

export interface SequenceResult {
  nextPageKey?: number;
  pageSize?: number;
  totalCount?: number;
  states: Sequence[];
}
