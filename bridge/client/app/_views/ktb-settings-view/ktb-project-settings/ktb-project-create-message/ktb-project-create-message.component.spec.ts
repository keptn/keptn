import { ComponentFixture, TestBed } from '@angular/core/testing';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbProjectSettingsModule } from '../ktb-project-settings.module';
import { KtbProjectCreateMessageComponent } from './ktb-project-create-message.component';

describe('KtbProjectCreateMessageComponent', () => {
  let component: KtbProjectCreateMessageComponent;
  let fixture: ComponentFixture<KtbProjectCreateMessageComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectSettingsModule, RouterTestingModule],
      declarations: [],
    }).compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbProjectCreateMessageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
