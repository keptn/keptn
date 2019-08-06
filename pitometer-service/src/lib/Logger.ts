export class Logger {
  static log(keptnContext: string, eventId: string, message: string, logLevel: string = 'INFO'): void {
    if (process.env.NODE_ENV === 'production') {
      console.log(JSON.stringify({
        keptnContext,
        eventId,
        keptnService: 'pitometer-service',
        timestamp: Date.now(),
        logLevel,
        message,
      }));
    } else {
      console.log(JSON.stringify({
        keptnContext,
        eventId,
        keptnService: 'pitometer-service',
        timestamp: Date.now(),
        logLevel,
        message,
      }, null, 2));
    }
  }
}
