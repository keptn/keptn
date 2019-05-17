import { inject, injectable } from 'inversify';
import axios from 'axios';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';
import { Logger } from '../lib/log/Logger';
import moment from 'moment';

@injectable()
export class MessageService {

  constructor(@inject('ChannelReconciler') private readonly channelReconciler: ChannelReconciler) {}

  public async sendMessage(
    message: KeptnRequestModel,
    keptnContext: string = ''): Promise<boolean> {
    let channelUri;
    const eventType = message.type;
    if (eventType !== undefined) {
      const split = message.type.split('.');
      if (split.length < 4) {
        return false;
      }
      const channelName = split[3];
      channelUri = await this.channelReconciler.resolveChannel(channelName);
    }
    message.time = moment().format();
    if (channelUri === '' || channelUri === undefined) {
      Logger.error(keptnContext, `Could not find channel URI for event of type ${message.type}`);
      return false;
    }
    Logger.info(keptnContext, `Sending message to ${channelUri}`);

    try {
      const result = await axios.post(`http://${channelUri}`, message);
      Logger.debug(
        message.shkeptncontext,
        `Sent request to channel. Response: ${JSON.stringify(result.data)}`,
      );
    } catch (e) {
      Logger.error(message.shkeptncontext, `Error while sending request: ${e}`);
      return false;
    }

    return true;
  }
}
