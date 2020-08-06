import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import {EvaluationBoardComponent} from './evaluation-board.component';
import {AppModule} from '../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('ProjectBoardComponent', () => {
  let component: EvaluationBoardComponent;
  let fixture: ComponentFixture<EvaluationBoardComponent>;

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
    fixture = TestBed.createComponent(EvaluationBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
