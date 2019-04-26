export class Logger {
  static info(keptnContext: string, message: string) {
    try {
      const msg = JSON.stringify({
        keptnContext,
        message,
        keptnService: 'control',
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
        keptnService: 'control',
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
        keptnService: 'control',
        logLevel: 'ERROR',
      });
      console.log(msg);
    } catch (e) {
      console.log(e);
    }
  }
}
