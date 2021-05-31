import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbServiceRemediationListComponent } from './ktb-service-remediation-list.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbServiceRemediationListComponent', () => {
  let component: KtbServiceRemediationListComponent;
  let fixture: ComponentFixture<KtbServiceRemediationListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [ KtbServiceRemediationListComponent ],
      imports: [AppModule, HttpClientTestingModule]
    })
    .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbServiceRemediationListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
