import { Component, EventEmitter, Input, Output } from '@angular/core';
import { AppUtils } from '../../../_utils/app.utils';
import { IGitBasicConfiguration, IGitSshConfiguration, IGitSshData } from '../../../../../shared/interfaces/project';

@Component({
  selector: 'ktb-project-settings-git-ssh',
  templateUrl: './ktb-project-settings-git-ssh.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitSshComponent {
  public gitUpstream?: IGitBasicConfiguration;
  public sshKeyData?: IGitSshData;
  public gitInputData?: IGitBasicConfiguration;
  public sshInputData?: IGitSshData;

  @Input()
  public isLoading = false;

  @Input()
  public set gitInputSshData(data: IGitSshConfiguration | undefined) {
    if (data) {
      this.gitInputData = {
        remoteURL: data.remoteURL,
        user: data.user,
      };
      this.gitUpstream = AppUtils.copyObject(this.gitInputData);
      this.sshInputData = {
        privateKeyPass: data.ssh?.privateKeyPass ?? '',
        privateKey: data.ssh?.privateKey ?? '',
      };
    }
  }
  @Output()
  public sshChange = new EventEmitter<IGitSshConfiguration | undefined>();

  public get data(): IGitSshConfiguration | undefined {
    return this.gitUpstream && this.sshKeyData
      ? {
          remoteURL: this.gitUpstream.remoteURL,
          user: this.gitUpstream.user,
          ssh: {
            ...this.sshKeyData,
          },
        }
      : undefined;
  }

  public sshChanged(): void {
    this.sshChange.emit(this.data);
  }
}
