import { injectable } from 'inversify';
import * as Jwt from 'jsonwebtoken';
import * as UUID from 'uuid';
import { CredentialsService } from './CredentialsService';
import { WebSocketChannelInfo } from '../lib/types/WebSocketChannelInfo';
import * as WebSocket from 'ws';

export class WebSocketService {
  private static instance: WebSocketService;
  private credentialsService: CredentialsService;
  private constructor() {
    this.credentialsService = CredentialsService.getInstance();
  }

  static getInstance(): WebSocketService {
    if (WebSocketService.instance === undefined) {
      WebSocketService.instance = new WebSocketService();
    }
    return WebSocketService.instance;
  }

  public async createChannel(): Promise<WebSocketChannelInfo> {
    const channelInfo: WebSocketChannelInfo = {} as WebSocketChannelInfo;
    const channelId = UUID.v4();
    const token = Jwt.sign(
      { channelId },
      await this.credentialsService.getKeptnApiToken(),
      {
        expiresIn: 1 * 24 * 60 * 60 * 1000,
      },
    );
    channelInfo.channelId = channelId;
    channelInfo.token = token;
    return channelInfo;
  }

  public async verifyToken(info, cb) {
    cb(true);
    /*
    const token = info.req.headers.token;
    const secretKey = await this.credentialsService.getKeptnApiToken();
    if (!token) {
      cb(false, 401, 'Unauthorized');
    } else {
      Jwt.verify(token, secretKey, (err, decoded) => {
        if (err) {
          cb(false, 401, 'Unauthorized');
        } else {
          info.req.user = decoded;
          cb(true);
        }
      });
    }
    */
  }
}
