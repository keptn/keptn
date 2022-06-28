import { IGitDataExtended, IGitHTTPSConfiguration, IGitSSHConfiguration } from '../../../shared/interfaces/project';
import {
  IGitData,
  IRequiredGitData,
} from '../_components/ktb-project-settings/ktb-project-settings-git/ktb-project-settings-git.utils';

export function isGitHTTPS(data: IGitDataExtended): data is IGitHTTPSConfiguration {
  return !isGitInputWithSSH(data) || data.hasOwnProperty('https');
}

export function isGitSSH(data: IGitDataExtended): data is IGitSSHConfiguration {
  return isGitInputWithSSH(data) || data.hasOwnProperty('ssh');
}

export function isGitInputWithSSH(gitCredentials: IGitDataExtended | undefined): boolean {
  return !!gitCredentials?.remoteURL.startsWith('ssh://');
}

export function isGitUpstreamValidSet(gitUpstream: Omit<IGitData, 'valid'>): gitUpstream is IRequiredGitData {
  return !!(gitUpstream.token && gitUpstream.remoteURL);
}
