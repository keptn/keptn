import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: POLLING_INTERVAL_MILLIS, useValue: 0 }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
