export class Logger {
  static info(keptnContext: string, message: string, keptnEntry: boolean = false) {
    try {
      const msg = JSON.stringify({
        keptnContext,
        keptnEntry,
        message,
        keptnService: 'eventbroker',
        logLevel: 'INFO',
      });
      console.log(msg);
    } catch (e) {
      console.log(e);
    }
  }

  static debug(keptnContext: string, message: string) {
    try {
      const msg = JSON.stringify({
        keptnContext,
        message,
        keptnService: 'eventbroker',
        logLevel: 'DEBUG',
      });
      console.log(msg);
    } catch (e) {
      console.log(e);
    }
  }

  static error(keptnContext: string, message: string) {
    try {
      const msg = JSON.stringify({
        keptnContext,
        message,
        keptnService: 'eventbroker',
        logLevel: 'ERROR',
      });
      console.log(msg);
    } catch (e) {
      console.log(e);
    }
  }
}
