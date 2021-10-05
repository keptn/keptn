import { Trace } from '../models/trace';

export interface EventResult {
  events: Trace[];
  totalCount: number;
  pageSize: number;
  nextPageKey: number;
}
