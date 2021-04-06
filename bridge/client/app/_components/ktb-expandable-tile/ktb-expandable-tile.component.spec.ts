import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {KtbExpandableTileComponent, KtbExpandableTileHeader} from './ktb-expandable-tile.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbEventsListComponent} from "../ktb-events-list/ktb-events-list.component";

describe('KtbExpandableTileComponent', () => {
  let component: KtbExpandableTileComponent;
  let fixture: ComponentFixture<KtbExpandableTileComponent>;

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
        fixture = TestBed.createComponent(KtbExpandableTileComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
