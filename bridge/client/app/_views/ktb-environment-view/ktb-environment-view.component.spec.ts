import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbEnvironmentViewComponent', () => {
  let component: KtbEnvironmentViewComponent;
  let fixture: ComponentFixture<KtbEnvironmentViewComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbEnvironmentViewComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbEnvironmentViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
