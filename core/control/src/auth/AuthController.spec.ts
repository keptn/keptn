import 'reflect-metadata';
import { AuthController } from './AuthController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';

describe('AuthController', () => {
  let authController: AuthController;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    authController = new AuthController();
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should return true if a message has been forwarded', async () => {
    const messageServiceStub = sinon
      .stub()
      .returns(true);

    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    await authController.verifyToken(request, response, next);

    expect(responseSendSpy.calledWith({ status: 'OK' })).is.true;
  });
});
