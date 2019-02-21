import express = require('express');
import { AuthRequest } from '../types/authRequest';
const crypto = require('crypto');
const bufferEq = require('buffer-equal-constant-time');

const authRouter = express.Router();

function sign(data: string) {
  const signature =
    `sha1=${crypto.createHmac('sha1', process.env.SECRET_TOKEN || '')
    .update(data).digest('hex')}`;

  console.log(`Calculated signature: ${signature}`);
  return signature;
}

function verify(authRequest: AuthRequest) {
  return bufferEq(Buffer.from(authRequest.signature), Buffer.from(sign(authRequest.payload)));
}

authRouter.post('/', (request: express.Request, response: express.Response) => {
  console.log('Starting authentication');
  const authRequest: AuthRequest = request.body;
  console.log(`Received auth request: ${JSON.stringify(authRequest)}`);

  const authResult = {
    authenticated: verify(authRequest),
  };

  console.log(`Response: ${JSON.stringify(authResult)}`);

  response.send(authResult);
});
// add more route handlers here
// e.g. authRouter.post('/', (req,res,next)=> {/*...*/})
export = authRouter;
