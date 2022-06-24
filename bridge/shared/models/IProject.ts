import { IStage } from '../interfaces/stage';

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

export interface IGitHTTPSData {
  token: string;
  certificate?: string;
  insecureSkipTLS: boolean;
  proxy?: IProxy;
}

export interface IGitSshData {
  privateKey: string;
  privateKeyPass?: string;
}

export interface IGitHTTPSConfiguration extends IGitBasicConfiguration {
  https: IGitHTTPSData;
}

export interface IGitSSHConfiguration extends IGitBasicConfiguration {
  ssh: IGitSshData;
}

export type IGitDataExtended = IGitSSHConfiguration | IGitHTTPSConfiguration | IGitBasicConfiguration;

export interface IProject {
  projectName: string;
  gitCredentials?: IGitDataExtended; // optional because of configuration-service.
  shipyardVersion?: string;
  creationDate: string;
  shipyard: string;
  stages: IStage[];
}
