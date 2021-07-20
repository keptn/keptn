enum NotificationType {
  Info= 'info',
  Success = 'success',
  Warning = 'warning',
  Error = 'error'
}

class Notification {
  constructor(public type: NotificationType, public message: string ) {
  }
}

export {Notification, NotificationType};
