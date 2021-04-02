import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbHttpLoadingBarComponent } from './ktb-http-loading-bar.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbHorizontalSeparatorComponent} from "../ktb-horizontal-separator/ktb-horizontal-separator.component";

describe('HttpLoadingBarComponent', () => {
  let component: KtbHttpLoadingBarComponent;
  let fixture: ComponentFixture<KtbHttpLoadingBarComponent>;

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
        fixture = TestBed.createComponent(KtbHttpLoadingBarComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
