import express = require('express');
import { AuthRequest } from '../lib/types/AuthRequest';
import { BearerAuthRequest } from '../lib/types/BearerAuthRequest';
import axios from 'axios';

const AUTH_URL = 'http://authenticator.keptn.svc.cluster.local';

const authenticator: express.RequestHandler = async (
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) => {
  // TODO: also add authentication for Dynatrace endpoint!
  if (request.url !== undefined && request.url.indexOf('dynatrace') > 0) {
    await handleDynatraceRequest(request, response, next);
    return;
  } else if (request.url !== undefined && request.url.indexOf('event') > 0) {
    await handleExtEventRequest(request, response, next);
    return;
  }
  await handleGitHubRequest(request, response, next);
};

async function handleExtEventRequest(
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) {
  console.log('Starting authentication');
  console.log(JSON.stringify(request.body));
  const signature: string =
    request.headers !== undefined ?
      request.headers['x-keptn-signature'] as string : undefined;
  console.log(signature);
  if (signature === undefined) {
    response.status(401);
    response.end();
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
    response.end();
  }
}

async function handleDynatraceRequest(
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) {
  let token;
  const authHeaderValue = request.headers['authorization'];
  if (authHeaderValue === undefined || authHeaderValue === '') {
    response.status(401);
    response.end();
    return;
  }
  const split = authHeaderValue.split(' ');
  if (split.length < 2) {
    response.status(401);
    response.end();
    return;
  }
  token = split[1];
  const authRequest: BearerAuthRequest = {
    token,
  };

  console.log(`Sending Auth request: ${JSON.stringify(authRequest)}`);
  let authResult;
  try {
    authResult = await axios.post(`${AUTH_URL}/auth/token`, authRequest);
    console.log(`auth result: ${JSON.stringify(authResult.data)}`);
    if (authResult.data.authenticated) {
      next();
    } else {
      response.status(401);
      response.end();
    }
  } catch (e) {
    console.log(`Authentication request failed: ${e}`);
    response.status(401);
    response.end();
  }
}

async function handleGitHubRequest(
  request: express.Request,
  response: express.Response,
  next: express.NextFunction,
) {
  console.log('Starting authentication');
  console.log(JSON.stringify(request.body));
  // get signature from header
  const signature: string =
    request.headers !== undefined ?
      request.headers['x-hub-signature'] as string : undefined;
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

  console.log(`Sending Auth request: ${JSON.stringify(authRequest)}`);
  let authResult;
  try {
    authResult = await axios.post(`${AUTH_URL}/auth`, authRequest);
    console.log(`auth result: ${JSON.stringify(authResult.data)}`);
    if (authResult.data.authenticated) {
      next();
    } else {
      response.status(401);
      response.end();
    }
  } catch (e) {
    console.log(`Authentication request failed: ${e}`);
    response.status(401);
    response.end();
  }
}

export = authenticator;
