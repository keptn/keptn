import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';
import { AppModule } from '../../app.module';

describe('KtbIntegrationViewComponent', () => {
  let component: KtbIntegrationViewComponent;
  let fixture: ComponentFixture<KtbIntegrationViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: POLLING_INTERVAL_MILLIS, useValue: 0},
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbIntegrationViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
