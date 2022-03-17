import { Project } from '../_models/project';
import { IGitData, IGitDataExtended, IGitHttps, IGitHTTPSProxy, IRequiredGitData } from '../_interfaces/git-upstream';

export function isGitHTTPS(data: IGitDataExtended): data is IGitHttps {
  return data.hasOwnProperty('https');
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
