export interface IGitData {
  gitRemoteURL?: string;
  gitUser?: string;
  gitToken?: string;
  gitFormValid?: boolean;
}

export interface IRequiredGitData extends IGitData {
  gitRemoteURL: string;
  gitToken: string;
}

interface IGitHttpsWithoutProxy {
  https: {
    gitRemoteURL: string;
    gitUser?: string;
    gitToken: string;
    gitPemCertificate?: string;
  };
}

export interface IProxy {
  gitProxyInsecure: boolean;
  gitProxyScheme: 'https' | 'http';
  gitProxyUrl: string; // Host and Port
  gitProxyUser?: string;
  gitProxyPassword?: string;
}

export interface IGitHTTPSProxy extends IGitHttpsWithoutProxy {
  https: IGitHttpsWithoutProxy['https'] & IProxy;
}

export interface IGitSshData {
  gitRemoteURL: string;
  gitUser?: string;
}

export interface ISshKeyData {
  gitPrivateKey: string;
  gitPrivateKeyPass?: string;
}

export interface IGitSsh {
  ssh: IGitSshData & ISshKeyData;
}

export type IGitHttps = IGitHttpsWithoutProxy | IGitHTTPSProxy;

export type IGitDataExtended = IGitHttps | IGitSsh;
