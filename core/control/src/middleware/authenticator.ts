import express = require('express');
import { AuthRequest } from '../lib/types/AuthRequest';
import axios from 'axios';

const AUTH_URL = process.env.NODE_ENV === 'production' ?
  'http://authenticator.keptn.svc.cluster.local/auth' : 'http://localhost:3000/auth';

const authenticator: express.RequestHandler = (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  console.log('Starting authentication');
  console.log(JSON.stringify(request.body));
  // TODO: insert call to authenticator.keptn.svc.cluster.local here
  // get signature from header
  const signature: string = request.headers['x-keptn-signature'] as string;
  console.log(signature);
  if (signature === undefined) {
    response.status(401);
    return;
  }
  const payload = JSON.stringify(request.body);

  const authRequest: AuthRequest = {
    signature,
    payload,
  };

  console.log(`Sending auth request: ${JSON.stringify(authRequest)}`);

  axios.post(AUTH_URL, authRequest)
    .then((authResult) => {
      console.log(authResult);
      if (authResult.data.authenticated) {
        next();
      } else {
        response.status(401);
        response.end();
      }
    })
    .catch(() => {
      console.log('Authentication request failed');
      response.status(401);
    });
};

export = authenticator;
