import { expect } from 'chai';
import 'mocha';
import { DynatraceService } from './DynatraceService';
import nock from 'nock';

describe('DynatraceService', function () {
  this.timeout(0);
  let dynatraceService: DynatraceService;

  beforeEach(() => {
    dynatraceService = new DynatraceService();
  });

  it('Should ...', async () => {
    
  });
});
