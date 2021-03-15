import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbServicesListComponent } from './ktb-services-list.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('KtbServicesListComponent', () => {
  let component: KtbServicesListComponent;
  let fixture: ComponentFixture<KtbServicesListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbServicesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
