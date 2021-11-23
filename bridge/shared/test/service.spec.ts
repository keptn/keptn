import { ServiceMock } from '../fixtures/service.mock';
import { Service } from '../models/service';

describe('Service', () => {
  let service: Service;
  beforeEach(() => {
    service = Service.fromJSON(ServiceMock);
  });

  it('should get latest service event', () => {
    const serviceEvent = service.getLatestEvent();
    expect(serviceEvent).toEqual({
      eventId: 'b1079092-4c45-4583-abcc-c248528f7dd4',
      keptnContext: '2e21574c-dcf7-4275-b677-6bc19214acd5',
      time: '1636115048858051676',
    });
  });

  it('should not get latest service event', () => {
    service.lastEventTypes = {};
    const serviceEvent = service.getLatestEvent();
    expect(serviceEvent).toBeUndefined();
  });

  it('should get image version', () => {
    const version = service.getImageVersion();
    expect(version).toEqual('0.12.3');
  });

  it('should not get image version', () => {
    service.deployedImage = undefined;
    const version = service.getImageVersion();
    expect(version).toBeUndefined();
  });

  it('should get short image name', () => {
    const imageName = service.getShortImageName();
    expect(imageName).toEqual('carts');
  });

  it('should not get short image name', () => {
    service.deployedImage = undefined;
    const imageName = service.getShortImageName();
    expect(imageName).toBeUndefined();
  });
});
