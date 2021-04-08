import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { AppHeaderComponent } from './app-header.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../app.module";
import {KtbServiceViewComponent} from "../_views/ktb-service-view/ktb-service-view.component";

describe('AppHeaderComponent', () => {
  let component: AppHeaderComponent;
  let fixture: ComponentFixture<AppHeaderComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(AppHeaderComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
