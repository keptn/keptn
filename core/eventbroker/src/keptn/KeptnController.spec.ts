import 'reflect-metadata';
import { KeptnController } from './KeptnController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { MessageService } from '../svc/MessageService';
import { cleanUpMetadata } from 'inversify-express-utils';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';

describe('KeptnController', () => {
  let keptnController: KeptnController;
  let messageService: MessageService;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    messageService = new MessageService(new ChannelReconciler);
    keptnController = new KeptnController(messageService);
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
    await keptnController.dispatchEvent(request, response, next);

    expect(messageServiceStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({
      success: true,
    })).is.true;
  });
});
