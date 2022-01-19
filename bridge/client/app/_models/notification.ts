import { Type } from '@angular/core';

enum NotificationType {
  INFO = 'info',
  SUCCESS = 'success',
  WARNING = 'warning',
  ERROR = 'error',
}

export interface ComponentInfo {
  component: Type<unknown>;
  data: Record<string, unknown>;
}

class Notification {
  constructor(
    public type: NotificationType,
    public message: string = '',
    public componentInfo?: ComponentInfo,
    public time?: number
  ) {}
}

export { Notification, NotificationType };
