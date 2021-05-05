import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from "rxjs";

import {Notification} from "../_models/notification";

@Injectable({
  providedIn: 'root'
})
export class NotificationsService {

  private _notifications = new BehaviorSubject<Notification[]>([]);

  constructor() { }

  get notifications(): Observable<Notification[]> {
    return this._notifications.asObservable();
  }

  addNotification(type, message) {
    let notification = Notification.fromJSON({ type, message });
    const duplicateNotifications = this._notifications.getValue().filter(n => n.type === notification.type && n.message === notification.message);
    if(duplicateNotifications.length === 0) {
      this._notifications.next([...this._notifications.getValue(), notification]);
    }
  }

  removeNotification(notification) {
    this._notifications.next(this._notifications.getValue().filter(n => n != notification));
  }
}
