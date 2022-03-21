import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitComponent } from './ktb-project-settings-git.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { IGitData } from '../../_interfaces/git-upstream';

describe('KtbProjectSettingsGitComponent', () => {
  let component: KtbProjectSettingsGitComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbProjectSettingsGitComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set token value to empty string if git uri and git user are set', () => {
    // given
    component.gitData = {
      gitRemoteURL: 'https://some-repo.git',
      gitUser: 'username',
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not set git token control when only git uri is set', () => {
    // given
    component.gitData = {
      gitRemoteURL: 'https://some-repo.git',
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not set git token control when only git user is set', () => {
    // given
    // when
    component.gitData = {
      gitUser: 'username',
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not show "Set Git upstream" button when create mode is true', () => {
    // given
    component.isCreateMode = true;

    // when
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button');

    // then
    expect(button).toBeFalsy();
  });

  it('should show "Set Git upstream" button when create mode is false', () => {
    // given
    component.isCreateMode = false;

    // when
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button'); //TODO: move to UI test

    // then
    expect(button).toBeTruthy();
  });

  it('should disable the inputs when loading in edit mode', () => {
    // given, when
    component.isCreateMode = false;
    component.isLoading = true;

    // then
    assertDisabledInputs(true);
  });

  it('should not disable the inputs when isLoading is not given', () => {
    // given, when
    component.isCreateMode = false;
    component.isLoading = undefined;

    // then
    assertDisabledInputs(false);
  });

  it('should not disable the inputs when isLoading is given but createMode is true', () => {
    // given, when
    component.isCreateMode = true;
    component.isLoading = true;

    // then
    assertDisabledInputs(false);
  });

  it('should enable the buttons when onChange isLoading is false', () => {
    // given, when
    component.isCreateMode = false;
    component.isLoading = true;

    // then
    assertDisabledInputs(true);

    // given, when
    component.isLoading = false;

    // then
    assertDisabledInputs(false);
  });

  it('should be an invalid form when no fields are set', () => {
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a invalid form when only remoteUri field is set', () => {
    // given
    component.gitUrlControl.setValue('https://some-repo.git');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a invalid form when only user field is set', () => {
    // given
    component.gitUserControl.setValue('username');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a invalid form when only token field is set', () => {
    // given
    component.gitTokenControl.setValue('testToken');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a invalid form when only remoteUri and username fields are set', () => {
    // given
    component.gitUserControl.setValue('username');
    component.gitUrlControl.setValue('https://some-repo.git');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a valid form when only remoteUri and token fields are set', () => {
    // given
    component.gitUrlControl.setValue('https://some-repo.git');
    component.gitTokenControl.setValue('testToken');

    // then
    expect(component.gitUpstreamForm.valid).toBe(true);
  });

  it('should be a invalid form when only username and token fields are set', () => {
    // given
    component.gitUserControl.setValue('username');
    component.gitTokenControl.setValue('testToken');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
  });

  it('should be a valid form when remoteUri, username and token fields are set', () => {
    // given
    component.gitUserControl.setValue('username');
    component.gitUrlControl.setValue('https://some-repo.git');
    component.gitTokenControl.setValue('testToken');

    // then
    expect(component.gitUpstreamForm.valid).toBe(true);
  });

  it('should show a disabled button when form is invalid', () => {
    // given
    component.isCreateMode = false;
    component.gitUserControl.setValue('username');
    component.gitTokenControl.setValue('testToken');
    component.gitUpstreamForm.markAsDirty();

    // when
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button');

    // then
    expect(button.disabled).toBeTruthy();
  });

  it('should show an enabled button when form is valid', () => {
    // given
    component.isCreateMode = false;
    component.gitUserControl.setValue('username');
    component.gitUrlControl.setValue('https://some-repo.git');
    component.gitTokenControl.setValue('testToken');
    component.gitUpstreamForm.markAsDirty();

    // when
    // component.gitUpstreamForm.updateValueAndValidity();
    fixture.detectChanges();
    const button = fixture.nativeElement.querySelector('button');

    // then
    expect(button.disabled).toBeFalsy();
  });

  it('should emit the changed git data when form is changed', () => {
    // given
    component.gitData = {
      gitRemoteURL: 'https://some-repo.git',
      gitUser: 'username',
    };

    // when
    const spy = jest.spyOn(component.gitDataChanged, 'emit');
    component.gitUrlControl.setValue('https://some-other-repo.git', { emitEvent: true });
    component.onGitUpstreamFormChange();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith({
      gitToken: '',
      gitFormValid: false,
      gitRemoteURL: 'https://some-other-repo.git',
      gitUser: 'username',
    } as IGitData);
  });

  it('should submit/emit form', () => {
    // given
    const emitSpy = jest.spyOn(component.gitUpstreamSubmit, 'emit');
    const url = 'https://my-git-repo.git';
    const user = 'myUser';
    const token = 'myToken';
    component.gitUrlControl.setValue(url);
    component.gitUserControl.setValue(user);
    component.gitTokenControl.setValue(token);

    // when
    component.setGitUpstream();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitRemoteURL: url,
      gitUser: user,
      gitToken: token,
    });
  });

  it('should reset form to input values', () => {
    // given
    component.gitData = {
      gitRemoteURL: 'https://some-repo.git',
      gitUser: 'username',
    };

    // when
    const url = 'https://my-git-repo.git';
    const user = 'myUser';
    const token = 'myToken';
    component.gitUrlControl.setValue(url);
    component.gitUserControl.setValue(user);
    component.gitTokenControl.setValue(token);

    component.reset();

    expect(component.gitUrlControl.value).toBe('https://some-repo.git');
    expect(component.gitUserControl.value).toBe('username');
    expect(component.gitTokenControl.value).toBe('');
  });

  function assertDisabledInputs(isDisabled: boolean): void {
    expect(component.gitUrlControl.disabled).toBe(isDisabled);
    expect(component.gitUserControl.disabled).toBe(isDisabled);
    expect(component.gitTokenControl.disabled).toBe(isDisabled);
  }
});
