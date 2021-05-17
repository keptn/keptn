import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSequenceTasksListComponent } from './ktb-sequence-tasks-list.component';
import { AppModule } from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {KtbProjectTileComponent} from "../ktb-project-tile/ktb-project-tile.component";

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceTasksListComponent;
  let fixture: ComponentFixture<KtbSequenceTasksListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbSequenceTasksListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
