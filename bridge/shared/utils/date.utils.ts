import { Trace } from '../models/trace';

export class DateUtil {
  static compareTraceTimesAsc<T extends Trace>(a: T, b: T): number {
    return DateUtil.compareTraceTimesDesc(a, b, -1);
  }

  static compareTraceTimesDesc<T extends Trace>(a?: T, b?: T, direction = 1): number {
    let result;
    if (a?.time && b?.time) {
      result = new Date(a.time).getTime() - new Date(b.time).getTime();
    } else if (a?.time && !b?.time) {
      result = 1;
    } else if (!a?.time && b?.time) {
      result = -1;
    } else {
      result = 0;
    }
    return result * direction;
  }
}
