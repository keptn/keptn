export default class DateUtil {

  static ONE_SECOND = 1000;
  static ONE_MINUTE = DateUtil.ONE_SECOND*60;
  static ONE_HOUR = DateUtil.ONE_MINUTE*60;
  static ONE_DAY = DateUtil.ONE_HOUR*24;

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
}
