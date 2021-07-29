import { EventResult as er } from '../../../shared/interfaces/event-result';
import { Trace } from '../_models/trace';

export interface EventResult extends er {
  events: Trace[];
}
