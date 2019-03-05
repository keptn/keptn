import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { KeptnRequestModel } from '../lib/types/KeptnRequestModel';

@injectable()
export class GitHubService {

  constructor() {}

  public async handleGitHubEvent(gitHubEventType: string, githubEventPayload: any): Promise<void> {
    // for now, only handle events of type 'push'
    if (gitHubEventType !== 'push') {
      return;
    }

    const refSplit = githubEventPayload.ref.split('/');
    const environment = refSplit[refSplit.length - 1];
    const commitMessage = githubEventPayload.head_commit.message;
    const githubOrg = githubEventPayload.repository.owner.name;

    const keptnEvent: KeptnRequestModel = new KeptnRequestModel();
    keptnEvent.type = KeptnRequestModel.EVENT_TYPES.CONFIGURATION_CHANGED;
    axios.post('http://event-broker', keptnEvent);
  }
}
