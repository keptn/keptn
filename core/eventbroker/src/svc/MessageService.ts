import { inject, injectable } from 'inversify';
import axios from 'axios';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';

@injectable()
export class MessageService {

  constructor(@inject('ChannelReconciler') private readonly channelReconciler: ChannelReconciler) {}

  public async sendMessage(message: KeptnRequestModel): Promise<boolean> {
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
    if (channelUri === '' || channelUri === undefined) {
      console.log(`Could not find channel URI for event of type ${message.type}`);
      return false;
    }
    console.log(`Sending message to ${channelUri}`);

    axios.post(`http://${channelUri}`, message).then().catch((e) => { console.log(e); });

    return true;
  }
}
