import 'reflect-metadata';
import { AuthController } from './AuthController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { AuthService } from './AuthService';
import { cleanUpMetadata } from 'inversify-express-utils';

describe('AuthController', () => {

  let authService: AuthService;
  let authController: AuthController;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    authService = new AuthService();
    authController = new AuthController(authService);
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });

  it('Should return a auth-true response for valid requests', async () => {
    const verifyStub = sinon.stub().returns(true);

    authService.verify = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      payload: '2344',
      signature: 'sha1=8a12fb3402b3d203aee37701f39e22a104b9a2a0',
    };

    await authController.authenticate(request, response, next);

    expect(verifyStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      authenticated: true,
    })).is.true;
  });

  it('Should return a auth-false response for valid requests', async () => {
    const verifyStub = sinon.stub().returns(false);

    authService.verify = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      payload: '2344',
      signature: 'sha1=invalid',
    };

    await authController.authenticate(request, response, next);

    expect(verifyStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      authenticated: false,
    })).is.true;
  });

  it('Should return 422 response for invalid payload', async () => {
    const verifyStub = sinon.fake();
    verifyStub.returnValues = [false];
    authService.verify = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    request.body = {
      foo: 'bar',
    };

    await authController.authenticate(request, response, next);

    expect(responseStatusSpy.calledWith(422)).is.true;
    expect(verifyStub.called).is.false;
  });

  it('Should return 422 response for invalid payload', async () => {
    const verifyStub = sinon.fake();
    verifyStub.returnValues = [false];

    authService.verify = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    request.body = {
      signature: {},
      payload: '',
    };

    await authController.authenticate(request, response, next);

    expect(responseStatusSpy.calledWith(422)).is.true;
    expect(verifyStub.called).is.false;
  });
  // BEARER TOKEN TESTS
  it('Should return a auth-true response for valid bearer token requests', async () => {
    const verifyStub = sinon.stub().returns(true);

    authService.verifyBearerToken = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      token: 'avalidtoken',
    };

    await authController.authenticateToken(request, response, next);

    expect(verifyStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      authenticated: true,
    })).is.true;
  });

  it('Should return a auth-false response for invalid bearer token requests', async () => {
    const verifyStub = sinon.stub().returns(false);

    authService.verifyBearerToken = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      token: 'invalidtoken',
    };

    await authController.authenticateToken(request, response, next);

    expect(verifyStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      authenticated: false,
    })).is.true;
  });

  it('Should return 422 response for invalid bearer token payload', async () => {
    const verifyStub = sinon.fake();
    verifyStub.returnValues = [false];
    authService.verifyBearerToken = verifyStub;

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    request.body = {
      foo: 'bar',
    };

    await authController.authenticateToken(request, response, next);

    expect(responseStatusSpy.calledWith(422)).is.true;
    expect(verifyStub.called).is.false;
  });
});
