import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbHttpLoadingBarComponent} from "../ktb-http-loading-bar/ktb-http-loading-bar.component";

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

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
        fixture = TestBed.createComponent(KtbKeptnServicesListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
