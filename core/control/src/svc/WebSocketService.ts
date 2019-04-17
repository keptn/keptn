import { injectable } from 'inversify';
import * as Jwt from 'jsonwebtoken';
import * as UUID from 'uuid';
import { CredentialsService } from './CredentialsService';
import { WebSocketChannelInfo } from '../lib/types/WebSocketChannelInfo';
import * as WebSocket from 'ws';
import { MessageQueue } from '../lib/types/MessageQueue';

export class WebSocketService {
  private static instance: WebSocketService;
  private credentialsService: CredentialsService;
  private static messageQueues: MessageQueue[] = [];
  private static connections: any[];
  private constructor() {
    this.credentialsService = CredentialsService.getInstance();
    this.verifyToken = this.verifyToken.bind(this);
  }

  static getInstance(): WebSocketService {
    if (WebSocketService.instance === undefined) {
      WebSocketService.instance = new WebSocketService();
    }
    return WebSocketService.instance;
  }

  public handleConnection(wss, ws, req): void {
    const channelId = req.headers['x-keptn-ws-channel-id'];
    WebSocketService.connections.push({
      channelId,
      client: ws,
    });
    const index = WebSocketService.messageQueues
      .findIndex(queue => queue.channelId === channelId);

    if (index > -1) {
      WebSocketService.messageQueues[index].messages.forEach((msg) => {
        ws.send(`${msg}`);
      });
      WebSocketService.messageQueues.splice(index, 1);
    }
    ws.on('message', (message) => {
      console.log('received: %s', message);
      let msgPayload = undefined;
      try {
        msgPayload = JSON.parse(message);
      } catch (e) {
        msgPayload = undefined;
      }

      let keptnContext = '';
      if (msgPayload !== undefined) {
        keptnContext = msgPayload.shkeptncontext;
      }
      let clientFound = false;
      const connection =
        WebSocketService.connections.find(connection => connection.channelId === keptnContext);
      if (connection !== undefined && connection.client !== undefined) {
        clientFound = true;
        connection.client.send(`${message}`);
      }
      /*
      wss.clients.forEach((client: WebSocket) => {
        if (client.readyState === WebSocket.OPEN && ws !== client) {
          clientFound = true;
          client.send(`${message}`);
        }
      });
      */
      if (!clientFound && msgPayload) {
        const index =
          WebSocketService.messageQueues.findIndex(queue => queue.channelId === keptnContext);
        if (index > -1) {
          WebSocketService.messageQueues[index].messages.push(message);
        }
      }
    });
  }

  public async createChannel(keptnContext: string): Promise<WebSocketChannelInfo> {
    const channelInfo: WebSocketChannelInfo = {} as WebSocketChannelInfo;
    const channelId = keptnContext;
    const token = Jwt.sign(
      { channelId },
      await this.credentialsService.getKeptnApiToken(),
      {
        expiresIn: 1 * 24 * 60 * 60 * 1000,
      },
    );
    channelInfo.channelId = channelId;
    channelInfo.token = token;
    const messageQueue = {} as MessageQueue;
    messageQueue.channelId = channelId;
    WebSocketService.messageQueues.push(messageQueue);
    return channelInfo;
  }

  public async verifyToken(info, cb): Promise<void> {
    if (process.env.NODE_ENV !== 'production') {
      console.log('Skipping verification in dev mode');
      cb(true);
      return;
    }
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
  }
}
