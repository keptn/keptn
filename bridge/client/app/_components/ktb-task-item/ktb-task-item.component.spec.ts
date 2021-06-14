import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import {KtbTaskItemComponent} from './ktb-task-item.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventItemComponent', () => {
  let component: KtbTaskItemComponent;
  let fixture: ComponentFixture<KtbTaskItemComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbTaskItemComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
