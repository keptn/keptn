import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { AuthRequest } from '../lib/types/AuthRequest';
import { BearerAuthRequest } from '../lib/types/BearerAuthRequest';
const nock = require('nock');

const authenticator = require('./authenticator');

const AUTH_URL = 'http://authenticator.keptn.svc.cluster.local';

describe('authenticator', () => {
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should call next() if authentication was successful', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.headers = {};
    request.headers['x-hub-signature'] = 'sha1=123';

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    const authRequest: AuthRequest = {} as AuthRequest;
    authRequest.signature = request.headers['x-hub-signature'] as string;
    authRequest.payload = JSON.stringify(request.body);

    nock(AUTH_URL, {
      filteringScope: () => {
        return true;
      },
    })
      .post('/auth')
      .reply(200, { authenticated: true });

    await authenticator(request, response, next);

    expect(nextSpy.called).is.true;
  });

  it('should call next() if Dynatrace token authentication was successful', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.url = '/dynatrace';
    request.headers = {};
    request.headers['authorization'] = 'Bearer mytokenvalue';

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    const authRequest: BearerAuthRequest = {} as BearerAuthRequest;
    authRequest.token = 'mytokenvalue';

    nock(AUTH_URL, {
      filteringScope: () => {
        return true;
      },
    })
      .post('/auth/token')
      .reply(200, { authenticated: true });

    await authenticator(request, response, next);

    expect(nextSpy.called).is.true;
  });

  it('should return a 401 if no signature header was provided', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
  });

  it('should return a 401 if no auth token for /dynatrace was provided', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.body = {};
    request.url = '/dynatrace';
    request.headers = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
  });

  it('should return a 401 if no valid auth token for /dynatrace was provided', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.body = {};
    request.url = '/dynatrace';
    request.headers = {};
    request.headers['authorization'] = 'Bearer';

    const nextSpy = sinon.spy();
    next = nextSpy;

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
  });

  it('should return a 401 if the auth service could not verify the signature', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.headers = {};
    request.headers['x-keptn-signature'] = 'sha1=123';

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    const authRequest: AuthRequest = {} as AuthRequest;
    authRequest.signature = request.headers['x-hub-signature'] as string;
    authRequest.payload = JSON.stringify(request.body);

    nock(AUTH_URL)
      .post('/auth')
      .reply(200, { authenticated: false });

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
    expect(nextSpy.called).is.false;
  });

  it('should return a 401 if the auth service could not verify the token', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.url = '/dynatrace';
    request.headers = {};
    request.headers['authorization'] = 'Bearer invalidtoken';

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    const authRequest: BearerAuthRequest = {} as BearerAuthRequest;
    authRequest.token = 'mytokenvalue';

    nock(AUTH_URL)
      .post('/auth/token')
      .reply(200, { authenticated: false });

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
    expect(nextSpy.called).is.false;
  });

  it('should not try to verify requests for getting the swagger doc', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.body = {};

    request.url = '/swagger.json';

    const nextSpy = sinon.spy();
    next = nextSpy;

    await authenticator(request, response, next);

    expect(nextSpy.called).is.true;
  });

  it('should return a 401 if the auth service call fails', async () => {
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    const responseEndSpy = sinon.spy();
    response.end = responseEndSpy;

    request.headers = {};
    request.headers['x-hub-signature'] = 'sha1=123';

    request.body = {};

    const nextSpy = sinon.spy();
    next = nextSpy;

    const authRequest: AuthRequest = {} as AuthRequest;
    authRequest.signature = request.headers['x-keptn-signature'] as string;
    authRequest.payload = JSON.stringify(request.body);

    nock(AUTH_URL)
      .post('/auth')
      .reply(500);

    await authenticator(request, response, next);

    expect(responseStatusSpy.calledWith(401)).is.true;
    expect(nextSpy.called).is.false;
  });
});
