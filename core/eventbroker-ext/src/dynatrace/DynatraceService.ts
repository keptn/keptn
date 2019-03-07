import 'reflect-metadata';
import { injectable } from 'inversify';
import axios  from 'axios';
import { DynatraceRequestModel } from './DynatraceRequestModel';
import { KeptnRequestModel } from '../lib/types/KeptnRequestModel';

@injectable()
export class DynatraceService {

  constructor() {}

  public async handleDynatraceEvent(dtEventPayload: DynatraceRequestModel): Promise<boolean> {

    const keptnEvent: KeptnRequestModel = new KeptnRequestModel();
    keptnEvent.data = dtEventPayload;
    keptnEvent.type = KeptnRequestModel.EVENT_TYPES.PROBLEM;
    console.log(`Sending keptn event ${JSON.stringify(keptnEvent)}`);
    await axios.post('http://event-broker.keptn.svc.cluster.local/keptn', keptnEvent);
    return true;
  }
}
