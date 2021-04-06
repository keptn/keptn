import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import { KtbRootEventsListComponent } from './ktb-root-events-list.component';
import {KtbEventsListComponent} from "../ktb-events-list/ktb-events-list.component";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppModule} from "../../app.module";
import {KtbProjectTileComponent} from "../ktb-project-tile/ktb-project-tile.component";

describe('KtbEventsListComponent', () => {
  let component: KtbRootEventsListComponent;
  let fixture: ComponentFixture<KtbRootEventsListComponent>;

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
        fixture = TestBed.createComponent(KtbRootEventsListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
