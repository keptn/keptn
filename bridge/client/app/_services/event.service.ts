import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { DeletionProgressEvent, DeletionTriggeredEvent } from '../_interfaces/delete';

@Injectable({
  providedIn: 'root',
})
export class EventService {
  public deletionTriggeredEvent = new Subject<DeletionTriggeredEvent>();
  public deletionProgressEvent = new Subject<DeletionProgressEvent>();
}
