import { IStage } from './stage';

export interface IGitBasicConfiguration {
  remoteURL: string;
  user?: string;
}

export interface IProxy {
  url: `${string /*host*/}:${number /*port*/}`;
  scheme: 'http' | 'https';
  user?: string;
  password?: string;
}

export interface IGitHttpsData {
  token: string;
  certificate?: string;
  insecureSkipTLS: boolean;
  proxy?: IProxy;
}

export interface IGitSshData {
  privateKey: string;
  privateKeyPass?: string;
}

export interface IGitHttpsConfiguration extends IGitBasicConfiguration {
  https?: IGitHttpsData;
}

export interface IGitSshConfiguration extends IGitBasicConfiguration {
  ssh?: IGitSshData;
}

export type IGitDataExtended = IGitSshConfiguration | IGitHttpsConfiguration;

export interface IProject {
  projectName: string;
  gitCredentials?: IGitDataExtended; // optional because of configuration-service.
  shipyardVersion?: string;
  creationDate: string;
  shipyard: string;
  stages: IStage[];
}
