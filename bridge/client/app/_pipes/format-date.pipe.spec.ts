import { FormatDatePipe } from './format-date.pipe';

describe('FormatDatePipe', () => {
  const pipe = new FormatDatePipe();

  it('should format date with default format', () => {
    expect(pipe.transform('2022-04-18T14:12:45.000Z')).toBe('2022-04-18 16:12');
  });

  it('should format date with custom format', () => {
    expect(pipe.transform('2022-04-18T14:12:45.000Z', 'YYYY-MM-DD HH:mm:ss')).toBe('2022-04-18 16:12:45');
  });

  it('should not format undefined', () => {
    expect(pipe.transform(undefined)).toBe(undefined);
  });

  it('should not format empty string', () => {
    expect(pipe.transform('')).toBe('');
  });
});
