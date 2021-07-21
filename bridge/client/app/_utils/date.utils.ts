import moment from 'moment';
import {Trace} from '../_models/trace';
import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class DateUtil {

  public readonly DEFAULT_DATE_FORMAT = 'YYYY-MM-DD';
  public readonly DEFAULT_TIME_FORMAT = 'HH:mm';

  static compareTraceTimesAsc(a: Trace, b: Trace): number {
    return DateUtil.compareTraceTimesDesc(a, b, -1);
  }

  static compareTraceTimesDesc(a?: Trace, b?: Trace, direction = 1): number {
    let result;
    if (a?.time && b?.time) {
      result = new Date(a.time).getTime() - new Date(b.time).getTime();
    }
    else if (a?.time && !b?.time) {
      result = 1;
    }
    else if (!a?.time && b?.time) {
      result = -1;
    }
    else {
      result = 0;
    }
    return result * direction;
  }

  public getDurationFormatted(start: string | Date, end?: string | Date) {
    const diff = moment(end).diff(moment(start));
    const duration = moment.duration(diff);
    const days = Math.floor(duration.asDays());
    const hours = Math.floor(duration.asHours() % 24);
    const minutes = Math.floor(duration.asMinutes() % 60);
    const seconds = Math.floor(duration.asSeconds() % 60);

    let result = seconds + ' seconds';
    if (minutes > 0) {
      result = minutes + ' minutes ' + result;
    }
    if (hours > 0) {
      result = hours + ' hours ' + result;
    }
    if (days > 0) {
      result = days + ' days ' + result;
    }

    return result;
  }

  public getCalendarFormats(showSeconds?: boolean) {
    if (showSeconds) {
      return {
        lastDay : '[yesterday at] HH:mm:ss',
        sameDay : '[today at] HH:mm:ss',
        nextDay : '[tomorrow at] HH:mm:ss',
        lastWeek : '[last] dddd [at] HH:mm:ss',
        nextWeek : 'dddd [at] HH:mm:ss',
        sameElse : 'YYYY-MM-DD HH:mm:ss'
      };
    }
    return {
      lastDay : '[yesterday at] HH:mm',
      sameDay : '[today at] HH:mm',
      nextDay : '[tomorrow at] HH:mm',
      lastWeek : '[last] dddd [at] HH:mm',
      nextWeek : 'dddd [at] HH:mm',
      sameElse : 'YYYY-MM-DD HH:mm'
    };
  }

  public getDateTimeFormat(): string {
    return [this.DEFAULT_DATE_FORMAT, this.DEFAULT_TIME_FORMAT].join(' ');
  }

  public getTimeFormat(): string {
    return this.DEFAULT_TIME_FORMAT;
  }
}
