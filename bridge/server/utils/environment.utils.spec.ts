import { EnvironmentUtils } from './environment.utils';
import { ServerFeatureFlags } from '../feature-flags';

describe('Test environment utils', () => {
  it('should format number and return undefined', () => {
    for (const input of [undefined, 'undefined', 'not a number']) {
      expect(EnvironmentUtils.getNumber(input)).toBe(undefined);
    }
  });

  it('should format number and return positive integer', () => {
    expect(EnvironmentUtils.getNumber('10')).toBe(10);
    expect(EnvironmentUtils.getNumber('-10')).toBe(10);
    expect(EnvironmentUtils.getNumber('2.55')).toBe(2);
  });

  it('should format number and return integer', () => {
    expect(EnvironmentUtils.getNumber('10', false)).toBe(10);
    expect(EnvironmentUtils.getNumber('-10', false)).toBe(-10);
  });

  it('should set server feature flags to true', () => {
    // given
    const flags = { OAUTH_ENABLED: false };
    // when
    EnvironmentUtils.setFeatureFlags(
      {
        OAUTH_ENABLED: 'true',
      },
      flags
    );
    // then
    expect(flags).toEqual({ OAUTH_ENABLED: true });
  });

  it('should set server feature flag to false', () => {
    for (const input of ['false', '', 'FALSE', 'TRUE']) {
      // given
      const flags = { OAUTH_ENABLED: true };
      // when
      EnvironmentUtils.setFeatureFlags(
        {
          OAUTH_ENABLED: input,
        },
        flags
      );
      // then
      expect(flags).toEqual({ OAUTH_ENABLED: false });
    }
  });

  it('should use default server feature flag if not provided', () => {
    // given
    const flags = new ServerFeatureFlags();
    // when
    EnvironmentUtils.setFeatureFlags({}, flags);
    // then
    expect(flags).toEqual({ OAUTH_ENABLED: false });
  });

  it('should not add additional server feature flags', () => {
    // given
    const flags = { OAUTH_ENABLED: false };
    // when
    EnvironmentUtils.setFeatureFlags(
      {
        OAUTH_ENABLED: 'true',
        ANOTHER_FLAG: 'true',
      },
      flags
    );
    // then
    expect(flags).toEqual({ OAUTH_ENABLED: true });
  });
});
