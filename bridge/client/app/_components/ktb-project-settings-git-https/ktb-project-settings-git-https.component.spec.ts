import { ComponentFixture, TestBed } from '@angular/core/testing';

import { KtbProjectSettingsGitHttpsComponent } from './ktb-project-settings-git-https.component';
import { IGitData, IGitHttps, IProxy } from '../../_interfaces/git-upstream';
import { AppModule } from '../../app.module';
import { AppUtils } from '../../_utils/app.utils';

describe('KtbProjectSettingsGitHttpsComponent', () => {
  let component: KtbProjectSettingsGitHttpsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitHttpsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [AppModule],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbProjectSettingsGitHttpsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set input data with proxy correctly', () => {
    component.gitInputData = getInputDataWithProxy();
    const iProxy: IProxy = {
      gitProxyUrl: '0.0.0.0',
      gitProxyScheme: 'https',
      gitProxyInsecure: false,
      gitProxyPassword: '',
      gitProxyUser: 'myProxyUser',
    };
    expect(component.proxyEnabled).toBe(true);
    expect(component.proxyInput).toEqual(iProxy);
    expect(component.gitInputData).toEqual(getInputDataWithProxy());
    expect(component.certificateInput).toBe(btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'));
    expect(component.gitDataRequired).toEqual({
      gitUser: 'myUser',
      gitRemoteURL: 'https://myGitUrl.com',
      gitToken: '',
    });
  });

  it('should have undefined proxy if input data does not contain proxy information', () => {
    // given

    // when
    component.gitInputData = getInputDataWithoutProxy();

    // then
    expect(component.proxyEnabled).toBe(false);
    expect(component.proxy).toBe(undefined);
    expect(component.gitInputData).toEqual(getInputDataWithoutProxy());
    expect(component.certificateInput).toBe(btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'));
    expect(component.gitDataRequired).toEqual({
      gitUser: 'myUser',
      gitRemoteURL: 'https://myGitUrl.com',
      gitToken: '',
    });
  });

  it('should only set gitUpstream if the data is valid', () => {
    const invalidGitUpstreams: IGitData[] = [
      {},
      {
        gitToken: '',
        gitRemoteURL: '',
      },
      {
        gitToken: '',
        gitRemoteURL: '',
        gitFormValid: false,
      },
      {
        gitRemoteURL: 'https://myGitUrl.com',
        gitFormValid: false,
      },
      {
        gitToken: 'myToken',
        gitUser: 'myUser',
        gitFormValid: false,
      },
      {
        gitToken: 'myToken',
        gitRemoteURL: 'myGitUrl.com',
        gitFormValid: false,
      },
      {},
    ];
    for (const gitData of invalidGitUpstreams) {
      component.gitUpstreamChanged(gitData);
      // eslint-disable-next-line @typescript-eslint/dot-notation
      expect(component['gitUpstream']).toBe(undefined);
    }

    const validUpstreams: IGitData[] = [
      {
        gitUser: 'myUser',
        gitToken: 'myToken',
        gitRemoteURL: 'http://myGitUrl.com',
        gitFormValid: true,
      },
      {
        gitToken: 'myToken',
        gitRemoteURL: 'http://myGitUrl.com',
        gitFormValid: true,
      },
    ];
    for (const gitData of validUpstreams) {
      component.gitUpstreamChanged(AppUtils.copyObject(gitData)); // just make sure that it isn't a reference. We use the same object to validate it again
      const { gitFormValid, ...requiredData } = gitData;
      // eslint-disable-next-line @typescript-eslint/dot-notation
      expect(component['gitUpstream']).toEqual(requiredData);
    }
  });

  it('should remove proxy information when proxy switch is switched', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithProxy();
    component.proxy = {
      gitProxyUrl: '0.0.0.0',
      gitProxyScheme: 'https',
      gitProxyInsecure: false,
      gitProxyPassword: '',
      gitProxyUser: 'myProxyUser',
    };
    expect(component.proxyInput).not.toBe(undefined);
    expect(component.proxy).not.toBe(undefined);

    // when
    component.gitUpstreamChanged({
      gitUser: 'myUser',
      gitToken: 'myToken',
      gitRemoteURL: 'https://myGitUrl.com',
      gitFormValid: true,
    });
    component.certificate = btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----');
    component.proxyEnabled = false;
    component.inputChanged();

    // then
    expect(component.proxy).not.toBe(undefined);
    expect(emitSpy).toHaveBeenCalledWith(getInputDataWithoutProxy('myToken'));
  });

  it('should add cached proxy information when proxy switch is disabled and enabled again', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithProxy();

    // when
    component.proxy = {
      gitProxyUrl: '0.0.0.0',
      gitProxyScheme: 'https',
      gitProxyInsecure: false,
      gitProxyPassword: '',
      gitProxyUser: 'myProxyUser',
    };
    component.gitUpstreamChanged({
      gitUser: 'myUser',
      gitToken: 'myToken',
      gitRemoteURL: 'https://myGitUrl.com',
      gitFormValid: true,
    });
    component.certificate = btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----');
    component.proxyEnabled = false;
    component.inputChanged();
    component.proxyEnabled = true;
    component.inputChanged();

    // then
    expect(component.proxy).not.toBe(undefined);
    expect(emitSpy).toHaveBeenCalledWith(getInputDataWithProxy('myToken'));
  });

  it('should emit data if proxy and certificate are not set', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithoutProxy('myToken');

    // when
    component.gitUpstreamChanged({
      gitToken: 'myToken',
      gitRemoteURL: 'https://myGitUrl.com',
      gitFormValid: true,
    });
    component.inputChanged();

    // then
    expect(component.certificate).toBe(undefined);
    expect(component.proxy).toBe(undefined);
    expect(emitSpy).toHaveBeenCalledWith({
      https: {
        gitRemoteURL: 'https://myGitUrl.com',
        gitToken: 'myToken',
      },
    });
  });

  it('should emit undefined if proxy is invalid', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithoutProxy('myToken');

    // when
    component.gitUpstreamChanged({
      gitToken: 'myToken',
      gitRemoteURL: 'https://myGitUrl.com',
      gitFormValid: true,
    });
    component.proxyEnabled = true;
    component.inputChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  it('should emit undefined if certificate is invalid', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithoutProxy('myToken');

    // when
    component.gitUpstreamChanged({
      gitToken: 'myToken',
      gitRemoteURL: 'https://myGitUrl.com',
      gitFormValid: true,
    });
    component.isCertificateValid = false;
    component.inputChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  function getInputDataWithProxy(gitToken = ''): IGitHttps {
    return {
      https: {
        gitRemoteURL: 'https://myGitUrl.com',
        gitToken,
        gitProxyUrl: '0.0.0.0',
        gitProxyScheme: 'https',
        gitProxyInsecure: false,
        gitProxyPassword: '',
        gitUser: 'myUser',
        gitProxyUser: 'myProxyUser',
        gitPemCertificate: btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'),
      },
    };
  }

  function getInputDataWithoutProxy(gitToken = ''): IGitHttps {
    return {
      https: {
        gitRemoteURL: 'https://myGitUrl.com',
        gitToken,
        gitUser: 'myUser',
        gitPemCertificate: btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'),
      },
    };
  }
});
