import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { DynatraceRequestModel } from './DynatraceRequestModel';
import { KeptnRequestModel } from '../lib/types/KeptnRequestModel';

const uuidv4 = require('uuid/v4');

@injectable()
export class DynatraceService {

  constructor() {}

  public async handleDynatraceEvent(dtEventPayload: DynatraceRequestModel): Promise<boolean> {

    if (dtEventPayload.shkeptncontext === undefined) {
      dtEventPayload.shkeptncontext = uuidv4();
    }
    console.log(JSON.stringify({
      keptnContext: dtEventPayload.shkeptncontext,
      keptnService: 'eventbroker-ext',
      logLevel: 'INFO',
      message: `Sending keptn event ${JSON.stringify(dtEventPayload)}`,
    }));
    await axios.post('http://event-broker.keptn.svc.cluster.local/keptn', dtEventPayload);
    return true;
  }
}
