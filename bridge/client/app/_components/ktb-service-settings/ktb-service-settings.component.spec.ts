import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsComponent } from './ktb-service-settings.component';
import { KtbServiceSettingsModule } from './ktb-service-settings.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('KtbServiceSettingsComponent', () => {
  let component: KtbServiceSettingsComponent;
  let fixture: ComponentFixture<KtbServiceSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceSettingsModule, RouterTestingModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbServiceSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
