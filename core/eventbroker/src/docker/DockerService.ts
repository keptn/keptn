import 'reflect-metadata';
import { injectable, inject } from 'inversify';
import axios  from 'axios';
import { MessageService } from '../svc/MessageService';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { DockerRequestModel } from './DockerRequestModel';

@injectable()
export class DockerService {

  constructor(@inject('MessageService') private readonly messageService: MessageService) {}

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
    const image = `${eventPayload.request.host}/${eventPayload.target.repository}`;
    const msgPayload = {
      project,
      service,
      image,
      tag,
    };

    const msg: KeptnRequestModel = new KeptnRequestModel();
    msg.data = msgPayload;
    msg.type = KeptnRequestModel.EVENT_TYPES.NEW_ARTEFACT;
    return await this.messageService.sendMessage(msg);
  }
}
