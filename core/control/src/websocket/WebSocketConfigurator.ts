import * as Jwt from 'jsonwebtoken';
import * as express from 'express';
const WebSocket = require('ws');

import { WebSocketService } from '../svc/WebSocketService';
import { WebSocketHandler } from './websocketHandler';

export class WebSocketConfigurator {
  private app: any;
  private server: any;
  static instance: WebSocketConfigurator;

  private constructor(app: any, server: any) {
    this.app = app;
    this.server = server;
  }

  static getInstance(app, server) {
    if (WebSocketConfigurator.instance === undefined) {
      WebSocketConfigurator.instance = new WebSocketConfigurator(app, server);
    }
    return WebSocketConfigurator.instance;
  }

  configure() {
    const server = this.server;
    const wss = new WebSocket.Server({
      server,
      verifyClient: WebSocketService.getInstance().verifyToken,
    });
    wss.on('connection', (ws, req) => {
      WebSocketService.getInstance().handleConnection(wss, ws, req);
    });
  }
}
