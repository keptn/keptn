import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { MessageService } from './MessageService';
const nock = require('nock');

describe('MessageService', () => {
  let messageService: MessageService;
  beforeEach(() => {
    cleanUpMetadata();
    process.env.CHANNEL_URI = 'channel';
    messageService = new MessageService();
  });
  it('should return true if a message has been forwarded', async () => {

    const message = {
      foo: 'bar',
    };
    nock(`http://${process.env.CHANNEL_URI}`)
      .post('/', message)
      .reply(200, {});
    const result = await messageService.sendMessage(message);
    expect(result).to.be.true;
  });
  /*
  it('should return false if a message has not been forwarded', async () => {
    const message = {
      foo: 'bar',
    };
    nock(`http://${process.env.CHANNEL_URI}`)
      .post('/', message)
      .reply(503, {});
    const result = await messageService.sendMessage(message);
    expect(result).to.be.false;
  });
  */
  it('should return false if no channel uri has been set', async () => {
    process.env.CHANNEL_URI = '';
    messageService = new MessageService();
    const message = {
      foo: 'bar',
    };

    const result = await messageService.sendMessage(message);
    expect(result).to.be.false;
  });
});
