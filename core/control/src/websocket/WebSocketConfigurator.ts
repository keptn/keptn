import * as Jwt from 'jsonwebtoken';
import * as express from 'express';
const WebSocket = require('ws');
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
    const server = this.server;
    const wss = new WebSocket.Server({ server });
    wss.on('connection', function connection(ws) {
      ws.on('message', function incoming(message) {
        console.log('received: %s', message);
      });
      ws.send('something');
    });
    /*
    this.server.on('upgrade', (request, socket, head) => {
      wss.handleUpgrade(request, socket, head, function done(ws) {
        wss.emit('connection', ws, request);
      });
    });
    */
    /*
    require('express-ws')(this.app, this.server, {
      wsOptions: {
        verifyClient: WebSocketService.getInstance().verifyToken,
      },
    });
    this.app.ws('/comm', websocketHandler);
    */
  }
}
