import 'reflect-metadata';
import { DockerController } from './DockerController';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { DockerService } from './DockerService';
import { cleanUpMetadata } from 'inversify-express-utils';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';
import { MessageService } from '../svc/MessageService';

describe('DockerController', () => {
  let dockerController: DockerController;
  let dockerService: DockerService;
  let request: express.Request;
  let response: express.Response;
  let next: express.NextFunction;

  beforeEach(() => {
    cleanUpMetadata();
    dockerService = new DockerService(
        new MessageService(
          new ChannelReconciler()));
    dockerController = new DockerController(dockerService);
    request = {} as express.Request;
    response = {} as express.Response;
    next = {} as express.NextFunction;
  });
  it('should return true if a message has been forwarded', async () => {
    const dockerServiceStub = sinon
      .stub()
      .returns(true);

    dockerService.handleDockerRequest = dockerServiceStub;
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    request.body = {
      events: [
        {
          id: 'a24e1fe3-efc9-42e3-b274-3c736f015552',
          timestamp: '2019-03-05T14:52:38.292839945Z',
          action: 'push',
          target: {
            mediaType: 'application/vnd.docker.distribution.manifest.v2+json',
            size: 2223,
            digest: 'sha256:d1b654481b04da5f1f69dc4e3bd72f4b592a60c6fb5618a9096eeac870cd3fe6',
            length: 2223,
            repository: 'keptn/keptn-event-broker-ext',
            url: 'http://docker-registry.keptn.svc.cluster.local:5000/',
            tag: 'latest',
          },
          request: {
            id: '164cafa5-ec03-4445-83b3-3ba6f7f28772',
            addr: '127.0.0.1:37402',
            host: 'docker-registry.keptn.svc.cluster.local:5000',
            method: 'PUT',
            useragent: 'kaniko/unset',
          },
          actor: {

          },
          source: {
            addr: 'docker-registry-55bd8d967c-hztw9:5000',
            instanceID: '2871c1e7-78b9-4fa5-b749-dc5e5fff8f9c',
          },
        },
      ],
    };

    await dockerController.handleDockerEvent(request, response, next);

    expect(dockerServiceStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({ messageSent: true })).is.true;
    expect(responseStatusSpy.calledWith(200)).is.true;
  });

  it('should return false if a message has not been forwarded', async () => {
    const dockerServiceStub = sinon
      .stub()
      .returns(false);

    dockerService.handleDockerRequest = dockerServiceStub;
    const responseSendSpy = sinon.spy();
    response.send = responseSendSpy;

    const responseStatusSpy = sinon.spy();
    response.status = responseStatusSpy;

    request.body = {
      events: [
        {
        },
      ],
    };

    await dockerController.handleDockerEvent(request, response, next);

    expect(dockerServiceStub.calledWith(request.body)).is.true;
    expect(responseSendSpy.calledWith({ messageSent: false })).is.true;
    expect(responseStatusSpy.calledWith(200)).is.true;
  });
});
