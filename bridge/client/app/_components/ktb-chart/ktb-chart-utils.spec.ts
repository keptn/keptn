import { getColor, getIconStyle, getTooltipPosition, replaceSpace } from './ktb-chart-utils';

describe('KtbChartUtils', () => {
  it('should return a color', () => {
    // given
    const index = 10_000;

    // when
    const color = getColor(index);

    // then
    expect(color).toBe('#ef651f');
  });

  it('should get an icon style by index', () => {
    // given
    const index = 12;

    // when
    const iconStyle = getIconStyle(index, false);

    // then
    expect(iconStyle).toBe('--dt-icon-color: #93060e');
  });

  it('should get disabled icon style by index', () => {
    // given
    const index = 12;

    // when
    const iconStyle = getIconStyle(index, true);

    // then
    expect(iconStyle).toBe('--dt-icon-color: #cccccc');
  });

  it('should replaces spaces', () => {
    // given
    const name = 'A metric  with spaces';

    // when
    const actual = replaceSpace(name);

    // then
    expect(actual).toBe('A-metric--with-spaces');
  });

  describe(getTooltipPosition.name, () => {
    it('should calculate the tooltip position', () => {
      // given
      const window = { width: 1111, height: 999 };
      const tooltip = { width: 155, height: 177 };
      const barArea = { width: 50, height: 100, top: 21, left: 33 };

      // when
      const actual = getTooltipPosition(window, tooltip, barArea);

      // then
      expect(actual).toEqual({ top: 141, left: 33 });
    });

    it('should calculate the tooltip position when there is no space', () => {
      // given
      const window = { width: 1111, height: 999 };
      const tooltip = { width: 155, height: 177 };
      const barArea = { width: 50, height: 100, top: 980, left: 1100 };

      // when
      const actual = getTooltipPosition(window, tooltip, barArea);

      // then
      expect(actual).toEqual({ top: 812, left: 946 });
    });
  });
});
