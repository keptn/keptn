import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbNotificationBarComponent} from './ktb-notification-bar.component';
import {Component} from "@angular/core";
import {By} from "@angular/platform-browser";
import {NotificationsService} from "../../_services/notifications.service";
import {NotificationType} from "../../_models/notification";
import {DtIconModule} from "@dynatrace/barista-components/icon";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbMarkdownComponent} from "../ktb-markdown/ktb-markdown.component";

describe('KtbNotificationBarComponent', () => {
  let service: NotificationsService;
  let component: SimpleKtbNotificationBarComponent;
  let fixture: ComponentFixture<SimpleKtbNotificationBarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        SimpleKtbNotificationBarComponent,
        KtbNotificationBarComponent,
      ],
      imports: [
        HttpClientTestingModule,
        DtIconModule,
        DtIconModule.forRoot({
          svgIconLocation: `/assets/icons/{{name}}.svg`,
        }),
      ],
      providers: [NotificationsService]
    })
    .compileComponents()
    .then(() => {
      fixture = TestBed.createComponent(SimpleKtbNotificationBarComponent);
      component = fixture.componentInstance;
      service = TestBed.get(NotificationsService);
      fixture.detectChanges();
    });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));

  it('should add and remove notifications', () => {
    let notifications = fixture.debugElement.queryAll(By.css('.page-note'));
    expect(notifications.length).toBe(0);

    service.addNotification(NotificationType.Info, "Information");
    service.addNotification(NotificationType.Success, "Success");
    service.addNotification(NotificationType.Warning, "Warning");
    service.addNotification(NotificationType.Error, "Error");
    fixture.detectChanges();

    notifications = fixture.debugElement.queryAll(By.css('.page-note'));
    expect(notifications.length).toBe(4);

    expect(notifications[0].nativeElement.classList).toContain('info-note');
    expect(notifications[0].nativeElement.textContent).toContain('Information');
    expect(notifications[1].nativeElement.classList).toContain('success-note');
    expect(notifications[2].nativeElement.classList).toContain('warning-note');
    expect(notifications[3].nativeElement.classList).toContain('error-note');
  });
});

/** Simple component for testing the KtbNotificationBarComponent */
@Component({
  template: `<ktb-notification-bar></ktb-notification-bar>`,
})
class SimpleKtbNotificationBarComponent {}
