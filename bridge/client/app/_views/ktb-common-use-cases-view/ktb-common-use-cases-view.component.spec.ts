import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbCommonUseCasesViewComponent } from './ktb-common-use-cases-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { KtbCommonUseCasesViewModule } from './ktb-common-use-cases-view.module';

describe('KtbIntegrationViewComponent', () => {
  let component: KtbCommonUseCasesViewComponent;
  let fixture: ComponentFixture<KtbCommonUseCasesViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbCommonUseCasesViewModule, HttpClientTestingModule],
      providers: [{ provide: POLLING_INTERVAL_MILLIS, useValue: 0 }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbCommonUseCasesViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
