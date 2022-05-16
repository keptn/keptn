import * as gitUtils from './git-upstream.utils';
import { Project } from '../_models/project';

describe('GitUpstreamUtils', () => {
  it('should be HTTPS configuration', () => {
    expect(
      gitUtils.isGitHTTPS({
        https: {
          gitRemoteURL: 'https://myRemoteUrl',
          gitToken: '',
        },
      })
    ).toBe(true);
  });

  it('should be SSH configuration', () => {
    expect(
      gitUtils.isGitHTTPS({
        ssh: {
          gitRemoteURL: 'ssh://myGitUrl',
          gitPrivateKey: btoa('myPrivateKey'),
        },
      })
    ).toBe(false);
  });

  it('should be valid git upstream', () => {
    expect(
      gitUtils.isGitUpstreamValidSet({
        gitRemoteURL: 'https://myGitUrl',
        gitToken: 'myToken',
      })
    ).toBe(true);

    expect(
      gitUtils.isGitUpstreamValidSet({
        gitRemoteURL: 'https://myGitUrl',
        gitToken: 'myToken',
        gitUser: 'myUser',
      })
    ).toBe(true);
  });

  it('should not be valid git upstream', () => {
    expect(
      gitUtils.isGitUpstreamValidSet({
        gitRemoteURL: '',
        gitToken: 'myToken',
      })
    ).toBe(false);

    expect(
      gitUtils.isGitUpstreamValidSet({
        gitRemoteURL: 'https://myGitUrl',
        gitToken: '',
      })
    ).toBe(false);
  });

  it('should be HTTPS with proxy configuration', () => {
    expect(
      gitUtils.isGitWithProxy({
        https: {
          gitToken: 'myToken',
          gitRemoteURL: 'https://myGitUrl',
          gitProxyUser: '',
          gitProxyPassword: '',
          gitProxyScheme: 'https',
          gitProxyInsecure: false,
          gitProxyUrl: '0.0.0.0',
        },
      })
    ).toBe(true);

    expect(
      gitUtils.isGitWithProxy({
        https: {
          gitToken: 'myToken',
          gitRemoteURL: 'https://myGitUrl',
          gitProxyScheme: 'https',
          gitProxyInsecure: false,
          gitProxyUrl: '0.0.0.0',
        },
      })
    ).toBe(true);
  });

  it('should be HTTPS without proxy configuration', () => {
    expect(
      gitUtils.isGitWithProxy({
        https: {
          gitToken: 'myToken',
          gitRemoteURL: 'https://myGitUrl',
        },
      })
    ).toBe(false);

    expect(
      gitUtils.isGitWithProxy({
        https: {
          gitToken: 'myToken',
          gitRemoteURL: 'https://myGitUrl',
          gitProxyScheme: 'https',
          gitProxyUser: 'myUser',
          gitProxyPassword: 'myPassword',
          gitProxyInsecure: false,
        },
      })
    ).toBe(false);
  });

  it('should be project with HTTPS configuration', () => {
    expect(
      gitUtils.isGitInputWithHTTPS({
        gitRemoteURI: 'https://myGitUrl',
      } as Project)
    ).toBe(true);

    expect(
      gitUtils.isGitInputWithHTTPS({
        gitRemoteURI: 'http://myGitUrl',
      } as Project)
    ).toBe(true);
  });

  it('should be project with SSH configuration', () => {
    expect(
      gitUtils.isGitInputWithHTTPS({
        gitRemoteURI: 'ssh://myGitUrl',
      } as Project)
    ).toBe(false);
  });

  it('should return true if the https git remote url is empty', () => {
    expect(
      gitUtils.isRemoteUrlEmpty({
        https: {
          gitToken: '',
          gitRemoteURL: '',
        },
      })
    ).toBe(true);
  });

  it('should return false if the https git remote url is not empty', () => {
    expect(
      gitUtils.isRemoteUrlEmpty({
        https: {
          gitToken: '',
          gitRemoteURL: 'https://myGitUrl',
        },
      })
    ).toBe(false);
  });

  it('should return true if the ssh git remote url is empty', () => {
    expect(
      gitUtils.isRemoteUrlEmpty({
        ssh: {
          gitRemoteURL: '',
          gitPrivateKey: '',
        },
      })
    ).toBe(true);
  });

  it('should return false if the ssh git remote url is not empty', () => {
    expect(
      gitUtils.isRemoteUrlEmpty({
        ssh: {
          gitRemoteURL: 'ssh://myGitUrl',
          gitPrivateKey: '',
        },
      })
    ).toBe(false);
  });

  it('should return true if it is not a https or an ssh configuration', () => {
    expect(gitUtils.isRemoteUrlEmpty({ noupstream: '' })).toBe(true);
  });
});
