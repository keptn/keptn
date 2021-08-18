import { Notification, NotificationType } from './notification';

describe('Notification', () => {
  it('should create a new instance', () => {
    // given
    const notification = new Notification(NotificationType.Info, 'test');

    // then
    expect(notification).toBeTruthy();
  });
});
