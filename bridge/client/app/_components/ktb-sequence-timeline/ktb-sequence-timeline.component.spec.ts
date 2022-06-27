import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbSequenceTimelineModule } from './ktb-sequence-timeline.module';

describe('KtbSequenceTimelineComponent', () => {
  let component: KtbSequenceTimelineComponent;
  let fixture: ComponentFixture<KtbSequenceTimelineComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbSequenceTimelineModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
