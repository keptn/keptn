import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbIntegrationViewComponent', () => {
  let component: KtbIntegrationViewComponent;
  let fixture: ComponentFixture<KtbIntegrationViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbIntegrationViewComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbIntegrationViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
