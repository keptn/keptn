import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbProjectSettingsModule } from '../ktb-project-settings.module';
import { KtbProjectSettingsGitComponent } from './ktb-project-settings-git.component';
import { IGitData } from './ktb-project-settings-git.utils';

describe('KtbProjectSettingsGitComponent', () => {
  let component: KtbProjectSettingsGitComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectSettingsModule, HttpClientTestingModule],
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
      remoteURL: 'https://some-repo.git',
      user: 'username',
      valid: false,
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not set git token control when only git uri is set', () => {
    // given
    component.gitData = {
      remoteURL: 'https://some-repo.git',
      valid: false,
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not set git token control when only git user is set', () => {
    // given
    // when
    component.gitData = {
      user: 'username',
      valid: false,
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

  it('should emit the changed git data when form is changed', () => {
    // given
    component.gitData = {
      remoteURL: 'https://some-repo.git',
      user: 'username',
      valid: false,
    };

    // when
    const spy = jest.spyOn(component.gitDataChanged, 'emit');
    component.gitUrlControl.setValue('https://some-other-repo.git', { emitEvent: true });
    component.onGitUpstreamFormChange();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith({
      token: '',
      valid: false,
      remoteURL: 'https://some-other-repo.git',
      user: 'username',
    } as IGitData);
  });

  it('should reset form to input values', () => {
    // given
    component.gitData = {
      remoteURL: 'https://some-repo.git',
      user: 'username',
      valid: false,
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
