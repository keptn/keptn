import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { ExtEventRequestModel } from './ExtEventRequestModel';

@injectable()
export class ExtEventService {

  constructor() {}

  public async handleExtEvent(extEventPayload: ExtEventRequestModel): Promise<boolean> {

    console.log(JSON.stringify({
      keptnContext: extEventPayload.shkeptncontext,
      keptnService: 'control',
      logLevel: 'INFO',
      message: `Sending keptn event ${JSON.stringify(extEventPayload)}`,
    }));
    await axios.post('http://event-broker.keptn.svc.cluster.local/keptn', extEventPayload);
    return true;
  }
}
