import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {KtbHorizontalSeparatorComponent, KtbHorizontalSeparatorTitle} from './ktb-horizontal-separator.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";

describe('KtbHorizontalSeparatorComponent', () => {
  let component: KtbHorizontalSeparatorComponent;
  let fixture: ComponentFixture<KtbHorizontalSeparatorComponent>;

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
    fixture = TestBed.createComponent(KtbHorizontalSeparatorComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
