import { init as initApp } from './app';
import { BridgeConfiguration } from './interfaces/configuration';
import { envToConfiguration, getConfiguration } from './utils/configuration';
import { ComponentLogger, logger } from './utils/logger';

const PORT = normalizePort(process.env.PORT || '3000');
const HOST = process.env.HOST || '0.0.0.0';
let configuration: BridgeConfiguration;
try {
  configuration = getConfiguration(envToConfiguration(process.env));
} catch (e) {
  console.log(`Error while configuring the application. Cause: ${e}`);
  process.exit(1);
}

// init destination and debug flags
logger.configure(
  configuration.logging.destination,
  configuration.logging.enabledComponents,
  configuration.logging.defaultLogLevel
);

const log = new ComponentLogger('Server');

if (typeof PORT === 'number') {
  (async (): Promise<void> => {
    try {
      const app = await initApp(configuration);
      app.set('port', PORT);

      /**
       * Listen on provided port, on all network interfaces.
       */
      log.info(`Running on http://${HOST}:${PORT}`);
      app.listen(PORT, HOST);
    } catch (e) {
      log.error(`Error while starting the application. Cause : ${e}`);
      process.exit(1);
    }
  })();
} else {
  log.error(`Error while starting the application. Invalid port`);
  process.exit(1);
}

/**
 * Normalize a port into a number or false.
 */
function normalizePort(val: string): number | boolean {
  const parsedPort = parseInt(val, 10);
  return !isNaN(parsedPort) && parsedPort >= 0 ? parsedPort : false;
}
