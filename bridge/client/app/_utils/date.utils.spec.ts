import { DateUtil } from './date.utils';

describe('DateUtils', () => {
  it('gets calendar formats without seconds, and not uppercase', () => {
    const expectedFormat = {
      lastDay: '[yesterday at] HH:mm',
      sameDay: '[today at] HH:mm',
      nextDay: '[tomorrow at] HH:mm',
      lastWeek: '[last] dddd [at] HH:mm',
      nextWeek: 'dddd [at] HH:mm',
      sameElse: 'YYYY-MM-DD HH:mm',
    };
    expect(new DateUtil().getCalendarFormats()).toEqual(expectedFormat);
  });

  it('gets calendar formats with seconds, and not uppercase', () => {
    const expectedFormat = {
      lastDay: '[yesterday at] HH:mm:ss',
      sameDay: '[today at] HH:mm:ss',
      nextDay: '[tomorrow at] HH:mm:ss',
      lastWeek: '[last] dddd [at] HH:mm:ss',
      nextWeek: 'dddd [at] HH:mm:ss',
      sameElse: 'YYYY-MM-DD HH:mm:ss',
    };
    expect(new DateUtil().getCalendarFormats(true)).toEqual(expectedFormat);
  });

  it('gets calendar formats without seconds, and uppercase', () => {
    const expectedFormat = {
      lastDay: '[Yesterday at] HH:mm',
      sameDay: '[Today at] HH:mm',
      nextDay: '[Tomorrow at] HH:mm',
      lastWeek: '[Last] dddd [at] HH:mm',
      nextWeek: 'dddd [at] HH:mm',
      sameElse: 'YYYY-MM-DD HH:mm',
    };
    expect(new DateUtil().getCalendarFormats(false, true)).toEqual(expectedFormat);
  });

  it('gets calendar formats with seconds, and uppercase', () => {
    const expectedFormat = {
      lastDay: '[Yesterday at] HH:mm:ss',
      sameDay: '[Today at] HH:mm:ss',
      nextDay: '[Tomorrow at] HH:mm:ss',
      lastWeek: '[Last] dddd [at] HH:mm:ss',
      nextWeek: 'dddd [at] HH:mm:ss',
      sameElse: 'YYYY-MM-DD HH:mm:ss',
    };
    expect(new DateUtil().getCalendarFormats(true, true)).toEqual(expectedFormat);
  });
});
