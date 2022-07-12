export interface IGitData {
  remoteURL?: string;
  user?: string;
  token?: string;
  valid: boolean;
}

export interface IRequiredGitData extends Omit<IGitData, 'valid'> {
  remoteURL: string;
  token: string;
}
