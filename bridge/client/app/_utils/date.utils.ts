import moment from "moment";
import {Trace} from "../_models/trace";
import {Injectable} from "@angular/core";

@Injectable({
  providedIn: 'root'
})
export class DateUtil {

  public DEFAULT_DATE_FORMAT = 'YYYY-MM-DD';
  public DEFAULT_TIME_FORMAT = 'HH:mm';

  public getDurationFormatted(start, end?) {
    let diff = moment(end).diff(moment(start));
    let duration = moment.duration(diff);

    let days = Math.floor(duration.asDays());
    let hours = Math.floor(duration.asHours()%24);
    let minutes = Math.floor(duration.asMinutes()%60);
    let seconds = Math.floor(duration.asSeconds()%60);

    let result = seconds+' seconds';
    if(minutes > 0)
      result = minutes+' minutes '+result;
    if(hours > 0)
      result = hours+' hours '+result;
    if(days > 0)
      result = days+' days '+result;

    return result;
  }

  public getCalendarFormats(showSeconds?: boolean) {
    if(showSeconds) {
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

  public getDateTimeFormat() {
    return [this.DEFAULT_DATE_FORMAT, this.DEFAULT_TIME_FORMAT].join(" ");
  }

  public getTimeFormat() {
    return this.DEFAULT_TIME_FORMAT;
  }

  static compareTraceTimesAsc(a: Trace, b: Trace) {
    return new Date(b.time).getTime() - new Date(a.time).getTime();
  }

  static compareTraceTimesDesc(a: Trace, b: Trace) {
    return new Date(a.time).getTime() - new Date(b.time).getTime();
  }
}
