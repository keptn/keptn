import {async, ComponentFixture, fakeAsync, TestBed} from '@angular/core/testing';

import {ProjectBoardComponent} from './project-board.component';
import {AppModule} from '../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {DashboardComponent} from "../dashboard/dashboard.component";

describe('ProjectBoardComponent', () => {
  let component: ProjectBoardComponent;
  let fixture: ComponentFixture<ProjectBoardComponent>;

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
        fixture = TestBed.createComponent(ProjectBoardComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
