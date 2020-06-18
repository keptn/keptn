import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbHttpLoadingBarComponent } from './ktb-http-loading-bar.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('HttpLoadingBarComponent', () => {
  let component: KtbHttpLoadingBarComponent;
  let fixture: ComponentFixture<KtbHttpLoadingBarComponent>;

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
    fixture = TestBed.createComponent(KtbHttpLoadingBarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create an instance', () => {
    expect(component).toBeTruthy();
  });
});
