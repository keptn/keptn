import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { of } from 'rxjs';
import { ApiService } from '../../../_services/api.service';
import { ApiServiceMock } from '../../../_services/api.service.mock';
import { KtbProjectSettingsModule } from '../ktb-project-settings.module';

import { GitFormType, KtbProjectSettingsGitExtendedComponent } from './ktb-project-settings-git-extended.component';
import { IGitHttpsConfiguration, IGitSshConfiguration } from '../../../../../shared/interfaces/project';

describe('KtbProjectSettingsGitExtendedComponent', () => {
  let component: KtbProjectSettingsGitExtendedComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitExtendedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectSettingsModule, HttpClientTestingModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ projectName: 'sockshop' })),
          },
        },
        {
          provide: ApiService,
          useClass: ApiServiceMock,
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProjectSettingsGitExtendedComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should default select HTTPS', () => {
    // given, when
    fixture.detectChanges();

    // then
    expect(component.selectedForm).toBe(GitFormType.HTTPS);
  });

  it('should default select NO_UPSTREAM if git upstream is not required', () => {
    // given
    component.gitUpstreamRequired = false;

    // when
    fixture.detectChanges();

    // then
    expect(component.selectedForm).toBe(GitFormType.NO_UPSTREAM);
  });

  it('should select HTTPS form on init with https data given', () => {
    // given
    component.gitInputData = getDefaultHttpsData();

    // when
    fixture.detectChanges();

    // then
    expect(component.selectedForm).toBe(GitFormType.HTTPS);
  });

  it('should select SSH form on init with ssh data given', () => {
    // given
    component.gitInputData = getDefaultSshData();

    // when
    fixture.detectChanges();

    // then
    expect(component.selectedForm).toBe(GitFormType.SSH);
  });

  it('should select NO_UPSTREAM form on init if not data given and git upstream is not required', () => {
    // given
    component.gitUpstreamRequired = false;
    component.gitInputData = getDefaultHttpsData();
    component.gitInputData.remoteURL = '';

    // when
    fixture.detectChanges();

    // then
    expect(component.selectedForm).toBe(GitFormType.NO_UPSTREAM);
  });

  it('should select another form and data should be invalidated', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    fixture.detectChanges();

    // when
    setSelectedForm(GitFormType.SSH);

    // then
    expect(component.selectedForm).toBe(GitFormType.SSH);
    expect(component.gitData).toBe(undefined);
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should correctly update and emit data', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    fixture.detectChanges();

    // when
    component.dataChanged(GitFormType.SSH, getDefaultSshData());

    // then
    expect(emitSpy).toHaveBeenCalledWith(getDefaultSshData());
  });

  it('should correctly return data if input is HTTPS', () => {
    // given
    component.gitInputData = getDefaultHttpsData();

    // when
    fixture.detectChanges();

    // then
    expect(component.gitInputDataHttps).toEqual(getDefaultHttpsData());
    expect(component.gitInputDataSsh).toBe(undefined);
  });

  it('should correctly return data if input is SSH', () => {
    // given
    component.gitInputData = getDefaultSshData();

    // when
    fixture.detectChanges();

    // then
    expect(component.gitInputDataSsh).toEqual(getDefaultSshData());
    expect(component.gitInputDataHttps).toBe(undefined);
  });

  it('should correctly return data if no upstream is selected', () => {
    // given
    const spy = jest.spyOn(component.gitDataChange, 'emit');
    component.gitUpstreamRequired = false;
    fixture.detectChanges();

    // when
    setSelectedForm(GitFormType.NO_UPSTREAM);

    // then
    expect(component.gitInputDataSsh).toEqual(undefined);
    expect(component.gitInputDataHttps).toBe(undefined);
    expect(spy).toHaveBeenCalledWith(null);
  });

  it('should return undefined if input is undefined', () => {
    // given
    component.gitInputData = undefined;

    // when
    fixture.detectChanges();

    // then
    expect(component.gitInputDataSsh).toBe(undefined);
    expect(component.gitInputDataHttps).toBe(undefined);
  });

  it('should emit cached https data if it switched back', () => {
    // given
    fixture.detectChanges();
    component.dataChanged(GitFormType.HTTPS, getDefaultHttpsData());

    // when
    setSelectedForm(GitFormType.SSH);
    expect(component.gitData).toBe(undefined);

    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    setSelectedForm(GitFormType.HTTPS);

    // then
    expect(component.gitData).toEqual(getDefaultHttpsData());
    expect(emitSpy).toHaveBeenCalledWith(getDefaultHttpsData());
  });

  it('should emit cached ssh data if it switched back', () => {
    // given
    fixture.detectChanges();
    setSelectedForm(GitFormType.SSH);
    component.dataChanged(GitFormType.SSH, getDefaultSshData());

    // when
    setSelectedForm(GitFormType.HTTPS);
    expect(component.gitData).toBe(undefined);
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    setSelectedForm(GitFormType.SSH);

    // then
    expect(component.gitData).toEqual(getDefaultSshData());
    expect(emitSpy).toHaveBeenCalledWith(getDefaultSshData());
  });

  function setSelectedForm(type: GitFormType): void {
    component.setSelectedForm({ value: type } as DtRadioChange<GitFormType>);
  }
});

export function getDefaultSshData(): IGitSshConfiguration {
  return {
    remoteURL: 'ssh://git@github.com/keptn/keptn',
    ssh: {
      privateKeyPass: '',
      privateKey: '',
    },
  };
}

export function getDefaultHttpsData(): IGitHttpsConfiguration {
  return {
    remoteURL: 'https://github.com/keptn/keptn',
    https: {
      token: '',
      insecureSkipTLS: false,
    },
  };
}
