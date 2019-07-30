export class Logger {
  static log(keptnContext: string, message: string, logLevel: string = 'INFO'): void {
    if (process.env.NODE_ENV === 'production') {
      console.log(JSON.stringify({
        keptnContext,
        logLevel,
        message,
        keptnService: 'pitometer-service',
      }));
    } else {
      console.log(JSON.stringify({
        keptnContext,
        logLevel,
        message,
        keptnService: 'pitometer-service',
      }, null, 2));
    }
  }
}
