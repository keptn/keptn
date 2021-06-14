import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

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
