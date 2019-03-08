import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { DockerService } from './DockerService';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';
import { MessageService } from '../svc/MessageService';
import { DockerRequestModel } from './DockerRequestModel';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';

describe('DockerService', () => {
  let dockerService: DockerService;
  let channelReconciler: ChannelReconciler;
  let messageService: MessageService;
  beforeEach(() => {
    cleanUpMetadata();
    channelReconciler = new ChannelReconciler();
    messageService = new MessageService(channelReconciler);
    dockerService = new DockerService(messageService);
  });
  it('should return true if a message has been forwarded', async () => {
    const message: DockerRequestModel = {
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
            repository: 'library/keptn/keptn-event-broker-ext',
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

    const messageServiceSendMessageStub = sinon.stub().resolves(true);

    messageService.sendMessage = messageServiceSendMessageStub;

    const result = await dockerService.handleDockerRequest(message);

    const expectedMessage = new KeptnRequestModel();
    expectedMessage.data = {
      project: 'keptn',
      service: 'keptn-event-broker-ext',
      image: 'docker-registry.keptn.svc.cluster.local:5000/library/keptn/keptn-event-broker-ext',
      tag: 'latest',
    };
    expectedMessage.type = KeptnRequestModel.EVENT_TYPES.NEW_ARTEFACT;
    expect(messageServiceSendMessageStub.calledWithMatch(expectedMessage)).is.true;
    expect(result).to.be.true;
  });

  it(
    'should return true if a message has been forwarded (arbitrary amount of url components)',
    async () => {
      const message: DockerRequestModel = {
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

      const messageServiceSendMessageStub = sinon.stub().resolves(true);

      messageService.sendMessage = messageServiceSendMessageStub;

      const result = await dockerService.handleDockerRequest(message);

      const expectedMessage = new KeptnRequestModel();
      expectedMessage.data = {
        project: 'keptn',
        service: 'keptn-event-broker-ext',
        image: 'docker-registry.keptn.svc.cluster.local:5000/keptn/keptn-event-broker-ext',
        tag: 'latest',
      };
      expectedMessage.type = KeptnRequestModel.EVENT_TYPES.NEW_ARTEFACT;
      expect(messageServiceSendMessageStub.calledWithMatch(expectedMessage)).is.true;
      expect(result).to.be.true;
    });

  it('should return false if a message has not been forwarded', async () => {
    const message: DockerRequestModel = {
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

    const messageServiceSendMessageStub = sinon.stub().resolves(false);

    messageService.sendMessage = messageServiceSendMessageStub;

    const result = await dockerService.handleDockerRequest(message);
    expect(result).to.be.false;
  });

  it('should return false if a message that is not a push event is received', async () => {
    const message: DockerRequestModel = {
      events: [
        {
          id: 'a24e1fe3-efc9-42e3-b274-3c736f015552',
          timestamp: '2019-03-05T14:52:38.292839945Z',
          action: 'pull',
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

    const messageServiceSendMessageStub = sinon.stub().resolves(false);

    messageService.sendMessage = messageServiceSendMessageStub;

    const result = await dockerService.handleDockerRequest(message);
    expect(result).to.be.false;
  });

  it(
    'should return false if a message with an invalid format is received (empty events array)',
    async () => {
      const message: DockerRequestModel = {
        events: [
        ],
      };

      const messageServiceSendMessageStub = sinon.stub().resolves(false);

      messageService.sendMessage = messageServiceSendMessageStub;

      const result = await dockerService.handleDockerRequest(message);
      expect(result).to.be.false;
    });

  it(
    'should return false if a message with an invalid format is received (no target)',
    async () => {
      const message: any = {
        events: [
          {
            id: 'a24e1fe3-efc9-42e3-b274-3c736f015552',
            timestamp: '2019-03-05T14:52:38.292839945Z',
            action: 'push',
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

      const messageServiceSendMessageStub = sinon.stub().resolves(false);

      messageService.sendMessage = messageServiceSendMessageStub;

      const result = await dockerService.handleDockerRequest(message);
      expect(result).to.be.false;
    });

  it(
    'should return false if a message with an invalid format is \
      received (invalid repository string)',
    async () => {
      const message: any = {
        events: [
          {
            id: 'a24e1fe3-efc9-42e3-b274-3c736f015552',
            timestamp: '2019-03-05T14:52:38.292839945Z',
            action: 'push',
            request: {
              id: '164cafa5-ec03-4445-83b3-3ba6f7f28772',
              addr: '127.0.0.1:37402',
              host: 'docker-registry.keptn.svc.cluster.local:5000',
              method: 'PUT',
              useragent: 'kaniko/unset',
            },
            target: {
              mediaType: 'application/vnd.docker.distribution.manifest.v2+json',
              size: 2223,
              digest: 'sha256:d1b654481b04da5f1f69dc4e3bd72f4b592a60c6fb5618a9096eeac870cd3fe6',
              length: 2223,
              repository: 'keptn-event-broker-ext',
              url: 'http://docker-registry.keptn.svc.cluster.local:5000/',
              tag: 'latest',
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

      const messageServiceSendMessageStub = sinon.stub().resolves(false);

      messageService.sendMessage = messageServiceSendMessageStub;

      const result = await dockerService.handleDockerRequest(message);
      expect(result).to.be.false;
    });
});
