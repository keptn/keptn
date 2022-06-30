import { ComponentFixture, TestBed } from '@angular/core/testing';
import { TestUtils } from '../../_utils/test.utils';
import { KtbSshKeyInputComponent } from './ktb-ssh-key-input.component';
import { KtbSshKeyInputModule } from './ktb-ssh-key-input.module';

describe('KtbSshKeyInputComponent', () => {
  let component: KtbSshKeyInputComponent;
  let fixture: ComponentFixture<KtbSshKeyInputComponent>;
  const validFileContent = '-----BEGIN OPENSSH PRIVATE KEY-----\n-----END OPENSSH PRIVATE KEY-----';

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbSshKeyInputModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbSshKeyInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set input data correctly', () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');

    // when
    component.sshInput = {
      privateKeyPass: 'myPassphrase',
      privateKey: btoa(validFileContent),
    };

    // then
    expect(emitSpy).not.toHaveBeenCalled();

    // when
    component.sshDataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      privateKeyPass: 'myPassphrase',
      privateKey: btoa(validFileContent),
    });
  });

  it('should emit valid data', () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');
    component.privateKeyControl.setValue(validFileContent);
    component.sshKeyForm.controls.privateKeyPassword.setValue('myPassphrase');

    // when
    component.sshDataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      privateKeyPass: 'myPassphrase',
      privateKey: btoa(validFileContent),
    });
  });

  it('should emit valid data if passphrase is empty', () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');
    component.privateKeyControl.setValue(validFileContent);

    // when
    component.sshDataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      privateKeyPass: '',
      privateKey: btoa(validFileContent),
    });
  });

  it('should emit undefined if privateKey is invalid', () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');
    component.privateKeyControl.setValue('myPrivateKey');

    // when
    component.sshDataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if privateKey is empty', () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');

    // when
    component.sshDataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should set text of file and show error', async () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');

    // when
    await component.validateSshPrivateKey(TestUtils.createFileList('myFile'));

    // then
    expect(component.privateKeyControl.valid).toBe(false);
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should set text of file and not show error', async () => {
    // given
    const emitSpy = jest.spyOn(component.sshDataChange, 'emit');

    // when
    await component.validateSshPrivateKey(TestUtils.createFileList(validFileContent));

    // then
    expect(component.privateKeyControl.valid).toBe(true);
    expect(emitSpy).toHaveBeenCalledWith({
      privateKeyPass: '',
      privateKey: btoa(validFileContent),
    });
  });
});
