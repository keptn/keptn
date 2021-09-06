import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbWebhookSettingsComponent', () => {
  let component: KtbWebhookSettingsComponent;
  let fixture: ComponentFixture<KtbWebhookSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbWebhookSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
