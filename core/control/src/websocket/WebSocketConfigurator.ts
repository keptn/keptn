import * as Jwt from 'jsonwebtoken';
import * as express from 'express';
import websocketHandler = require('../websocket/websocketHandler');
import { WebSocketService } from '../svc/WebSocketService';

export class WebSocketConfigurator {
  private app: any;
  static instance: WebSocketConfigurator;

  private constructor(app: any) {
    this.app = app;
  }

  static getInstance(app) {
    if (WebSocketConfigurator.instance === undefined) {
      WebSocketConfigurator.instance = new WebSocketConfigurator(app);
    }
    return WebSocketConfigurator.instance;
  }

  configure() {
    require('express-ws')(this.app, undefined, {
      wsOptions: {
        verifyClient: WebSocketService.getInstance().verifyToken,
      },
    });
    this.app.ws('/comm', websocketHandler);
  }
}
