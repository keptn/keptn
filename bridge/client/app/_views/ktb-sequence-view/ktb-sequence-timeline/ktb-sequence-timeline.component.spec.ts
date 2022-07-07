import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSequenceViewModule } from '../ktb-sequence-view.module';

describe('KtbSequenceTimelineComponent', () => {
  let component: KtbSequenceTimelineComponent;
  let fixture: ComponentFixture<KtbSequenceTimelineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbSequenceViewModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
