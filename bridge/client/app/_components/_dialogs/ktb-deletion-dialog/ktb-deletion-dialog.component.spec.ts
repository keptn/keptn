import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbDeletionDialogComponent } from './ktb-deletion-dialog.component';
import { MAT_DIALOG_DATA, MatDialogModule, MatDialogRef } from '@angular/material/dialog';
import { AppModule } from '../../../app.module';
import { DeleteResult, DeleteType } from '../../../_interfaces/delete';
import { EventService } from '../../../_services/event.service';


describe('KtbDeletionDialogComponent', () => {
  let component: KtbDeletionDialogComponent;
  let fixture: ComponentFixture<KtbDeletionDialogComponent>;
  const dialogData = {name: 'sockshop', type: DeleteType.PROJECT};
  let eventService: EventService;
  const dialogRefMock = {
    close: () => { }
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [KtbDeletionDialogComponent],
      imports: [AppModule, MatDialogModule],
      providers: [
        EventService,
        {provide: MAT_DIALOG_DATA, useValue: dialogData},
        {provide: MatDialogRef, useValue: dialogRefMock}
      ]
    })
      .compileComponents();
  });

  beforeEach(() => {
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
    expect(component.deletionConfirmationForm.valid).toBeTrue();
  });

  it('should be an invalid form when input is empty', () => {
    // given
    const input = component.deletionConfirmationControl;

    // when
    input.setValue('');
    component.deletionConfirmationForm.updateValueAndValidity();

    // then
    expect(component.deletionConfirmationForm.invalid).toBeTrue();
  });

  it('should be an invalid form when name and input do not match', () => {
    // given
    const input = component.deletionConfirmationControl;

    input.setValue('sock');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('$ock');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('1');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('Sockshop');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('sockshoP');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('sOckshop');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();

    input.setValue('');
    component.deletionConfirmationForm.updateValueAndValidity();
    expect(component.deletionConfirmationForm.invalid).toBeTrue();
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
    expect(button.disabled).toBeTrue();
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
    expect(button.disabled).toBeFalse();
  });

  it('should trigger an deletion event when deletion button is clicked', () => {
    // given
    const button = fixture.nativeElement.querySelector('.danger-button');
    const spy = spyOn(component, 'deleteProject');

    // when
    button.dispatchEvent(new Event('click'));
    fixture.detectChanges();

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should call EventService.deletionTriggeredEvent when deleteProject is called', () => {
    // given
    eventService = TestBed.inject(EventService);
    const spy = spyOn(eventService.deletionTriggeredEvent, 'next');

    // when
    component.deleteProject();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy.calls.mostRecent().args[0]).toEqual(dialogData);
  });

  it('should close the dialog when deletionProgressEvent result is SUCCESS', () => {
    // given
    eventService = TestBed.inject(EventService);
    const spy = spyOn(component.dialogRef, 'close').and.callThrough();

    // when
    eventService.deletionProgressEvent.next({isInProgress: false, result: DeleteResult.SUCCESS});

    // then
    expect(spy).toHaveBeenCalled();
  });
});
