import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { MessageService } from './MessageService';
import { KeptnRequestModel } from '../keptn/KeptnRequestModel';
const nock = require('nock');

describe('MessageService', () => {
  let messageService: MessageService;
  beforeEach(() => {
    cleanUpMetadata();
    process.env.CHANNEL_URI = 'channel';
    messageService = new MessageService();
  });
  it('should return true if a message has been forwarded', async () => {

    const message: KeptnRequestModel = {} as KeptnRequestModel;
    message.type = 'sh.keptn.events.new-artefact';
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
    messageService = new MessageService();
    const message: KeptnRequestModel = {} as KeptnRequestModel;
    message.type = 'unknown';

    const result = await messageService.sendMessage(message);
    expect(result).to.be.false;
  });
});
