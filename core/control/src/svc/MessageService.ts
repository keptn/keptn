import { injectable } from 'inversify';
import axios from 'axios';

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

    axios.post(`http://${this.channelUri}`, message).then().catch(() => {});

    return true;
  }
}
