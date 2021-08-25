import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbWebhookSettingsComponent } from './ktb-webhook-settings.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbProjectSettingsGitComponent', () => {
  let component: KtbWebhookSettingsComponent;
  let fixture: ComponentFixture<KtbWebhookSettingsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbWebhookSettingsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set token value to empty string if git uri and git user are set', () => {
    // given
    component.gitData = {
      remoteURI: 'https://some-repo.git',
      gitUser: 'username',
    };

    // then
    expect(component.gitTokenControl.value).toEqual('');
  });

  it('should not set git token control when only git uri is set', () => {
    // given
    component.gitData = {
      remoteURI: 'https://some-repo.git',
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
    const button = fixture.nativeElement.querySelector('button');

    // then
    expect(button).toBeTruthy();
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

  it('should be a invalid form when only remoteUri and token fields are set', () => {
    // given
    component.gitUrlControl.setValue('https://some-repo.git');
    component.gitTokenControl.setValue('testToken');

    // then
    expect(component.gitUpstreamForm.invalid).toBe(true);
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
      remoteURI: 'https://some-repo.git',
      gitUser: 'username',
    };

    // when
    const spy = jest.spyOn(component.gitDataChanged, 'emit');
    component.gitUrlControl.setValue('https://some-other-repo.git', {emitEvent: true});
    component.onGitUpstreamFormChange();

    // then
    expect(spy).toHaveBeenCalled();
    expect(spy).toHaveBeenCalledWith({gitToken: '', gitFormValid: false, remoteURI: 'https://some-other-repo.git', gitUser: 'username'});
  });
});
