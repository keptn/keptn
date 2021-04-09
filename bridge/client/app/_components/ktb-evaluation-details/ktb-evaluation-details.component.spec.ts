import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbEvaluationDetailsComponent } from './ktb-evaluation-details.component';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbCopyToClipboardComponent} from "../ktb-copy-to-clipboard/ktb-copy-to-clipboard.component";

describe('KtbEvaluationDetailsComponent', () => {
  let component: KtbEvaluationDetailsComponent;
  let fixture: ComponentFixture<KtbEvaluationDetailsComponent>;

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
        fixture = TestBed.createComponent(KtbEvaluationDetailsComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
