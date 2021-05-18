import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbUniformViewComponent } from './ktb-uniform-view.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbUniformViewComponent;
  let fixture: ComponentFixture<KtbUniformViewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbUniformViewComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbUniformViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
