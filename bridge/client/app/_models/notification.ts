enum NotificationType {
  Info= 'info',
  Success = 'success',
  Warning = 'warning',
  Error = 'error'
}

class Notification {
  type: NotificationType;
  message: string;

  static fromJSON(data: any) {
    return Object.assign(new this, data);
  }

}

export {Notification, NotificationType}
