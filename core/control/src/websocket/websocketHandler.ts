import express = require('express');
import * as WebSocket from 'ws';

const websocketHandler = async (
  ws: WebSocket.Server,
  request: express.Request,
) => {
  ws.on('message', (msg) => {
    console.log(msg);
  });
  console.log('socket', request);
};

export = websocketHandler;
