import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsComponent } from './ktb-service-settings.component';
import { AppModule } from '../../app.module';

describe('KtbServiceSettingsComponent', () => {
  let component: KtbServiceSettingsComponent;
  let fixture: ComponentFixture<KtbServiceSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbServiceSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
