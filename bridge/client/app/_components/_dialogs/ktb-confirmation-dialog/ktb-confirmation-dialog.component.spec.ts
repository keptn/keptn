import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbConfirmationDialogComponent } from './ktb-confirmation-dialog.component';
import { AppModule } from '../../../app.module';


describe('KtbDeletionDialogComponent', () => {
  let component: KtbConfirmationDialogComponent;
  let fixture: ComponentFixture<KtbConfirmationDialogComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbConfirmationDialogComponent],
      imports: [AppModule],
    })
      .compileComponents();
  });

  beforeEach(() => {
    fixture = TestBed.createComponent(KtbConfirmationDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
