import { AppUtils } from './app.utils';

describe('AppUtils', () => {
  it('formats Number 0.333333333 to 3 decimals', () => {
    expect(AppUtils.formatNumber(0.333333333)).toBe(0.333);
  });

  it('formats Number 17.333333333 to 2 decimals', () => {
    expect(AppUtils.formatNumber(17.333333333)).toBe(17.33);
  });

  it('formats Number 170.333333333 to 1 decimals', () => {
    expect(AppUtils.formatNumber(170.333333333)).toBe(170.3);
  });

  it('formats Number 1700.333333333 to 0 decimals', () => {
    expect(AppUtils.formatNumber(1700.333333333)).toBe(1700);
  });

  it('formats Number -0.333333333 to 3 decimals', () => {
    expect(AppUtils.formatNumber(-0.333333333)).toBe(-0.333);
  });

  it('formats Number -17.333333333 to 2 decimals', () => {
    expect(AppUtils.formatNumber(-17.333333333)).toBe(-17.33);
  });

  it('formats Number -170.333333333 to 1 decimals', () => {
    expect(AppUtils.formatNumber(-170.333333333)).toBe(-170.3);
  });

  it('formats Number -1700.333333333 to 0 decimals', () => {
    expect(AppUtils.formatNumber(-1700.333333333)).toBe(-1700);
  });

  it('should correctly set host and port', () => {
    for (const port of ['5000', '1', '0', '123456']) {
      // when
      const urlData = AppUtils.splitURLPort(`https://0.0.0.0:${port}`);

      // then
      expect(urlData).toEqual({
        host: 'https://0.0.0.0',
        port: port,
      });
    }
  });

  it('should return empty port if found port is not a number', () => {
    // when
    const urlData = AppUtils.splitURLPort('https://0.0.0.0:asdf');

    // then
    expect(urlData).toEqual({
      host: 'https://0.0.0.0:asdf',
      port: '',
    });
  });

  it('should correctly set host and not set port if input data does not contain port', () => {
    for (const host of ['0.0.0.0', 'https://0.0.0.0']) {
      // when
      const urlData = AppUtils.splitURLPort(host);

      // then
      expect(urlData).toEqual({
        host,
        port: '',
      });
    }
  });
});
