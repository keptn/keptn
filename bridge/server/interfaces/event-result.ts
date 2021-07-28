import { Trace } from '../models/trace.js';

export interface EventResult {
  events: Trace[];
  totalCount: number;
}
