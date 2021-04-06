import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbServicesListComponent } from './ktb-services-list.component';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbSequenceTimelineComponent} from "../ktb-sequence-timeline/ktb-sequence-timeline.component";

describe('KtbServicesListComponent', () => {
  let component: KtbServicesListComponent;
  let fixture: ComponentFixture<KtbServicesListComponent>;

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
        fixture = TestBed.createComponent(KtbServicesListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
