import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitSshInputComponent } from './ktb-project-settings-git-ssh-input.component';
import { AppModule } from '../../app.module';

describe('KtbProjectSettingsGitSshInputComponent', () => {
  let component: KtbProjectSettingsGitSshInputComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitSshInputComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProjectSettingsGitSshInputComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should correctly set input data and emit data', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');

    // when
    component.gitInputData = {
      gitRemoteURL: 'ssh://myGitUrl',
      gitUser: 'myUser',
    };

    // then
    expect(component.gitUrlControl.value).toBe('ssh://myGitUrl');
    expect(component.gitUpstreamForm.controls.gitUser.value).toBe('myUser');
    expect(emitSpy).not.toHaveBeenCalled();

    // when
    component.dataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitRemoteURL: 'ssh://myGitUrl',
      gitUser: 'myUser',
    });
  });

  it('should emit data if data is valid', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    component.gitUrlControl.setValue('ssh://myGitUrl');
    component.gitUpstreamForm.controls.gitUser.setValue('myUser');

    // when
    component.dataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitRemoteURL: 'ssh://myGitUrl',
      gitUser: 'myUser',
    });
  });

  it('should emit data if user is not set', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    component.gitUrlControl.setValue('ssh://myGitUrl');

    // when
    component.dataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith({
      gitRemoteURL: 'ssh://myGitUrl',
      gitUser: '',
    });
  });

  it('should not emit data if URL does not begin with ssh://', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    component.gitUrlControl.setValue('http://myGitUrl');

    // when
    component.dataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should not emit data if URL is not set', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    component.gitUpstreamForm.controls.gitUser.setValue('myUser');

    // when
    component.dataChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });
});
