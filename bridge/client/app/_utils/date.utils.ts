import moment from 'moment';
import { Trace } from '../_models/trace';
import { Injectable } from '@angular/core';
import { DateUtil as dtl } from '../../../shared/utils/date.utils';

@Injectable({
  providedIn: 'root',
})
export class DateUtil {
  public readonly DEFAULT_DATE_FORMAT = 'YYYY-MM-DD';
  public readonly DEFAULT_TIME_FORMAT = 'HH:mm';

  static compareTraceTimesAsc(a: Trace, b: Trace): number {
    return dtl.compareTraceTimesDesc(a, b, -1);
  }

  static compareTraceTimesDesc(a?: Trace, b?: Trace, direction = 1): number {
    return dtl.compareTraceTimesDesc(a, b, direction);
  }

  public getDurationFormatted(start: string | Date, end?: string | Date): string {
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

  public getCalendarFormats(
    showSeconds?: boolean,
    startUppercase?: boolean
  ): {
    lastDay: string;
    sameDay: string;
    nextDay: string;
    lastWeek: string;
    nextWeek: string;
    sameElse: string;
  } {
    if (showSeconds) {
      const calendarFormatsWithSeconds = {
        lastDay: '[yesterday at] HH:mm:ss',
        sameDay: '[today at] HH:mm:ss',
        nextDay: '[tomorrow at] HH:mm:ss',
        lastWeek: '[last] dddd [at] HH:mm:ss',
        nextWeek: 'dddd [at] HH:mm:ss',
        sameElse: 'YYYY-MM-DD HH:mm:ss',
      };
      if (startUppercase) {
        return {
          ...calendarFormatsWithSeconds,
          lastDay: '[Yesterday at] HH:mm:ss',
          sameDay: '[Today at] HH:mm:ss',
          nextDay: '[Tomorrow at] HH:mm:ss',
          lastWeek: '[Last] dddd [at] HH:mm:ss',
        };
      }

      return calendarFormatsWithSeconds;
    }

    const calendarFormats = {
      lastDay: '[yesterday at] HH:mm',
      sameDay: '[today at] HH:mm',
      nextDay: '[tomorrow at] HH:mm',
      lastWeek: '[last] dddd [at] HH:mm',
      nextWeek: 'dddd [at] HH:mm',
      sameElse: 'YYYY-MM-DD HH:mm',
    };
    if (startUppercase) {
      return {
        ...calendarFormats,
        lastDay: '[Yesterday at] HH:mm',
        sameDay: '[Today at] HH:mm',
        nextDay: '[Tomorrow at] HH:mm',
        lastWeek: '[Last] dddd [at] HH:mm',
      };
    }

    return calendarFormats;
  }

  public getDateTimeFormat(): string {
    return [this.DEFAULT_DATE_FORMAT, this.DEFAULT_TIME_FORMAT].join(' ');
  }

  public getTimeFormat(): string {
    return this.DEFAULT_TIME_FORMAT;
  }
}
