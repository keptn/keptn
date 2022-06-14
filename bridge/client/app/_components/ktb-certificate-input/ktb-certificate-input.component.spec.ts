import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbCertificateInputComponent } from './ktb-certificate-input.component';
import { TestUtils } from '../../_utils/test.utils';
import { KtbCertificateInputModule } from './ktb-certificate-input.module';

describe('KtbCertificateInputComponent', () => {
  let component: KtbCertificateInputComponent;
  let fixture: ComponentFixture<KtbCertificateInputComponent>;
  const validFileContent = '-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----';

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbCertificateInputModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbCertificateInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set text of file and show error', async () => {
    // given
    const emitSpy = jest.spyOn(component.certificateChange, 'emit');

    // when
    await component.validateCertificate(TestUtils.createFileList('myFile'));

    // then
    expect(component.certificateControl.valid).toBe(false);
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should set text of file and not show error', async () => {
    // given
    const emitSpy = jest.spyOn(component.certificateChange, 'emit');

    // when
    await component.validateCertificate(TestUtils.createFileList(validFileContent));

    // then
    expect(component.certificateControl.valid).toBe(true);
    expect(emitSpy).toHaveBeenCalledWith(btoa(validFileContent));
  });

  it('should not allow invalid certificate', () => {
    // given
    const emitSpy = jest.spyOn(component.certificateChange, 'emit');

    // when
    component.certificateControl.setValue('myInvalidValue');
    component.certificateChanged();

    // then
    expect(component.certificateControl.valid).toBe(false);
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should allow empty certificates', () => {
    // given
    const emitSpy = jest.spyOn(component.certificateChange, 'emit');

    // when
    component.certificateControl.setValue('');
    component.certificateChanged();

    // then
    expect(component.certificateControl.valid).toBe(true);
    expect(emitSpy).toHaveBeenCalledWith('');
  });

  it('should be valid certificate', () => {
    // given
    const emitSpy = jest.spyOn(component.certificateChange, 'emit');

    // when
    component.certificateControl.setValue(validFileContent);
    component.certificateChanged();

    // then
    expect(component.certificateControl.valid).toBe(true);
    expect(emitSpy).toHaveBeenCalledWith(btoa(validFileContent));
  });
});
