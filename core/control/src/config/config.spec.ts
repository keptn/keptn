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

  beforeEach(() => {
    cleanUpMetadata();
    messageService = new MessageService();
    configController = new ConfigController(messageService);
  });
  it('should return true if a message has been forwarded', async () => {
    // tslint:disable-next-line: prefer-const
    let request: express.Request = {} as express.Request;
    const response: express.Response = {} as express.Response;
    let next: express.NextFunction;

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
      status: 'OK',
    })).is.true;
  });
});
