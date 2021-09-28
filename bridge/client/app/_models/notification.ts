enum NotificationType {
  INFO = 'info',
  SUCCESS = 'success',
  WARNING = 'warning',
  ERROR = 'error',
}

enum TemplateRenderedNotifications {
  CREATE_PROJECT = 'message_create_project',
}

class Notification {
  isTemplateRendered?: boolean;
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  data?: any;

  constructor(public type: NotificationType, public message: string) {}
}

export { Notification, NotificationType, TemplateRenderedNotifications };
