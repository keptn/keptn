import { ComponentFixture, TestBed, waitForAsync } from '@angular/core/testing';

import { KtbStageOverviewComponent } from './ktb-stage-overview.component';
import {AppModule} from '../../app.module';
import {HttpClientTestingModule} from '@angular/common/http/testing';

describe('KtbStageOverviewComponent', () => {
  let component: KtbStageOverviewComponent;
  let fixture: ComponentFixture<KtbStageOverviewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [ KtbStageOverviewComponent ],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbStageOverviewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });
});
