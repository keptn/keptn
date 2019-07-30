import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { Service } from './Service';
import { RequestModel } from './RequestModel';

describe('Service', () => {
  let service: Service;
  beforeEach(() => {
    cleanUpMetadata();
  });
  it('should ...', async () => {
    expect(true).is.true;
  });
});
