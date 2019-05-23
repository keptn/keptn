import 'reflect-metadata';
import { ExtEventController } from './ExtEventController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { ExtEventService } from './ExtEventService';
import { cleanUpMetadata } from 'inversify-express-utils';
import { doesNotReject } from 'assert';

describe('ExtEventController', () => {
  let extEventController: ExtEventController;
  let extEventService: ExtEventService;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    extEventService = new ExtEventService();
    extEventController = new ExtEventController(extEventService);
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should invoke the message service', async () => {
    const extEventServiceStub = sinon
      .stub()
      .returns(true);

    extEventService.handleExtEvent = extEventServiceStub;
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

    expect(extEventServiceStub.calledWith(request.body)).is.true;
  }).timeout(5000);
});
