import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';
import { EvaluationBoardComponent } from './evaluation-board.component';
import { AppModule } from '../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('ProjectBoardComponent', () => {
  let component: EvaluationBoardComponent;
  let fixture: ComponentFixture<EvaluationBoardComponent>;

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
        fixture = TestBed.createComponent(EvaluationBoardComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));

});
