import { Component, EventEmitter, Input, Output } from '@angular/core';
import { AppUtils } from '../../../_utils/app.utils';
import { isGitUpstreamValidSet } from '../../../_utils/git-upstream.utils';
import { IGitHTTPSConfiguration, IProxy } from '../../../../../shared/interfaces/project';
import { FormControl } from '@angular/forms';
import { IGitData, IRequiredGitData } from '../ktb-project-settings-git/ktb-project-settings-git.utils';

@Component({
  selector: 'ktb-project-settings-git-https',
  templateUrl: './ktb-project-settings-git-https.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitHttpsComponent {
  private gitUpstream?: IRequiredGitData;
  private _gitInputData?: IGitHTTPSConfiguration;
  public certificateInput?: string;
  public proxyEnabled = false;
  public proxy?: IProxy;
  public certificate?: string;
  public isCertificateValid = true;
  public gitDataRequired: IGitData = { remoteURL: '', valid: false };
  public proxyInput?: IProxy;
  public isInsecureControl = new FormControl(false);

  @Input()
  public isCreateMode = false;
  @Input()
  public isLoading = false;
  @Input()
  public set gitInputData(data: IGitHTTPSConfiguration | undefined) {
    this._gitInputData = data;
    if (data?.https?.proxy) {
      this.proxyEnabled = true;
      this.proxyInput = data.https.proxy;
      this.proxy = AppUtils.copyObject(this.proxyInput);
    }
    this.isInsecureControl.setValue(data?.https?.insecureSkipTLS ?? false);
    this.certificateInput = data?.https?.certificate;
    this.certificate = data?.https?.certificate;
    this.gitDataRequired = {
      user: data?.user,
      remoteURL: data?.remoteURL ?? '',
      valid: false,
    };
  }
  public get gitInputData(): IGitHTTPSConfiguration | undefined {
    return this._gitInputData;
  }
  @Output()
  public dataChange = new EventEmitter<IGitHTTPSConfiguration | undefined>();

  public get data(): IGitHTTPSConfiguration | undefined {
    return this.isValid && this.gitUpstream
      ? {
          remoteURL: this.gitUpstream.remoteURL,
          user: this.gitUpstream.user,
          https: {
            ...(this.proxyEnabled && this.proxy && { proxy: this.proxy }),
            token: this.gitUpstream.token,
            certificate: this.certificate,
            insecureSkipTLS: this.isInsecureControl.value,
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
    const { valid, ...gitUpstream } = data;
    if (valid && isGitUpstreamValidSet(gitUpstream)) {
      this.gitUpstream = gitUpstream;
    } else {
      this.gitUpstream = undefined;
    }
  }
}
