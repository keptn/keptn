import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbCreateSecretFormComponent } from './ktb-create-secret-form.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbCreateSecretFormComponent;
  let fixture: ComponentFixture<KtbCreateSecretFormComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbCreateSecretFormComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbCreateSecretFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
