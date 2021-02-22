import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSequenceTimelineComponent } from './ktb-sequence-timeline.component';
import {RouterTestingModule} from '@angular/router/testing';

describe('KtbSequenceTimelineComponent', () => {
  let component: KtbSequenceTimelineComponent;
  let fixture: ComponentFixture<KtbSequenceTimelineComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSequenceTimelineComponent ],
      imports: [
        RouterTestingModule
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSequenceTimelineComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
