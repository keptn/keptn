import { IGitDataExtended, IGitHttpsConfiguration, IGitSshConfiguration } from '../../../shared/interfaces/project';
import {
  IGitData,
  IRequiredGitData,
} from '../_components/ktb-project-settings/ktb-project-settings-git/ktb-project-settings-git.utils';

export function isGitHttps(data: IGitDataExtended): data is IGitHttpsConfiguration {
  return !isGitInputWithSsh(data) || data.hasOwnProperty('https');
}

export function isGitSsh(data: IGitDataExtended): data is IGitSshConfiguration {
  return isGitInputWithSsh(data) || data.hasOwnProperty('ssh');
}

export function isGitInputWithSsh(gitCredentials: IGitDataExtended | undefined): boolean {
  return !!gitCredentials?.remoteURL.startsWith('ssh://');
}

export function isGitUpstreamValidSet(gitUpstream: Omit<IGitData, 'valid'>): gitUpstream is IRequiredGitData {
  return !!(gitUpstream.token && gitUpstream.remoteURL);
}
