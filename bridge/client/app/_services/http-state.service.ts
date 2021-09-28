import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

import { HttpState } from '../_models/http-progress-state';

@Injectable({
  providedIn: 'root',
})
export class HttpStateService {
  public state = new BehaviorSubject<HttpState>({} as HttpState);
}
