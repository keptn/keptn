import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbUniformViewComponent } from './ktb-uniform-view.component';
import {AppModule} from '../../app.module';
import {HttpClient} from '@angular/common/http';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbUniformViewComponent;
  let fixture: ComponentFixture<KtbUniformViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbUniformViewComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbUniformViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
