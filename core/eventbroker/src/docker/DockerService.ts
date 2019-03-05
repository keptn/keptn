import 'reflect-metadata';
import { injectable, inject } from 'inversify';
import axios  from 'axios';
import { MessageService } from '../svc/MessageService';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';

@injectable()
export class DockerService {

  constructor(@inject('MessageService') private readonly messageService: MessageService) {}

  public async handleDockerRequest(event: any): Promise<boolean> {
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
    const project = repositorySplit[0];
    const service = repositorySplit[1];
    const tag = eventPayload.target.tag;
    const msgPayload = {
      project,
      service,
      tag,
    };

    const msg: KeptnRequestModel = new KeptnRequestModel();
    msg.data = msgPayload;
    await this.messageService.sendMessage(msg);
  }
}
