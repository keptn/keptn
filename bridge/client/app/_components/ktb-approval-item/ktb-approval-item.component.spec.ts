import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import {KtbApprovalItemComponent} from './ktb-approval-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventItemComponent', () => {
  let component: KtbApprovalItemComponent;
  let fixture: ComponentFixture<KtbApprovalItemComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents().then(() => {
      fixture = TestBed.createComponent(KtbApprovalItemComponent);
      component = fixture.componentInstance;
      fixture.detectChanges();
    });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
