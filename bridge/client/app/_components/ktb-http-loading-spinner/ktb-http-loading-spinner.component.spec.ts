import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbHttpLoadingSpinnerComponent } from './ktb-http-loading-spinner.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('HttpLoadingSpinnerComponent', () => {
  let component: KtbHttpLoadingSpinnerComponent;
  let fixture: ComponentFixture<KtbHttpLoadingSpinnerComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbHttpLoadingSpinnerComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
