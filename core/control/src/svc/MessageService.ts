import { injectable } from 'inversify';
import axios from 'axios';

@injectable()
export class MessageService {
  private channelUri: string;
  constructor() {
    this.channelUri = process.env.CHANNEL_URI || '';
  }

  public async sendMessage(message: any): Promise<boolean> {
    if (this.channelUri === '') {
      return false;
    }
    console.log(`Forwarding message to channel ${this.channelUri}`);
    const result = await axios.post(this.channelUri, message);
    console.log(result);
    return true;
  }
}
