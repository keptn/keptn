enum NotificationType {
  Info= 'info',
  Success = 'success',
  Warning = 'warning',
  Error = 'error'
}

enum TemplateRenderedNotifications {
  CREATE_PROJECT = 'message_create_project'
}

class Notification {
  isTemplateRendered?: boolean;
  // tslint:disable-next-line:no-any
  data?: any;

  constructor(public type: NotificationType, public message: string) {
  }
}

export {Notification, NotificationType, TemplateRenderedNotifications};

