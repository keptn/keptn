import {ComponentFixture, TestBed} from '@angular/core/testing';

import {KtbSequenceControlsComponent} from './ktb-sequence-controls.component';
import {AppModule} from '../../app.module';

describe('KtbSequenceStateListComponent', () => {
  let component: KtbSequenceControlsComponent;
  let fixture: ComponentFixture<KtbSequenceControlsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbSequenceControlsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
