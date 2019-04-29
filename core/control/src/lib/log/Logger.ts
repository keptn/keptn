export class Logger {
  static keptnContext;
  static info(keptnContext: string = '', message: string) {
    Logger.logMessage(keptnContext, message, 'INFO');
  }

  static debug(keptnContext: string = '', message: string) {
    Logger.logMessage(keptnContext, message, 'DEBUG');
  }

  static error(keptnContext: string = '', message: string) {
    Logger.logMessage(keptnContext, message, 'ERROR');
  }

  private static logMessage(keptnContext: string, message: string, logLevel: string) {
    try {
      if (keptnContext !== undefined && keptnContext !== '') {
        Logger.keptnContext = keptnContext;
      }
      const msg = JSON.stringify({
        message,
        logLevel,
        keptnContext: Logger.keptnContext,
        keptnService: 'control',
      });
      console.log(msg);
    } catch (e) {
      console.log(e);
    }
  }
}
