import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEnvironmentViewComponent } from './ktb-environment-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbEnvironmentViewModule } from './ktb-environment-view.module';
import { RouterTestingModule } from '@angular/router/testing';
import { POLLING_INTERVAL_MILLIS } from '../../_utils/app.utils';

describe('KtbEnvironmentViewComponent', () => {
  let component: KtbEnvironmentViewComponent;
  let fixture: ComponentFixture<KtbEnvironmentViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEnvironmentViewModule, RouterTestingModule, HttpClientTestingModule],
      providers: [{ provide: POLLING_INTERVAL_MILLIS, useValue: 0 }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbEnvironmentViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
