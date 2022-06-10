import { init as initApp } from './app';
import { getConfiguration } from './utils/configuration';
import { logger } from './utils/logger';

const PORT = normalizePort(process.env.PORT || '3000');
const HOST = process.env.HOST || '0.0.0.0';

const configuration = getConfiguration();

// init destination and debug flags
logger.configure(configuration.logging.destination, configuration.logging.enabledComponents);

if (typeof PORT === 'number') {
  (async (): Promise<void> => {
    try {
      const app = await initApp();
      app.set('port', PORT);

      /**
       * Listen on provided port, on all network interfaces.
       */
      console.log(`Running on http://${HOST}:${PORT}`);
      app.listen(PORT, HOST);
    } catch (e) {
      console.log(`Error while starting the application. Cause : ${e}`);
      process.exit(1);
    }
  })();
} else {
  console.log(`Error while starting the application. Invalid port`);
  process.exit(1);
}

/**
 * Normalize a port into a number or false.
 */
function normalizePort(val: string): number | boolean {
  const parsedPort = parseInt(val, 10);
  return !isNaN(parsedPort) && parsedPort >= 0 ? parsedPort : false;
}
