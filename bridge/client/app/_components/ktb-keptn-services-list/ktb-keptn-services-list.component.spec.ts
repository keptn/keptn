import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
        KtbKeptnServicesListComponent
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbKeptnServicesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
