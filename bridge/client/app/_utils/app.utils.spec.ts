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

  it('isValidJson should return true if given valid JSON string', () => {
    // given
    const jsonString = '{"name":"John"}';

    // when
    const isJson = AppUtils.isValidJson(jsonString);

    // then
    expect(isJson).toBeTruthy();
  });

  it('isValidJson should return false if given invalid JSON string', () => {
    // given
    for (const jsonString of ['{"name":"John",}', '{name:John}']) {
      // when
      const isJson = AppUtils.isValidJson(jsonString);

      // then
      expect(isJson).toBeFalsy();
    }
  });

  it('isValidUrl should return true if given valid URL', () => {
    // given
    for (const url of ['https://keptn.sh', 'https://tutorials.keptn.sh/']) {
      // when
      const isUrl = AppUtils.isValidUrl(url);

      // then
      expect(isUrl).toBeTruthy();
    }
  });

  it('isValidUrl should return false if given invalid URL', () => {
    // given
    const url = 'keptn.sh';

    // when
    const isUrl = AppUtils.isValidUrl(url);

    // then
    expect(isUrl).toBeFalsy();
  });
});
