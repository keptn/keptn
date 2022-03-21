import { ComponentFixture, TestBed } from '@angular/core/testing';

import { GitFormType, KtbProjectSettingsGitExtendedComponent } from './ktb-project-settings-git-extended.component';
import { AppModule } from '../../app.module';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { of } from 'rxjs';
import { IGitHttps, IGitSsh } from '../../_interfaces/git-upstream';

describe('KtbProjectSettingsGitExtendedComponent', () => {
  let component: KtbProjectSettingsGitExtendedComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitExtendedComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ projectName: 'sockshop' })),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProjectSettingsGitExtendedComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should default select HTTPS', () => {
    expect(component.selectedForm).toBe(GitFormType.HTTPS);
  });

  it('should update selected form to SSH and back to HTTPS', () => {
    // given
    component.gitInputData = getDefaultSshData();
    // then
    expect(component.selectedForm).toBe(GitFormType.SSH);

    // when
    component.gitInputData = undefined;

    // then
    expect(component.selectedForm).toBe(GitFormType.HTTPS);
  });

  it('should select another form and data should be invalidated', () => {
    // given
    const emitSpy = jest.spyOn(component.gitDataChange, 'emit');
    expect(component.selectedForm).toBe(GitFormType.HTTPS);

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

    // when
    component.dataChanged(getDefaultSshData());

    // then
    expect(emitSpy).toHaveBeenCalledWith(getDefaultSshData());
  });

  it('should not update gitUpstream if data is not set', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const updateUpstreamSpy = jest.spyOn(dataService, 'updateGitUpstream');

    // when
    component.updateUpstream();

    // then
    expect(updateUpstreamSpy).not.toHaveBeenCalled();
  });

  it('should update gitUpstream', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const updateUpstreamSpy = jest.spyOn(dataService, 'updateGitUpstream');

    // when
    component.gitData = getDefaultSshData();
    component.updateUpstream();

    // then
    expect(updateUpstreamSpy).toHaveBeenCalled();
  });

  it('should correctly return data if input is HTTPS', () => {
    component.gitInputData = getDefaultHttpsData();
    expect(component.gitInputDataHTTPS).toEqual(getDefaultHttpsData());

    expect(component.gitInputDataSSH).toBe(undefined);
  });

  it('should correctly return data if input is SSH', () => {
    component.gitInputData = getDefaultSshData();
    expect(component.gitInputDataSSH).toEqual(getDefaultSshData());
    expect(component.gitInputDataHTTPS).toBe(undefined);
  });

  it('should return undefined if input is undefined', () => {
    component.gitInputData = undefined;
    expect(component.gitInputDataSSH).toBe(undefined);
    expect(component.gitInputDataHTTPS).toBe(undefined);
  });

  function setSelectedForm(type: GitFormType): void {
    component.setSelectedForm({ value: type } as DtRadioChange<GitFormType>);
  }

  function getDefaultSshData(): IGitSsh {
    return {
      ssh: {
        gitPrivateKeyPass: '',
        gitPrivateKey: '',
        gitRemoteURL: '',
      },
    };
  }

  function getDefaultHttpsData(): IGitHttps {
    return {
      https: {
        gitRemoteURL: '',
        gitToken: '',
      },
    };
  }
});
