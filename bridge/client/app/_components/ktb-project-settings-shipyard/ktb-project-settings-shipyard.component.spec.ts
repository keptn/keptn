import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsShipyardComponent } from './ktb-project-settings-shipyard.component';
import { AppModule } from '../../app.module';
import { TestUtils } from '../../_utils/test.utils';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbProjectSettingsEditProjectComponent', () => {
  let component: KtbProjectSettingsShipyardComponent;
  let fixture: ComponentFixture<KtbProjectSettingsShipyardComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectSettingsShipyardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should hide the "Update shipyard" button when create mode is enabled', () => {
    // given
    component.isCreateMode = true;
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button.shipyard-upload-button');

    // then
    expect(button).toBeFalsy();
  });

  it('should show the "Update shipyard" button when create mode is disabled', () => {
    // given
    component.isCreateMode = false;
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button.shipyard-upload-button');

    // then
    expect(button).toBeTruthy();
  });

  it('should disable the "Update shipyard" button when no file is selected', () => {
    // given
    component.isCreateMode = false;
    component.shipyardFile = undefined;
    const button = fixture.nativeElement.querySelector('button.shipyard-upload-button');

    // then
    expect(button.disabled).toBe(true);
  });

  it('should enable the "Update shipyard" button when a file is selected', () => {
    // given
    component.isCreateMode = false;
    component.shipyardFile = new File(['test content'], 'test1.yaml');

    // when
    fixture.detectChanges();

    // then
    const button = fixture.nativeElement.querySelector('button.shipyard-upload-button');
    expect(button.disabled).toBe(false);
  });

  it('should show an error when an error handling is triggered', () => {
    // given
    const errorContainer = fixture.nativeElement.querySelector('p.drop-error');

    // when
    component.handleDragAndDropError('Drag and drop error');
    fixture.detectChanges();

    // then
    expect(errorContainer.innerText).toEqual('Drag and drop error');
  });

  it('should show an error when the file selected by the file input is a directory', () => {
    // given
    const errorContainer = fixture.nativeElement.querySelector('p.drop-error');
    const file = new File(['test content'], 'test-directory', { type: '' });
    Object.defineProperty(file.constructor.prototype, 'size', {
      value: 4096,
    });
    const fileList = TestUtils.createNewDropEventWithFiles([file]).dataTransfer?.files ?? null;

    // when
    component.validateAndUpdateFile(fileList);
    fixture.detectChanges();

    // then
    expect(errorContainer.innerText).toEqual(`Please select only files`);
  });

  it('should show an error when the file selected by the file input has not the right file extension', () => {
    // given
    const errorContainer = fixture.nativeElement.querySelector('p.drop-error');
    const fileList =
      TestUtils.createNewDropEventWithFiles([new File(['test content'], 'test1.png', { type: 'image/png' })])
        .dataTransfer?.files ?? null;

    // when
    component.validateAndUpdateFile(fileList);
    fixture.detectChanges();

    // then
    expect(errorContainer.innerText).toEqual(`Only ${component.allowedExtensions.join(', ')} files allowed`);
  });

  it('should update the files when dropped handler is called', () => {
    // given
    const fileList = TestUtils.createNewDropEventWithFiles([
      new File(['test content'], 'test1.png', { type: 'image/png' }),
    ]).dataTransfer?.files;

    // when
    component.updateFile(fileList);
    fixture.detectChanges();

    // then
    expect(component.shipyardFile).toEqual(fileList?.[0]);
  });

  it('should call validateAndUpdateFile when file input was changed', () => {
    // given
    const spy = jest.spyOn(component, 'validateAndUpdateFile');
    const input = fixture.nativeElement.querySelector('#shipyard-file-input');

    // when
    input.dispatchEvent(new Event('change'));

    // then
    expect(spy).toHaveBeenCalled();
  });

  it('should emit the updated files', () => {
    // given
    const fileList = TestUtils.createNewDropEventWithFiles([new File(['test content'], 'test1.yaml')]).dataTransfer
      ?.files;
    const spy = jest.spyOn(component.shipyardFileChanged, 'emit');

    // when
    component.updateFile(fileList);

    // then
    expect(spy).toHaveBeenCalledWith(fileList?.[0]);
  });

  it('should show the file when a shipyard file is set', () => {
    // given
    component.shipyardFile = new File(['test content'], 'test1.yaml');

    // when
    fixture.detectChanges();

    // then
    const div = fixture.nativeElement.querySelector('.shipyard-file-name');
    expect(div.textContent).toEqual('test1.yaml');
  });

  it('should show a delete button when a shipyard file is set', () => {
    // given
    component.shipyardFile = new File(['test content'], 'test1.yaml');

    // when
    fixture.detectChanges();

    // then
    const button = fixture.nativeElement.querySelector('.shipyard-delete-button');
    expect(button).toBeTruthy();
  });

  it('should not show a delete button when no shipyard file is set', () => {
    const button = fixture.nativeElement.querySelector('.shipyard-delete-button');
    expect(button).toBeFalsy();
  });

  it('should remove the shipyard file when delete button is clicked', () => {
    // given
    component.shipyardFile = new File(['test content'], 'test1.yaml');

    // when
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('.shipyard-delete-button');
    button.dispatchEvent(new Event('click'));
    fixture.detectChanges();

    // then
    const div = fixture.nativeElement.querySelector('.shipyard-file-name');
    expect(div).toBeFalsy();
    expect(component.shipyardFile).toBeUndefined();
  });

  it('should have a different file selected after a file was deleted and re-added', () => {
    // given
    component.shipyardFile = new File(['test content'], 'test1.yaml');

    // when
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('.shipyard-delete-button');
    button.dispatchEvent(new Event('click'));
    fixture.detectChanges();
    component.shipyardFile = new File(['test content'], 'test2.yaml');
    fixture.detectChanges();

    // then
    const div = fixture.nativeElement.querySelector('.shipyard-file-name');
    expect(div.textContent).toEqual('test2.yaml');
  });
});
