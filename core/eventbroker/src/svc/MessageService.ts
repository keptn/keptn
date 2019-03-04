import { injectable } from 'inversify';
import axios from 'axios';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';

@injectable()
export class MessageService {

  private channelReconciler: ChannelReconciler;

  constructor() {
    this.channelReconciler = new ChannelReconciler();
  }

  public async sendMessage(message: KeptnRequestModel): Promise<boolean> {
    let channelUri;
    const eventType = message.type;
    if (eventType !== undefined) {
      const split = message.type.split('.');
      if (split.length < 4) {
        return false;
      }
      const channelName = split[3];
      channelUri = this.channelReconciler.resolveChannel(channelName);
    }
    if (channelUri === '') {
      console.log(`Could not find channel URI for event of type ${message.type}`);
      return false;
    }
    console.log(`Sending message to ${channelUri}`);

    axios.post(`http://${channelUri}`, message).then().catch((e) => { console.log(e); });

    return true;
  }
}
