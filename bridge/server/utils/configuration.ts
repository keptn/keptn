import { EnabledComponents, LogDestination } from './logger';

/**
 * Option to configure the Bridge-Server
 */
export interface BridgeOption {
  logging?: LogOptions;
}

interface LogOptions {
  destination?: LogDestination;
  enabledComponents?: string;
}

/**
 * Configuration object
 */
export interface BridgeConfiguration {
  logging: LogConfiguration;
}

interface LogConfiguration {
  destination: LogDestination;
  enabledComponents: EnabledComponents;
}

/**
 * @param options Customization options that override env var options.
 * @returns Returns the Bridge-server configuration.
 */
export function getConfiguration(options?: BridgeOption): BridgeConfiguration {
  const logDestination = options?.logging?.destination ?? LogDestination.STDOUT;
  const loggingComponents = Object.create({}) as EnabledComponents;
  const loggingComponentsString = options?.logging?.enabledComponents ?? process.env.LOGGING_COMPONENTS ?? '';
  if (loggingComponentsString.length > 0) {
    const components = loggingComponentsString.split(',').map((s) => s.trim());
    for (const component of components) {
      const [name, value] = parseComponent(component);
      loggingComponents[name] = value;
    }
  }
  return {
    logging: {
      destination: logDestination,
      enabledComponents: loggingComponents,
    },
  };
}

function parseComponent(component: string): [string, boolean] {
  // we expect only componentName = bool
  const split = component.split('=', 3);
  return [split[0].trim(), toBool(split[1])];
}

/**
 * Convert string to boolean. If the input is equal to false or 0, it returns false. True otherwise.
 * @param v string to convert.
 */
function toBool(v: string): boolean {
  const val = v.toLowerCase();
  return val !== '0' && val !== 'false';
}
