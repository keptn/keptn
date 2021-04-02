import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline.component';
import {RouterTestingModule} from '@angular/router/testing';
import {AppModule} from "../../app.module";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbSequenceTasksListComponent} from "../ktb-sequence-tasks-list/ktb-sequence-tasks-list.component";

describe('KtbSequenceTimelineComponent', () => {
  let component: KtbSequenceTimelineComponent;
  let fixture: ComponentFixture<KtbSequenceTimelineComponent>;

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
        fixture = TestBed.createComponent(KtbSequenceTimelineComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
