import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { KeptnRequestModel } from '../lib/types/KeptnRequestModel';
import { ConfigurationModel } from '../lib/types/ConfigurationModel';

@injectable()
export class GitHubService {

  constructor() {}

  public async handleGitHubEvent(
    gitHubEventType: string,
    githubEventPayload: any): Promise<boolean> {
    // for now, only handle events of type 'push'
    if (gitHubEventType !== 'push') {
      console.log('Not a push event. retruning.');
      return false;
    }

    const refSplit = githubEventPayload.ref.split('/');
    const stage = refSplit[refSplit.length - 1];
    const commitMessage = githubEventPayload.head_commit.message;
    const project = githubEventPayload.repository.name;
    let gitHubOrg = '';
    if (githubEventPayload.repository.owner !== undefined) {
      gitHubOrg = githubEventPayload.repository.owner.name;
    }

    // a keptn-config-change message follows the following format:
    // [keptn-config-change]:<service>:<image-tag>
    if (commitMessage.indexOf('[keptn-config-change]') !== 0) {
      console.log('Not a keptn config change.');
      return false;
    }
    const commitMsgSplit = commitMessage.split(':');
    if (commitMsgSplit.length < 3) {
      return false;
    }

    const configChangeEvent: ConfigurationModel = {} as ConfigurationModel;

    configChangeEvent.service = commitMsgSplit[1];
    configChangeEvent.image = commitMsgSplit[2];
    configChangeEvent.project = project;
    configChangeEvent.stage = stage;
    configChangeEvent.gitHubOrg = gitHubOrg;

    console.log(`Sending ConfigChange event ${JSON.stringify(configChangeEvent)}`);

    const keptnEvent: KeptnRequestModel = new KeptnRequestModel();
    keptnEvent.data = configChangeEvent;
    keptnEvent.type = KeptnRequestModel.EVENT_TYPES.CONFIGURATION_CHANGED;
    await axios.post('http://event-broker.keptn.svc.cluster.local/keptn', keptnEvent);
    return true;
  }
}
