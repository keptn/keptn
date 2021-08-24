import { TestBed } from '@angular/core/testing';

import { EventService } from './event.service';
import { DeleteResult, DeleteType } from '../_interfaces/delete';

describe('EventService', () => {
  let service: EventService;

  beforeEach(() => {
    TestBed.configureTestingModule({});
    service = TestBed.inject(EventService);
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should have a Subject deletionTriggeredEvent open', () => {
    expect(service.deletionTriggeredEvent.closed).toEqual(false);
  });

  it('should have a Subject deletionProgressEvent open', () => {
    expect(service.deletionProgressEvent.closed).toEqual(false);
  });

  it('should deletionTriggeredEvent have been called with next value', () => {
    // given
    const data = {name: 'sockshop', type: DeleteType.PROJECT};
    const spy = jest.spyOn(service.deletionTriggeredEvent, 'next');

    // when
    service.deletionTriggeredEvent.next(data);

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith(data);
  });

  it('should deletionProgressEvent have been called with next value', () => {
    // given
    const data = {isInProgress: false, result: DeleteResult.ERROR, error: 'Error'};
    const spy = jest.spyOn(service.deletionProgressEvent, 'next');

    // when
    service.deletionProgressEvent.next(data);

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith(data);
  });
});
