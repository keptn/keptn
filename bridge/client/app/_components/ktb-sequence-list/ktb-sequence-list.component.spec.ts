import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbSequenceListComponent } from './ktb-sequence-list.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbSequenceListComponent', () => {
  let component: KtbSequenceListComponent;
  let fixture: ComponentFixture<KtbSequenceListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbSequenceListComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSequenceListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
