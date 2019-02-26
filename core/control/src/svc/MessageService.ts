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
    const result = await axios.post(this.channelUri, message);
    return true;
  }
}
