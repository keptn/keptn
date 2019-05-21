import 'reflect-metadata';
import { ExtEventController } from './ExtEventController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { MessageService } from '../svc/MessageService';
import { cleanUpMetadata } from 'inversify-express-utils';
import { doesNotReject } from 'assert';

describe('ExtEventController', () => {
  let extEventController: ExtEventController;
  let messageService: MessageService;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    messageService = new MessageService();
    extEventController = new ExtEventController(messageService);
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should invoke the message service', async () => {
    const messageServiceStub = sinon
      .stub()
      .returns(true);

    messageService.sendMessage = messageServiceStub;
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    request.body = {
      data: {
        project:  'my_prj', 
        service:  'my_svc',
        image:    'keptnexamples/carts:0.7.0',
        tag:      '0.7.0-8',
      },
    };
    await extEventController.handleExtEvent(request, response, next);

    expect(messageServiceStub.calledWith(request.body)).is.true;
  }).timeout(5000);
});
