export enum Level {
	Info 	= "info   ",
	Warning = "warning",
	Error 	= "error  ",
	Debug 	= "debug  "
}

export interface EnabledComponents {
	[component: string]: boolean | undefined;
}

class LoggerImpl {
	public configure(destination: (message?: any, ...optionalParams: any[]) => void, enabledComponents: EnabledComponents = Object.create(null)): void {
		this._log = (destination != null) ? destination : undefined;
		this._components = enabledComponents;
	}

	public log(level: Level, component: string, msg: string): void {
		if (this._log == null) {
			return;
		}
		// Print debug level messages only if the component is enabled
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
}