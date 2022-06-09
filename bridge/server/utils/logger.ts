import * as FS from 'fs';

export enum Level {
	Info 	= "info   ",
	Warning = "warning",
	Error 	= "error  ",
	Debug 	= "debug  "
}

export enum LogDestination {
	File   = "file",
	StdOut = "stdOutput",
}

export interface EnabledComponents {
	[component: string]: boolean | undefined;
}

class LoggerImpl {
	public configure(destination: LogDestination = LogDestination.StdOut, enabledComponents: EnabledComponents = Object.create(null)): void {
		this._log = (destination == LogDestination.StdOut) ? console.log : this.fileLog;
		this._components = enabledComponents;
	}

	private fileLog(msg: string) {
		FS.writeFile('./bridge-server.log', msg, { flag: 'a+' }, err => {
			console.error(err);
		});
	}

	public log(level: Level, component: string, msg: string): void {
		if (this._log == null) {
			return;
		}
		// Print debug levels only if the component is enabled
		if(level === Level.Debug && this._components != null && this._components[component] !== true) {
			return;
		}

		const date = new Date(Date.now()).toISOString();
		const message = `[Keptn] ${date} ${level} [${component}] ${msg}`;

		this._log(message);
	}

	public info  = (component: string, msg: string): void => this.log(Level.Info,  component, msg);
	public warning  = (component: string, msg: string): void => this.log(Level.Warning,  component, msg);
	public error = (component: string, msg: string): void => this.log(Level.Error, component, msg);
	public debug = (component: string, msg: string): void => this.log(Level.Debug, component, msg);

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
			.join(", ");
	}
}