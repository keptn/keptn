import {Injectable} from '@angular/core';
import {BehaviorSubject, Observable} from 'rxjs';
import { Notification, NotificationType } from '../_models/notification';

@Injectable({
  providedIn: 'root'
})
export class NotificationsService {

  private _notifications = new BehaviorSubject<Notification[]>([]);

  get notifications(): Observable<Notification[]> {
    return this._notifications.asObservable();
  }

  // tslint:disable-next-line:no-any
  addNotification(type: NotificationType, message: string, time?: number, isTemplateRendered = false, data?: any) {
    const notification = new Notification(type, message);
    notification.isTemplateRendered = isTemplateRendered;
    notification.data = data || null;

    if (time) {
      setTimeout(() => {
        this.removeNotification(notification);
      }, time);
    }

    // Check if the notification to add already exists
    const duplicateNotifications = this._notifications.getValue().filter(n => n.type === notification.type && n.message === notification.message);

    // Only show notification if it is not shown yet to prevent duplicates (issue #3896 - https://github.com/keptn/keptn/issues/3896)
    if (duplicateNotifications.length === 0) {
      this._notifications.next([...this._notifications.getValue(), notification]);
    }
  }

  removeNotification(notification: Notification) {
    this._notifications.next(this._notifications.getValue().filter(n => n !== notification));
  }
}
