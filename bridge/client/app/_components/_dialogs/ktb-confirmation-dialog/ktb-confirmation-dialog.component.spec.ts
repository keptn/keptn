import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbConfirmationDialogComponent } from './ktb-confirmation-dialog.component';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbConfirmationDialogModule } from './ktb-confirmation-dialog.module';

describe('KtbConfirmationDialogComponent', () => {
  let component: KtbConfirmationDialogComponent;
  let fixture: ComponentFixture<KtbConfirmationDialogComponent>;
  const dialogData = {
    sequence: {
      shkeptncontext: 'f6a38eb6-e99d-4d14-ab4c-3e94ed288b45',
      name: 'delivery',
    },
    confirmCallback: (): void => {},
  };
  const dialogRefMock = {
    close: (): void => {},
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbConfirmationDialogModule, MatDialogModule, HttpClientTestingModule],
      providers: [
        { provide: MAT_DIALOG_DATA, useValue: dialogData },
        { provide: MatDialogRef, useValue: dialogRefMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbConfirmationDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should call callback and close the dialog when user clicked on confirm', () => {
    // given
    const confirmButton = fixture.nativeElement.querySelector('.danger-button');
    const spyClose = jest.spyOn(component.dialogRef, 'close');
    const spyConfirmCallback = jest.spyOn(component.data, 'confirmCallback');

    // when
    confirmButton.dispatchEvent(new Event('click'));

    // then
    expect(spyClose).toHaveBeenCalled();
    expect(spyConfirmCallback).toHaveBeenCalled();
  });
});
