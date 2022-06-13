import { KtbDeletionDialogComponent } from './ktb-deletion-dialog.component';
import { DeleteResult, DeleteType } from '../../../_interfaces/delete';
import { EventService } from '../../../_services/event.service';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { KtbDeletionDialogModule } from './ktb-deletion-dialog.module';

describe('KtbDeletionDialogComponent', () => {
  let component: KtbDeletionDialogComponent;
  let fixture: ComponentFixture<KtbDeletionDialogComponent>;
  const dialogData = { name: 'sockshop', type: DeleteType.PROJECT };
  let eventService: EventService;
  const dialogRefMock = {
    close: (): void => {},
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbDeletionDialogModule, MatDialogModule, HttpClientTestingModule],
      providers: [
        EventService,
        { provide: MAT_DIALOG_DATA, useValue: dialogData },
        { provide: MatDialogRef, useValue: dialogRefMock },
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbDeletionDialogComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should be a valid form when input matches the given name', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('sockshop');
    component.deletionConfirmationForm.updateValueAndValidity();

    // then
    expect(component.deletionConfirmationForm.valid).toBe(true);
  });

  it('should be an invalid form when input is empty', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('');
    component.deletionConfirmationForm.updateValueAndValidity();

    // then
    expect(component.deletionConfirmationForm.invalid).toBe(true);
  });

  it('should be an invalid form when name and input do not match', () => {
    // given
    const input = component.deletionConfirmationControl;
    const values = ['sock', '$ock', '1', 'Sockshop', 'sockshoP', 'sOckshop', ''];

    values.forEach((val) => {
      input.setValue(val);
      component.deletionConfirmationForm.updateValueAndValidity();
      expect(component.deletionConfirmationForm.invalid).toBe(true);
    });
  });

  it('should have a disabled button when form is invalid', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('sock');
    component.deletionConfirmationForm.updateValueAndValidity();
    fixture.detectChanges();

    // then
    const button = fixture.nativeElement.querySelector('.danger-button');
    expect(button.disabled).toBe(true);
  });

  it('should have a disabled button when input was first valid and then is invalid', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('sockshop');
    component.deletionConfirmationForm.updateValueAndValidity();
    fixture.detectChanges();

    // then
    let button = fixture.nativeElement.querySelector('.danger-button');
    expect(button.disabled).toBe(false);

    // when
    input.setValue('sock');
    component.deletionConfirmationForm.updateValueAndValidity();
    fixture.detectChanges();

    // then
    button = fixture.nativeElement.querySelector('.danger-button');
    expect(button.disabled).toBe(true);
  });

  it('should have an enabled button when form is valid', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('sockshop');
    component.deletionConfirmationForm.updateValueAndValidity();
    fixture.detectChanges();

    // then
    const button = fixture.nativeElement.querySelector('.danger-button');
    expect(button.disabled).toBe(false);
  });

  it('should trigger an deletion event when deletion button is clicked', () => {
    // given
    const button = fixture.nativeElement.querySelector('.danger-button');
    const spy = jest.spyOn(component, 'deleteConfirm');

    // when
    button.dispatchEvent(new Event('click'));
    fixture.detectChanges();

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should call EventService.deletionTriggeredEvent when deleteProject is called', () => {
    // given
    eventService = TestBed.inject(EventService);
    const spy = jest.spyOn(eventService.deletionTriggeredEvent, 'next');

    // when
    component.deleteConfirm();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith({ ...dialogData });
  });

  it('should close the dialog when deletionProgressEvent result is SUCCESS', () => {
    // given
    eventService = TestBed.inject(EventService);
    const spy = jest.spyOn(component.dialogRef, 'close');

    // when
    eventService.deletionProgressEvent.next({ isInProgress: false, result: DeleteResult.SUCCESS });

    // then
    expect(spy).toHaveBeenCalled();
  });
});
