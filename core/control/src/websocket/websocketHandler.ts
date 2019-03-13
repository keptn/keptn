import express = require('express');
import * as WebSocket from 'ws';

export class WebSocketHandler {

  private wssInstance: any;
  constructor(wssInstance) {
    this.wssInstance = wssInstance;
  }

  handleMessage(
    ws: WebSocket.Server,
    request: express.Request,
  ) {
    console.log(this);
    const wssInstance = this.wssInstance;
    ws.on('message', (msg) => {
      console.log('Received message');
      wssInstance.clients.forEach((client) => {
        client.send(`${msg}`);
      });
    });
    console.log('socket', request);
  }
}

/*
const websocketHandler = async (
  ws: WebSocket.Server,
  request: express.Request,
) => {
  ws.on('message', (msg) => {
    console.log(msg);
    console.log(ws);
    ws.clients.forEach((client) => {
      client.send(`${msg}`);
    });
  });
  console.log('socket', request);
};

export = websocketHandler;
*/
