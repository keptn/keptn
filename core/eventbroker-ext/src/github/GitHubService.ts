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

    let project;
    let app;
    let branch; = eventPayload.ref.
    let commitMessage;

    const keptnEvent: KeptnRequestModel = new KeptnRequestModel();
    keptnEvent.type = KeptnRequestModel.EVENT_TYPES.CONFIGURATION_CHANGED;
    axios.post('http://event-broker', keptnEvent);
  }

}