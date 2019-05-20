import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { ExtEventRequestModel } from './ExtEventRequestModel';
import { KeptnRequestModel } from '../lib/types/KeptnRequestModel';

const uuidv4 = require('uuid/v4');

@injectable()
export class ExtEventService {

  constructor() {}

  public async handleExtEvent(extEventPayload: ExtEventRequestModel): Promise<boolean> {

    if (extEventPayload.shkeptncontext === undefined) {
      extEventPayload.shkeptncontext = uuidv4();
    }
    console.log(JSON.stringify({
      keptnContext: extEventPayload.shkeptncontext,
      keptnService: 'eventbroker-ext',
      logLevel: 'INFO',
      message: `Sending keptn event ${JSON.stringify(extEventPayload)}`,
    }));
    await axios.post('http://event-broker.keptn.svc.cluster.local/keptn', extEventPayload);
    return true;
  }
}
