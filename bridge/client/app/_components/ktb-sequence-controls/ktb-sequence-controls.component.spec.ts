import {ComponentFixture, TestBed} from '@angular/core/testing';

import {KtbSequenceControlsComponent} from './ktb-sequence-controls.component';
import {AppModule} from '../../app.module';

describe('KtbSequenceStateListComponent', () => {
  let component: KtbSequenceControlsComponent;
  let fixture: ComponentFixture<KtbSequenceControlsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbSequenceControlsComponent],
      imports: [
        AppModule,
      ],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSequenceControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
