import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeleteConfirmationComponent } from './ktb-delete-confirmation.component';
import { AppModule } from '../../../app.module';

describe('KtbDeleteConfirmationComponent', () => {
  let component: KtbDeleteConfirmationComponent;
  let fixture: ComponentFixture<KtbDeleteConfirmationComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbDeleteConfirmationComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
