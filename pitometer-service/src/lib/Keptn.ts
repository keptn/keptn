import { RequestModel } from '../svc/RequestModel';
import axios, { AxiosRequestConfig, AxiosPromise, AxiosError } from 'axios';
import { Logger } from './Logger';

export class Keptn {
  static async sendEvent(event: RequestModel): Promise<void> {
    if (!(process.env.NODE_ENV === 'production')) {
      return;
    }
    try {
      const headers = {
        'Content-Type': 'application/cloudevents+json',
      };

      await axios.post(
        'http://event-broker.keptn.svc.cluster.local/keptn',
        event,
        { headers },
      );
    } catch (e) {
      Logger.log(event.shkeptncontext, 'Could not send event to event-broker', 'ERROR');
    }
  }
}
