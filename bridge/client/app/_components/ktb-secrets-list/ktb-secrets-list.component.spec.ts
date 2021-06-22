import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbSecretsListComponent } from './ktb-secrets-list.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbUniformViewComponent', () => {
  let component: KtbSecretsListComponent;
  let fixture: ComponentFixture<KtbSecretsListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbSecretsListComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbSecretsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
