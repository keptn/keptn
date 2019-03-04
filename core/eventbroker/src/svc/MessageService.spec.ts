import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { MessageService } from './MessageService';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
import { ChannelReconciler } from '../lib/channel/ChannelReconciler';
const nock = require('nock');

describe('MessageService', () => {
  let messageService: MessageService;
  let channelReconciler: ChannelReconciler;
  beforeEach(() => {
    cleanUpMetadata();
    process.env.CHANNEL_URI = 'channel';
    channelReconciler = new ChannelReconciler();
    messageService = new MessageService(channelReconciler);
  });
  it('should return true if a message has been forwarded', async () => {
    const message: KeptnRequestModel = {} as KeptnRequestModel;
    message.type = 'sh.keptn.events.new-artefact';

    const channelReconcilerResolveStub = sinon.stub().resolves('any-url');

    channelReconciler.resolveChannel = channelReconcilerResolveStub;

    nock(`http://any-url`, {
      filteringScope: () => {
        return true;
      },
    })
      .post('/', message)
      .reply(200, {});
    const result = await messageService.sendMessage(message);
    expect(result).to.be.true;
  });
  /*
  it('should return false if a message has not been forwarded', async () => {
    const message: KeptnRequestModel = {} as KeptnRequestModel;
    nock(`http://any-url`, {
      filteringScope: () => {
        return true;
      },
    })
      .post('/', message)
      .reply(503, {});
    const result = await messageService.sendMessage(message);
    expect(result).to.be.true;
  });
  */
  it('should return false if no channel uri can be found ', async () => {
    const message: KeptnRequestModel = {} as KeptnRequestModel;
    message.type = 'unknown';

    const channelReconcilerResolveStub = sinon.stub().resolves('');

    channelReconciler.resolveChannel = channelReconcilerResolveStub;

    const result = await messageService.sendMessage(message);
    expect(result).to.be.false;
  });
  it('should return false if no channel uri can be found 2', async () => {
    const message: KeptnRequestModel = {} as KeptnRequestModel;
    message.type = 'unknown.keptn.something.notvalid';

    const channelReconcilerResolveStub = sinon.stub().resolves('');

    channelReconciler.resolveChannel = channelReconcilerResolveStub;

    const result = await messageService.sendMessage(message);
    expect(result).to.be.false;
  });
});
