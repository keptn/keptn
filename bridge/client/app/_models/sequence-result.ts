import {Sequence} from './sequence';

export class SequenceResult {
  nextPageKey: number;
  pageSize: number;
  totalCount: number;
  states: Sequence[];
}
