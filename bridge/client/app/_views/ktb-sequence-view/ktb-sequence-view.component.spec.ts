import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSequenceViewComponent } from './ktb-sequence-view.component';
import { AppModule } from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceViewComponent;
  let fixture: ComponentFixture<KtbSequenceViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSequenceViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
