import { injectable } from 'inversify';
import * as Jwt from 'jsonwebtoken';
import * as UUID from 'uuid';
import { CredentialsService } from './CredentialsService';
import { WebSocketChannelInfo } from '../lib/types/WebSocketChannelInfo';
import * as WebSocket from 'ws';
import { MessageQueue } from '../lib/types/MessageQueue';
import { Logger } from '../lib/log/Logger';

export class WebSocketService {
  private static instance: WebSocketService;
  private credentialsService: CredentialsService;
  private static messageQueues: MessageQueue[] = [];
  private static connections: any[] = [];
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
    this.handleCLIClientConnection(req, ws);
    ws.on('message', (message) => {
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
      Logger.debug(keptnContext, `received: ${message}`);
      Logger.debug(keptnContext, `Trying to find client for channel ${keptnContext}.`);
      let clientFound = false;
      const connection =
        WebSocketService.connections.find(connection => connection.channelId === keptnContext);
      if (connection !== undefined && connection.client !== undefined) {
        clientFound = true;
        Logger.debug(keptnContext, 'found client and sending message');
        connection.client.send(`${message}`);
      }

      if (!clientFound && msgPayload !== undefined) {
        Logger.debug(
          keptnContext,
          `No client for channel ${keptnContext} found. Putting msg into queue`,
        );
        const index =
          WebSocketService.messageQueues.findIndex(queue => queue.channelId === keptnContext);
        if (index > -1) {
          Logger.debug(keptnContext, `Queue found for ${keptnContext}.`);
          WebSocketService.messageQueues[index].messages.push(message);
        } else {
          Logger.debug(keptnContext, `No queue found for ${keptnContext}.`);
        }
      }
    });

    ws.on('close', () => {
      const conIdx = WebSocketService.connections.findIndex(connection => connection.client === ws);
      WebSocketService.connections.splice(conIdx, 1);
    });
  }

  private handleCLIClientConnection(req: any, ws: any) {
    const channelId = req.headers['x-keptn-ws-channel-id'];
    Logger.debug(channelId, `New connection with channelId ${channelId}`);
    if (channelId !== undefined) {
      WebSocketService.connections.push({
        channelId,
        client: ws,
      });
    }

    // check if there are already some messages in the buffer for this channel
    // (channel == CLI command execution)
    const index = WebSocketService.messageQueues
      .findIndex(queue => queue.channelId === channelId);
    if (index > -1) {
      WebSocketService.messageQueues[index].messages.forEach((msg) => {
        Logger.debug(channelId, `sending message from queue: ${msg}`);
        ws.send(`${msg}`);
      });
      WebSocketService.messageQueues.splice(index, 1);
    }
  }

  /**
   * Creates a new WebSocket Channel for logging a CLI command execution via WebSocket connection
   * Incoming connections for this channel will be validated using JWT
   *
   * @param keptnContext the ID used to identify the channel
   */
  public async createChannel(keptnContext: string): Promise<WebSocketChannelInfo> {
    const channelInfo: WebSocketChannelInfo = {} as WebSocketChannelInfo;
    const channelId = keptnContext;

    const apiToken = await this.credentialsService.getKeptnApiToken();
    if (apiToken === undefined || apiToken === '') {
      Logger.error(keptnContext, 'Could not establish WebSocket connection.');
      return channelInfo;
    }
    // create the JWT
    const token = Jwt.sign(
      { channelId },
      apiToken,
      {
        expiresIn: 1 * 24 * 60 * 60 * 1000,
      },
    );
    channelInfo.channelId = channelId;
    channelInfo.token = token;

    // create a new message queue
    // usually, the service executing the command will be connected and send logs before
    // the CLI client is connected. Until the CLI has established a connection, we need to buffer
    // the log messages sent by the service
    const messageQueue = {} as MessageQueue;
    messageQueue.channelId = channelId;
    messageQueue.messages = [];
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
