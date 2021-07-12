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
  type: NotificationType;
  message: string;
  isTemplateRendered?: boolean;
  data?: any;

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }

}

export {Notification, NotificationType, TemplateRenderedNotifications}
