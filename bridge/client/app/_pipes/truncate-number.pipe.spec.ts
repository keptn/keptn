import { TruncateNumberPipe } from './truncate-number';

describe('TruncateNumberPipe', () => {
  const pipe = new TruncateNumberPipe();

  it('truncates 99.999999999 to 1.99', () => {
    expect(pipe.transform(99.999999999, 2)).toBe(99.99);
  });

  it('truncates 70.333333333 to 1.33', () => {
    expect(pipe.transform(70.333333333, 2)).toBe(70.33);
  });

  it('truncates 3.141592653 to 3.14', () => {
    expect(pipe.transform(3.141592653, 2)).toBe(3.14);
  });

  it('works also with less decimals, e.g. display 99.1 as 99.1', () => {
    expect(pipe.transform(99.1, 2)).toBe(99.1);
  });

  it('works also without decimals, e.g. display 99 as 99', () => {
    expect(pipe.transform(99, 2)).toBe(99);
  });
});
