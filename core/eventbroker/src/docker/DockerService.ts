import 'reflect-metadata';
import { injectable, inject } from 'inversify';
import { MessageService } from '../svc/MessageService';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { DockerRequestModel } from './DockerRequestModel';
import { OrgToRepoMapper } from '../lib/org-to-repo-mapper/OrgToRepoMapper';
import { Logger } from '../lib/log/Logger';

@injectable()
export class DockerService {

  constructor(
    @inject('MessageService') private readonly messageService: MessageService,
    @inject('OrgToRepoMapper') private readonly orgToRepoMapper: OrgToRepoMapper,
  ) {}

  public async handleDockerRequest(event: DockerRequestModel): Promise<boolean> {
    if (event.events === undefined || event.events.length === 0) {
      return false;
    }
    const eventPayload = event.events[0];
    if (eventPayload.action !== 'push') {
      return false;
    }

    if (eventPayload.target === undefined || eventPayload.target.repository === undefined) {
      return false;
    }
    const repositorySplit = eventPayload.target.repository.split('/');
    if (repositorySplit.length < 2) {
      return false;
    }
    const project = repositorySplit[repositorySplit.length - 2];
    const service = repositorySplit[repositorySplit.length - 1];
    const tag = eventPayload.target.tag;

    if (tag === undefined) {
      return false;
    }
    const image = `${eventPayload.request.host}/${eventPayload.target.repository}`;

    const msg: KeptnRequestModel = new KeptnRequestModel();

    const repo = await this.orgToRepoMapper.getRepoForOrg(project);
    if (repo === '') {
      Logger.info(msg.shkeptncontext, `No repo found for organization ${project}`);
      return false;
    }
    Logger.info(msg.shkeptncontext, `Found mapping org ${project} to ${repo}`);
    Logger.info(
      msg.shkeptncontext,
      `Starting new pipeline run for ${project}/${service}:${tag}`,
      true,
    );

    const msgPayload = {
      service,
      image,
      tag,
      project: repo,
    };

    msg.data = msgPayload;
    msg.type = KeptnRequestModel.EVENT_TYPES.NEW_ARTEFACT;

    Logger.info(msg.shkeptncontext, JSON.stringify(msg));

    return await this.messageService.sendMessage(msg, msg.shkeptncontext);
  }
}
