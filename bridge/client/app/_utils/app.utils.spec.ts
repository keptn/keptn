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
});
