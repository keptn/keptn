import { Service } from './service';
import { waitForAsync } from '@angular/core/testing';

describe('Service', () => {
  it('should create instances from json', waitForAsync(() => {
    const service: Service =  Service.fromJSON({latestSequence: undefined, openRemediations: [], openApprovals: []});

    expect(service).toBeInstanceOf(Service);
  }));
});
