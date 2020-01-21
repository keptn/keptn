export default class DateUtil {
  static getCalendarFormats() {
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
