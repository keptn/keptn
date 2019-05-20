import { injectable } from 'inversify';
import axios from 'axios';
import { Logger } from '../lib/log/Logger';

@injectable()
export class MessageService {
  private channelUri: string;
  constructor() {
    this.channelUri = process.env.CHANNEL_URI || '';
  }

  public async sendMessage(message: any): Promise<boolean> {
    console.log(`Forwarding message to channel ${this.channelUri}`);
    if (this.channelUri === '') {
      return false;
    }

    try {
      const result = await axios.post(`http://${this.channelUri}`, message);
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
