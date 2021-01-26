import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSequenceOverviewComponent } from './ktb-sequence-overview.component';
import { AppModule } from '../../app.module';
import {HttpClientTestingModule} from "@angular/common/http/testing";

describe('KtbEventsListComponent', () => {
  let component: KtbSequenceOverviewComponent;
  let fixture: ComponentFixture<KtbSequenceOverviewComponent>;

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
    fixture = TestBed.createComponent(KtbSequenceOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
