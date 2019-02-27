import express = require('express');
import { AuthRequest } from '../lib/types/AuthRequest';
import axios from 'axios';

const AUTH_URL = 'http://authenticator.keptn.svc.cluster.local/auth';

const authenticator: express.RequestHandler = async (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  if (request.url !== undefined && request.url.indexOf('swagger') > 0) {
    console.log('Skipping auth for swagger doc');
    next();
    return;
  }
  console.log('Starting authentication');
  console.log(JSON.stringify(request.body));
  // TODO: insert call to authenticator.keptn.svc.cluster.local here
  // get signature from header
  const signature: string =
    request.headers !== undefined ?
      request.headers['x-keptn-signature'] as string : undefined;
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
  let authResult;
  try {
    authResult = await axios.post(AUTH_URL, authRequest);
    if (authResult.data.authenticated) {
      next();
    } else {
      response.status(401);
      response.end();
    }
  } catch (e) {
    console.log('Authentication request failed');
    response.status(401);
  }
};

export = authenticator;
