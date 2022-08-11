import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { DeleteData, DeletionProgressEvent } from '../_interfaces/delete';

@Injectable({
  providedIn: 'root',
})
export class EventService {
  public deletionTriggeredEvent = new Subject<DeleteData>();
  public deletionProgressEvent = new Subject<DeletionProgressEvent>();
}
