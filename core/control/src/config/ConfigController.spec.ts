import 'reflect-metadata';
import { ConfigController } from './ConfigController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { MessageService } from '../svc/MessageService';
import { cleanUpMetadata } from 'inversify-express-utils';

describe('ConfigController', () => {
  let configController: ConfigController;
  let messageService: MessageService;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    messageService = new MessageService();
    configController = new ConfigController(messageService);
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should return true if a message has been forwarded', async () => {
    const messageServiceStub = sinon
      .stub()
      .returns(true);

    messageService.sendMessage = messageServiceStub;
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      data: {
        org: 'my_org',
        user: 'my_user',
        token: 'my_token',
      },
    };
    await configController.setGithubConfig(request, response, next);

    expect(messageServiceStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      success: true,
    })).is.true;
  });
  it('should return false if a message has not been forwarded', async () => {
    const messageServiceStub = sinon
      .stub()
      .returns(false);

    messageService.sendMessage = messageServiceStub;
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      data: {
        org: 'my_org',
        user: 'my_user',
        token: 'my_token',
      },
    };
    await configController.setGithubConfig(request, response, next);

    expect(messageServiceStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      success: false,
    })).is.true;
  });
});
