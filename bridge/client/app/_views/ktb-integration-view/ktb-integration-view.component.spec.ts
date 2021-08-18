import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';
import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { INITIAL_DELAY_MILLIS } from '../../_utils/app.utils';
import { AppModule } from '../../app.module';

describe('KtbIntegrationViewComponent', () => {
  let component: KtbIntegrationViewComponent;
  let fixture: ComponentFixture<KtbIntegrationViewComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: INITIAL_DELAY_MILLIS, useValue: 0},
      ],
    })
      .compileComponents()
      .then(() => {
        fixture = TestBed.createComponent(KtbIntegrationViewComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
