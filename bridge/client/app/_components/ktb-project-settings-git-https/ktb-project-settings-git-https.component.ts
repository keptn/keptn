import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IGitData, IGitHttps, IProxy, IRequiredGitData } from '../../_interfaces/git-upstream';
import { isGitUpstreamValidSet, isGitWithProxy } from '../../_utils/git-upstream.utils';

@Component({
  selector: 'ktb-project-settings-git-https',
  templateUrl: './ktb-project-settings-git-https.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitHttpsComponent {
  private gitUpstream?: IRequiredGitData;
  private _gitInputData?: IGitHttps;
  public certificateInput?: string;
  public proxyEnabled = false;
  public proxy?: IProxy;
  public certificate?: string;
  public isCertificateValid = true;
  public gitDataRequired: IGitData = {};
  public proxyInput?: IProxy;

  @Input()
  public isCreateMode = false;
  @Input()
  public isLoading = false;
  @Input()
  public set gitInputData(data: IGitHttps | undefined) {
    this._gitInputData = data;
    if (data && isGitWithProxy(data)) {
      this.proxyEnabled = true;
      this.proxyInput = {
        gitProxyUrl: data.https.gitProxyUrl,
        gitProxyInsecure: data.https.gitProxyInsecure,
        gitProxyScheme: data.https.gitProxyScheme,
        gitProxyUser: data.https.gitProxyUser,
        gitProxyPassword: data.https.gitProxyPassword,
      };
    }
    this.certificateInput = data?.https.gitPemCertificate;
    this.gitDataRequired = {
      gitUser: data?.https.gitUser,
      gitToken: data?.https.gitToken,
      gitRemoteURL: data?.https.gitRemoteURL,
    };
  }
  public get gitInputData(): IGitHttps | undefined {
    return this._gitInputData;
  }
  @Output()
  public dataChange = new EventEmitter<IGitHttps | undefined>();

  public get data(): IGitHttps | undefined {
    return this.isValid && this.gitUpstream
      ? {
          https: {
            ...this.gitUpstream,
            ...(this.proxyEnabled && this.proxy),
            gitPemCertificate: this.certificate,
          },
        }
      : undefined;
  }

  private get isProxyValid(): boolean {
    return !this.proxyEnabled || !!this.proxy;
  }

  public get isValid(): boolean {
    return this.isProxyValid && this.isCertificateValid && !!this.gitUpstream;
  }

  public inputChanged(): void {
    this.dataChange.emit(this.data);
  }

  public gitUpstreamChanged(data: IGitData): void {
    const { gitFormValid, ...gitUpstream } = data;
    if (gitFormValid && isGitUpstreamValidSet(gitUpstream)) {
      this.gitUpstream = gitUpstream;
    } else {
      this.gitUpstream = undefined;
    }
  }
}
