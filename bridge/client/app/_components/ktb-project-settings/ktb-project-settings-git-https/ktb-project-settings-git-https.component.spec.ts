import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ComponentFixture, TestBed } from '@angular/core/testing';
import { IGitHTTPSConfiguration, IProxy } from 'shared/interfaces/project';
import { AppUtils } from '../../../_utils/app.utils';
import { KtbProjectSettingsModule } from '../ktb-project-settings.module';
import { KtbProjectSettingsGitHttpsComponent } from './ktb-project-settings-git-https.component';
import { IGitData } from '../ktb-project-settings-git/ktb-project-settings-git.utils';

describe('KtbProjectSettingsGitHttpsComponent', () => {
  let component: KtbProjectSettingsGitHttpsComponent;
  let fixture: ComponentFixture<KtbProjectSettingsGitHttpsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbProjectSettingsModule, HttpClientTestingModule],
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
      url: '0.0.0.0:5000',
      scheme: 'https',
      password: '',
      user: 'myProxyUser',
    };
    expect(component.proxyEnabled).toBe(true);
    expect(component.proxyInput).toEqual(iProxy);
    expect(component.gitInputData).toEqual(getInputDataWithProxy());
    expect(component.certificateInput).toBe(btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'));
    expect(component.gitDataRequired).toEqual({
      user: 'myUser',
      remoteURL: 'https://myGitUrl.com',
      valid: false,
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
      user: 'myUser',
      remoteURL: 'https://myGitUrl.com',
      valid: false,
    });
  });

  it('should only set gitUpstream if the data is valid', () => {
    const invalidGitUpstreams: IGitData[] = [
      {
        token: '',
        remoteURL: '',
        valid: false,
      },
      {
        remoteURL: 'https://myGitUrl.com',
        valid: false,
      },
      {
        token: 'myToken',
        remoteURL: 'myGitUrl.com',
        valid: false,
      },
    ];
    for (const gitData of invalidGitUpstreams) {
      component.gitUpstreamChanged(gitData);
      // eslint-disable-next-line @typescript-eslint/dot-notation
      expect(component['gitUpstream']).toBe(undefined);
    }

    const validUpstreams: IGitData[] = [
      {
        user: 'myUser',
        token: 'myToken',
        remoteURL: 'http://myGitUrl.com',
        valid: true,
      },
      {
        token: 'myToken',
        remoteURL: 'http://myGitUrl.com',
        valid: true,
      },
    ];
    for (const gitData of validUpstreams) {
      component.gitUpstreamChanged(AppUtils.copyObject(gitData)); // just make sure that it isn't a reference. We use the same object to validate it again
      const { valid, ...requiredData } = gitData;
      // eslint-disable-next-line @typescript-eslint/dot-notation
      expect(component['gitUpstream']).toEqual(requiredData);
    }
  });

  it('should remove proxy information when proxy switch is switched', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithProxy();
    component.proxy = {
      url: '0.0.0.0:5000',
      scheme: 'https',
      password: '',
      user: 'myProxyUser',
    };
    expect(component.proxyInput).not.toBe(undefined);
    expect(component.proxy).not.toBe(undefined);

    // when
    component.gitUpstreamChanged({
      user: 'myUser',
      token: 'myToken',
      remoteURL: 'https://myGitUrl.com',
      valid: true,
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
      url: '0.0.0.0:5000',
      scheme: 'https',
      password: '',
      user: 'myProxyUser',
    };
    component.gitUpstreamChanged({
      user: 'myUser',
      token: 'myToken',
      remoteURL: 'https://myGitUrl.com',
      valid: true,
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
    component.gitInputData = getInputDataWithoutProxyAndCertificate('myToken');

    // when
    component.gitUpstreamChanged({
      token: 'myToken',
      remoteURL: 'https://myGitUrl.com',
      valid: true,
    });
    component.inputChanged();

    // then
    expect(component.certificate).toBe(undefined);
    expect(component.proxy).toBe(undefined);
    expect(emitSpy).toHaveBeenCalledWith({
      remoteURL: 'https://myGitUrl.com',
      user: undefined,
      https: {
        certificate: undefined,
        insecureSkipTLS: false,
        token: 'myToken',
      },
    });
  });

  it('should emit undefined if proxy is invalid', () => {
    // given
    const emitSpy = jest.spyOn(component.dataChange, 'emit');
    component.gitInputData = getInputDataWithoutProxy('myToken');

    // when
    component.gitUpstreamChanged({
      token: 'myToken',
      remoteURL: 'https://myGitUrl.com',
      valid: true,
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
      token: 'myToken',
      remoteURL: 'https://myGitUrl.com',
      valid: true,
    });
    component.isCertificateValid = false;
    component.inputChanged();

    // then
    expect(emitSpy).toHaveBeenCalledWith(undefined);
  });

  function getInputDataWithProxy(token = ''): IGitHTTPSConfiguration {
    return {
      remoteURL: 'https://myGitUrl.com',
      user: 'myUser',
      https: {
        token,
        proxy: {
          url: '0.0.0.0:5000',
          scheme: 'https',
          password: '',
          user: 'myProxyUser',
        },
        insecureSkipTLS: false,
        certificate: btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'),
      },
    };
  }

  function getInputDataWithoutProxy(token = ''): IGitHTTPSConfiguration {
    return {
      remoteURL: 'https://myGitUrl.com',
      user: 'myUser',
      https: {
        token,
        insecureSkipTLS: false,
        certificate: btoa('-----BEGIN CERTIFICATE-----\n-----END CERTIFICATE-----'),
      },
    };
  }

  function getInputDataWithoutProxyAndCertificate(token = ''): IGitHTTPSConfiguration {
    return {
      remoteURL: 'https://myGitUrl.com',
      user: 'myUser',
      https: {
        token,
        insecureSkipTLS: false,
      },
    };
  }
});
