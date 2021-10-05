import { Service } from './service';

describe('Service', () => {
  it('should create instances from json', () => {
    const service: Service = Service.fromJSON({ latestSequence: undefined, openRemediations: [], openApprovals: [] });

    expect(service).toBeInstanceOf(Service);
  });
});
