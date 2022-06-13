import * as FS from 'fs';

export enum Level {
  INFO = 'info   ',
  WARNING = 'warning',
  ERROR = 'error  ',
  DEBUG = 'debug  ',
}

export enum LogDestination {
  FILE = 'file',
  STDOUT = 'stdOutput',
}

export interface EnabledComponents {
  [component: string]: boolean | undefined;
}

export class LoggerImpl {
  public configure(
    destination: LogDestination = LogDestination.STDOUT,
    enabledComponents: EnabledComponents = Object.create(null)
  ): void {
    this._log = destination == LogDestination.STDOUT ? console.log : this.fileLog;
    this._components = enabledComponents;
  }

  private fileLog(msg: string): void {
    FS.writeFile('./bridge-server.log', msg, { flag: 'a+' }, (err) => {
      console.error(err);
    });
  }

  public log(level: Level, component: string, msg: string): void {
    if (this._log == null) {
      return;
    }
    // Print debug levels only if the component is enabled
    if (level === Level.DEBUG && this._components != null && this._components[component] !== true) {
      return;
    }

    const date = new Date(Date.now()).toISOString();
    const message = `[Keptn] ${date} ${level} [${component}] ${msg}`;

    this._log(message);
  }

  public info = (component: string, msg: string): void => this.log(Level.INFO, component, msg);
  public warning = (component: string, msg: string): void => this.log(Level.WARNING, component, msg);
  public error = (component: string, msg: string): void => this.log(Level.ERROR, component, msg);
  public debug = (component: string, msg: string): void => this.log(Level.DEBUG, component, msg);

  private _log?: (msg: string) => void;
  private _components?: EnabledComponents;
}

export const logger = new LoggerImpl();

export interface Logger {
  debug(msg: string): void;
  info(msg: string): void;
  warning(msg: string): void;
  error(msg: string): void;
}

export class ComponentLogger implements Logger {
  constructor(public readonly componentName: string) {}

  public info(msg: string): void {
    logger.info(this.componentName, msg);
  }

  public warning(msg: string): void {
    logger.warning(this.componentName, msg);
  }

  public error(msg: string): void {
    logger.error(this.componentName, msg);
  }

  public debug(msg: string): void {
    logger.debug(this.componentName, msg);
  }

  /**
   * @param o the object
   * @returns a string with the key value pair properties.
   */
  public prettyPrint(o: object): string {
    return Object.entries(o)
      .filter(([, v]) => v != null)
      .map(([k, v]) => `${k}=${v}`)
      .join(', ');
  }
}
