import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {ProjectBoardComponent} from './project-board.component';
import {AppModule} from '../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

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
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
