import 'reflect-metadata';
import * as express from 'express';
import { expect } from 'chai';
import 'mocha';
import * as sinon from 'sinon';
import { cleanUpMetadata } from 'inversify-express-utils';
import { OrgToRepoMapper } from './OrgToRepoMapper';

describe('OrgToRepoMapper', () => {
  let orgToRepoMapper: OrgToRepoMapper;
  beforeEach(() => {
    cleanUpMetadata();
    orgToRepoMapper = new OrgToRepoMapper();
  });
  it('should return a mapping', async () => {
    orgToRepoMapper.getRepoForOrg('org');
  });
});
