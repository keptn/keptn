import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectSettingsGitSshComponent } from './ktb-project-settings-git-ssh.component';
import { AppModule } from '../../app.module';

describe('KtbProjectSettingsGitSshComponent', () => {
  let component: KtbProjectSettingsGitSshComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitSshComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
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
      ssh: {
        gitUser: 'myUser',
        gitRemoteURL: 'ssh://myGitUrl',
        gitPrivateKeyPass: '',
        gitPrivateKey: '',
      },
    };

    // then
    expect(component.gitInputData).toEqual({
      gitUser: 'myUser',
      gitRemoteURL: 'ssh://myGitUrl',
    });
    expect(component.sshInputData).toEqual({
      gitPrivateKeyPass: '',
      gitPrivateKey: '',
    });
  });

  it('should emit data if data is not undefined', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given
    component.gitUpstream = {
      gitUser: 'myUser',
      gitRemoteURL: 'ssh://myGitUrl',
    };
    component.sshKeyData = {
      gitPrivateKeyPass: 'myPassphrase',
      gitPrivateKey: 'myPrivateKey',
    };

    // when
    component.sshChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      ssh: {
        gitUser: 'myUser',
        gitRemoteURL: 'ssh://myGitUrl',
        gitPrivateKeyPass: 'myPassphrase',
        gitPrivateKey: 'myPrivateKey',
      },
    });
  });

  it('should emit undefined if data is changed and upstream is not set', () => {
    const emitSpy = jest.spyOn(component.sshChange, 'emit');
    // given
    component.sshKeyData = {
      gitPrivateKeyPass: 'myPassphrase',
      gitPrivateKey: 'myPrivateKey',
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
      gitPrivateKeyPass: 'myPassphrase',
      gitPrivateKey: 'myPrivateKey',
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
