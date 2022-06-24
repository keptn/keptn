import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectSettingsModule } from '../ktb-project-settings.module';
import { KtbProjectSettingsGitSshComponent } from './ktb-project-settings-git-ssh.component';

describe('KtbProjectSettingsGitSshComponent', () => {
  let component: KtbProjectSettingsGitSshComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitSshComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectSettingsModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProjectSettingsGitSshComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set input data correctly', () => {
    // given
    // when
    component.gitInputSshData = {
      user: 'myUser',
      remoteURL: 'ssh://myGitUrl',
      ssh: {
        privateKeyPass: '',
        privateKey: '',
      },
    };

    // then
    expect(component.gitInputData).toEqual({
      user: 'myUser',
      remoteURL: 'ssh://myGitUrl',
    });
    expect(component.sshInputData).toEqual({
      privateKeyPass: '',
      privateKey: '',
    });
  });

  it('should emit data if data is not undefined', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given
    component.gitUpstream = {
      user: 'myUser',
      remoteURL: 'ssh://myGitUrl',
    };
    component.sshKeyData = {
      privateKeyPass: 'myPassphrase',
      privateKey: 'myPrivateKey',
    };

    // when
    component.sshChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      user: 'myUser',
      remoteURL: 'ssh://myGitUrl',
      ssh: {
        privateKeyPass: 'myPassphrase',
        privateKey: 'myPrivateKey',
      },
    });
  });

  it('should emit undefined if data is changed and upstream is not set', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given
    component.sshKeyData = {
      privateKeyPass: 'myPassphrase',
      privateKey: 'myPrivateKey',
    };

    // when
    component.sshChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if data is changed and key data is not set', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given
    component.sshKeyData = {
      privateKeyPass: 'myPassphrase',
      privateKey: 'myPrivateKey',
    };

    // when
    component.sshChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if data is changed and neither key data or upstream is set', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given

    // when
    component.sshChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });
});
