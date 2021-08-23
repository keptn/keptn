import { ComponentFixture, fakeAsync, flush, TestBed, tick } from '@angular/core/testing';
import { KtbDeleteConfirmationComponent } from './ktb-delete-confirmation.component';
import { AppModule } from '../../../app.module';
import { TestUtils } from '../../../_utils/test.utils';

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

  it('should emit', fakeAsync(() => {
    // given
    const type = 'subscription';
    component.name = 'My item';
    component.type = type;
    component.dialogState = 'confirm';
    TestUtils.updateDialog(fixture);

    // when
    const spyEmit = jest.spyOn(component.confirmClicked, 'emit');
    const deleteButton: HTMLElement | null = document.querySelector('dt-confirmation-dialog-state[name=confirm] button[uitestid=dialogDeleteButton]');
    deleteButton?.click();
    TestUtils.updateDialog(fixture);

    // then
    expect(component.dialogState).toEqual('deleting');
    expect(document.querySelector('dt-confirmation-dialog-state[name=deleting]')?.textContent?.trim()).toEqual(`Deleting ${type} ...`);
    expect(spyEmit).toHaveBeenCalled();
    flush();
  }));

  it('should dismiss dialog after success and 2 seconds', fakeAsync(() => {
    // given
    component.name = 'My item';
    component.type = 'subscription';

    // when
    component.dialogState = 'success';
    TestUtils.updateDialog(fixture);
    expect(document.querySelector('dt-confirmation-dialog-state[name=success]')?.textContent?.trim()).toEqual('Subscription deleted successfully!');
    tick(2010);
    TestUtils.updateDialog(fixture);

    // then
    expect(component.dialogState).toBeNull();
    isDialogClosed();
    flush();
  }));

  it('should reset timeout on confirm', fakeAsync(() => {
    // given
    component.name = 'My item';
    component.type = 'subscription';
    component.dialogState = 'success';
    TestUtils.updateDialog(fixture);

    // when
    tick(1000);
    component.dialogState = 'confirm';

    // then
    tick(2000);
    expect(component.dialogState).toEqual('confirm');
  }));

  it('should dismiss dialog on cancel', fakeAsync(() => {
    // given
    component.name = 'My item';
    component.type = 'subscription';
    component.dialogState = 'confirm';
    TestUtils.updateDialog(fixture);
    // when
    const cancelButton: HTMLElement | null = document.querySelector('dt-confirmation-dialog-state[name=confirm] button[uitestid=dialogCancelButton]');
    cancelButton?.click();

    // then
    expect(component.dialogState).toBeNull();
    TestUtils.updateDialog(fixture);
    isDialogClosed();
    flush();
  }));

  function isDialogClosed() {
    for (const state of Array.from(document.querySelectorAll('dt-confirmation-dialog-state'))) {
      expect(state.textContent).toEqual('');
    }
  }
});
