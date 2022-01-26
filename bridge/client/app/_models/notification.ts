import { Type } from '@angular/core';

enum NotificationType {
  INFO = 'info',
  SUCCESS = 'success',
  WARNING = 'warning',
  ERROR = 'error',
}

export interface ComponentInfo<T> {
  component: Type<T>;
  data: { [key in keyof T]: unknown };
}

class Notification {
  /**
   * Timeout of the notification in milliseconds.
   * <br>If -1, the timeout is disabled and the notification is visualized as pinned.
   * <br>If not provided, the value is set to 5,000
   * @param time
   */
  public time: number;

  constructor(
    public severity: NotificationType,
    public message: string = '',
    public componentInfo?: ComponentInfo<unknown>,
    time?: number
  ) {
    this.time = time ?? 5_000;
  }

  get isPinned(): boolean {
    return this.time === -1;
  }
}

export { Notification, NotificationType };
