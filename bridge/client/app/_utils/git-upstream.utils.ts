import { Project } from '../_models/project';
import {
  IGitData,
  IGitDataExtended,
  IGitHttps,
  IGitHTTPSProxy,
  IGitSsh,
  IRequiredGitData,
} from '../_interfaces/git-upstream';

export function isGitHTTPS(data: IGitDataExtended): data is IGitHttps {
  return data.hasOwnProperty('https');
}

export function isGitSSH(data: IGitDataExtended): data is IGitSsh {
  return data.hasOwnProperty('ssh');
}

export function isGitWithProxy(data: IGitHttps): data is IGitHTTPSProxy {
  return data.https.hasOwnProperty('gitProxyUrl') && !!(data.https as Partial<IGitHTTPSProxy['https']>).gitProxyUrl;
}

export function isGitInputWithHTTPS(project: Project): boolean {
  return !project.gitRemoteURI?.startsWith('ssh://');
}

export function isGitUpstreamValidSet(gitUpstream: IGitData): gitUpstream is IRequiredGitData {
  return !!(gitUpstream.gitToken && gitUpstream.gitRemoteURL);
}

export function getGitData(data: IGitDataExtended): IGitHttps['https'] | IGitSsh['ssh'] | undefined {
  if (isGitHTTPS(data)) {
    return data.https;
  }
  if (isGitSSH(data)) {
    return data.ssh;
  }
  return undefined;
}

export function isRemoteUrlEmpty(gitInputData: IGitDataExtended): boolean {
  return (
    (isGitHTTPS(gitInputData) && !gitInputData.https.gitRemoteURL) ||
    (isGitSSH(gitInputData) && !gitInputData.ssh.gitRemoteURL) ||
    (!isGitHTTPS(gitInputData) && !isGitSSH(gitInputData))
  );
}
