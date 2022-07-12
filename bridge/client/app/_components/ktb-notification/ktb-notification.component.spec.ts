/* eslint-disable @typescript-eslint/dot-notation */
import { ComponentFactory, ComponentFactoryResolver, ComponentRef, ViewContainerRef } from '@angular/core';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Notification, NotificationType } from '../../_models/notification';
import { KtbProjectCreateMessageComponent } from '../../_views/ktb-settings-view/ktb-project-settings/ktb-project-create-message/ktb-project-create-message.component';
import { KtbNotificationComponent } from './ktb-notification.component';
import { KtbNotificationModule } from './ktb-notification.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

// eslint-disable-next-line @typescript-eslint/no-explicit-any
declare let global: any;

describe('KtbNotificationComponent', () => {
  let component: KtbNotificationComponent;
  let fixture: ComponentFixture<KtbNotificationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbNotificationModule, HttpClientTestingModule, BrowserAnimationsModule],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbNotificationComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    fixture.detectChanges();
    expect(component).toBeTruthy();
  });

  it('should start timeout', () => {
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    expect(component['fadeOutDelay']).not.toBeUndefined();
    expect(component['hideTimeout']).not.toBeUndefined();
  });

  it('should not start timeout', () => {
    setInputParameter(new Notification(NotificationType.ERROR, 'my message', undefined, -1));
    expect(component['fadeOutDelay']).toBeUndefined();
    expect(component['hideTimeout']).toBeUndefined();
  });

  it('should destroy listeners on hide', () => {
    jest.spyOn(global, 'clearTimeout');
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    component.ngOnDestroy();
    expect(clearTimeout).toHaveBeenCalledTimes(2);
  });

  it('should emit event if closed', () => {
    const spy = jest.spyOn(component.hide, 'next');
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    component.closeComponent();
    expect(spy).toHaveBeenCalled();
  });

  it('should set fadeOut status to out after timeout', () => {
    jest.useFakeTimers();
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    expect(component.fadeStatus).toBe('in');
    jest.runAllTimers();
    expect(component.fadeStatus).toBe('out');
  });

  it('should close component after timeout', () => {
    jest.useFakeTimers();
    const spy = jest.spyOn(component, 'closeComponent');
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    expect(spy).not.toHaveBeenCalled();
    jest.runAllTimers();
    expect(spy).toHaveBeenCalled();
  });

  it('should revert timeout and fadOut status on hover', () => {
    jest.useFakeTimers();
    const clearTimeoutSpy = jest.spyOn(global, 'clearTimeout');
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    jest.runAllTimers();
    expect(component.fadeStatus).toBe('out');
    component.onHover();
    expect(clearTimeoutSpy).toHaveBeenCalledTimes(2); //end and hover
    expect(component.fadeStatus).toBe('in');
  });

  it('should start timer again after leave', () => {
    const spy = jest.spyOn(global, 'setTimeout');
    setInputParameter(new Notification(NotificationType.ERROR, 'my message'));
    component.onHover();
    component.onLeave();
    expect(spy).toHaveBeenCalledTimes(4); // start and restart
  });

  it('should add component and its data', () => {
    const instance = {};
    component.componentViewRef = { createComponent: () => {}, clear: () => {} } as unknown as ViewContainerRef;
    jest
      .spyOn(TestBed.inject(ComponentFactoryResolver), 'resolveComponentFactory')
      .mockReturnValue({} as unknown as ComponentFactory<unknown>);
    jest
      .spyOn(component.componentViewRef as ViewContainerRef, 'createComponent')
      .mockReturnValue({ instance } as unknown as ComponentRef<unknown>);

    setInputParameter(
      new Notification(NotificationType.ERROR, 'my message', {
        component: KtbProjectCreateMessageComponent,
        data: {
          projectName: 'sockshop',
          routerLink: ['/'],
        },
      })
    );

    expect(instance).toMatchObject({
      projectName: 'sockshop',
      routerLink: ['/'],
    });
  });

  function setInputParameter(notification: Notification): void {
    component.notification = notification;
    // weird bug: fixture.detectChanges() sets/reverts the @Input() to empty string. That's why we are calling the life-cycle manually
    component.ngAfterViewInit();
  }
});
/* eslint-enable @typescript-eslint/dot-notation */
