import { expect } from 'chai';
import 'mocha';
import { GitHubService } from './GitHubService';
import nock from 'nock';

describe('GitHubService', function () {
  this.timeout(0);
  let gitHubService: GitHubService;

  beforeEach(() => {
    gitHubService = new GitHubService();
  });

  it('Should send a new ConfigChangeEvent when a valid commit is detected', async () => {
    const message = {
      ref: 'refs/heads/production',
      head_commit: {
        message: '[keptn-config-change]:carts:pr-0',
      },
      repository: {
        name: 'sockshop',
      },
    };
    nock('http://event-broker.keptn.svc.cluster.local', {
      filteringScope: () => {
        return true;
      },
    })
      .post('/keptn')
      .reply(200, {});

    const result = await gitHubService.handleGitHubEvent('push', message);
    expect(result).is.true;
  });
  it('Should not send a ConfigChangeEvent if no valid commit message is detected', async () => {
    const message = {
      ref: 'refs/heads/production',
      head_commit: {
        message: 'wubba lubba dub dub',
      },
      repository: {
        name: 'sockshop',
      },
    };

    const result = await gitHubService.handleGitHubEvent('push', message);
    expect(result).is.false;
  });
  it('Should not send a ConfigChangeEvent if no valid commit message is detected (2)', async () => {
    const message = {
      ref: 'refs/heads/production',
      head_commit: {
        message: '[keptn-config-change]:carts',
      },
      repository: {
        name: 'sockshop',
      },
    };

    const result = await gitHubService.handleGitHubEvent('push', message);
    expect(result).is.false;
  });
  it('Should not send a ConfigChangeEvent if the eventType is not set to push', async () => {
    const message = {
      ref: 'refs/heads/production',
      head_commit: {
        message: '[keptn-config-change]:carts:pr-0',
      },
      repository: {
        name: 'sockshop',
      },
    };

    const result = await gitHubService.handleGitHubEvent('pull', message);
    expect(result).is.false;
  });
  it('Should return false if an invalid event structure is received', async () => {
    const message = {
      foo: 'bar',
    };

    const result = await gitHubService.handleGitHubEvent('pull', message);
    expect(result).is.false;
  });
});
