import * as Jwt from 'jsonwebtoken';
import * as express from 'express';
import websocketHandler = require('../websocket/websocketHandler');
import { WebSocketService } from '../svc/WebSocketService';

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
    require('express-ws')(this.app, this.server, {
      wsOptions: {
        verifyClient: WebSocketService.getInstance().verifyToken,
      },
    });
    this.app.ws('/comm', websocketHandler);
  }
}
