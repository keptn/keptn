import { ClientFeatureFlags, ServerFeatureFlags } from '../feature-flags';

export class EnvironmentUtils {
  public static getNumber(num: string | undefined, positive = true): number | undefined {
    let result = undefined;

    if (num !== undefined) {
      const count = parseInt(num, 10);
      if (!isNaN(count)) {
        result = positive ? Math.abs(count) : count;
      }
    }
    return result;
  }

  public static setFeatureFlags(
    env: Record<string, string | undefined>,
    defaultFlags: ServerFeatureFlags | ClientFeatureFlags
  ): void {
    for (const flag of Object.keys(defaultFlags)) {
      if (flag in env) {
        (defaultFlags as unknown as Record<string, boolean>)[flag] = env[flag] === 'true';
      }
    }
  }
}
