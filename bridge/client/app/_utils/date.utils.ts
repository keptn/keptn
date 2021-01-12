import * as moment from "moment";
import {Trace} from "../_models/trace";

export default class DateUtil {
  static getDurationFormatted(start, end) {
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

  static getCalendarFormats(showSeconds?: boolean) {
    if(showSeconds) {
      return {
        lastDay : '[Yesterday at] HH:mm:ss',
        sameDay : '[Today at] HH:mm:ss',
        nextDay : '[Tomorrow at] HH:mm:ss',
        lastWeek : '[last] dddd [at] HH:mm:ss',
        nextWeek : 'dddd [at] HH:mm:ss',
        sameElse : 'YYYY-MM-DD HH:mm:ss'
      };
    }
    return {
      lastDay : '[Yesterday at] HH:mm',
      sameDay : '[Today at] HH:mm',
      nextDay : '[Tomorrow at] HH:mm',
      lastWeek : '[last] dddd [at] HH:mm',
      nextWeek : 'dddd [at] HH:mm',
      sameElse : 'YYYY-MM-DD HH:mm'
    };
  }

  static compareTraceTimesAsc(a: Trace, b: Trace) {
    return new Date(b.time).getTime() - new Date(a.time).getTime();
  }

  static compareTraceTimesDesc(a: Trace, b: Trace) {
    return new Date(a.time).getTime() - new Date(b.time).getTime();
  }
}
