import * as gitUtils from './git-upstream.utils';

describe('GitUpstreamUtils', () => {
  it('should be HTTPS configuration', () => {
    expect(
      gitUtils.isGitHTTPS({
        remoteURL: 'https://myRemoteUrl',
        https: {
          token: '',
          insecureSkipTLS: false,
        },
      })
    ).toBe(true);
  });

  it('should be SSH configuration', () => {
    expect(
      gitUtils.isGitHTTPS({
        remoteURL: 'ssh://myGitUrl',
        ssh: {
          privateKey: btoa('myPrivateKey'),
        },
      })
    ).toBe(false);
  });

  it('should be valid git upstream', () => {
    expect(
      gitUtils.isGitUpstreamValidSet({
        remoteURL: 'https://myGitUrl',
        token: 'myToken',
      })
    ).toBe(true);

    expect(
      gitUtils.isGitUpstreamValidSet({
        remoteURL: 'https://myGitUrl',
        token: 'myToken',
        user: 'myUser',
      })
    ).toBe(true);
  });

  it('should not be valid git upstream', () => {
    expect(
      gitUtils.isGitUpstreamValidSet({
        remoteURL: '',
        token: 'myToken',
      })
    ).toBe(false);

    expect(
      gitUtils.isGitUpstreamValidSet({
        remoteURL: 'https://myGitUrl',
        token: '',
      })
    ).toBe(false);
  });

  it('should be project with HTTPS configuration', () => {
    expect(
      gitUtils.isGitInputWithSSH({
        remoteURL: 'https://myGitUrl',
      })
    ).toBe(false);

    expect(
      gitUtils.isGitInputWithSSH({
        remoteURL: 'http://myGitUrl',
      })
    ).toBe(false);
  });

  it('should be project with SSH configuration', () => {
    expect(
      gitUtils.isGitInputWithSSH({
        remoteURL: 'ssh://myGitUrl',
      })
    ).toBe(true);
  });
});
