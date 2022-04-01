import { Component, EventEmitter, Input, Output } from '@angular/core';
import { IGitSsh, IGitSshData, ISshKeyData } from '../../_interfaces/git-upstream';

@Component({
  selector: 'ktb-project-settings-git-ssh',
  templateUrl: './ktb-project-settings-git-ssh.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitSshComponent {
  public gitUpstream?: IGitSshData;
  public sshKeyData?: ISshKeyData;
  public gitInputData?: IGitSshData;
  public sshInputData?: ISshKeyData;

  @Input()
  public isLoading = false;

  @Input()
  public set gitInputSshData(data: IGitSsh | undefined) {
    if (data) {
      this.gitInputData = {
        gitRemoteURL: data.ssh.gitRemoteURL,
        gitUser: data.ssh.gitUser,
      };
      this.sshInputData = {
        gitPrivateKeyPass: data.ssh.gitPrivateKeyPass,
        gitPrivateKey: data.ssh.gitPrivateKey,
      };
    }
  }
  @Output()
  public sshChange = new EventEmitter<IGitSsh | undefined>();

  public get data(): IGitSsh | undefined {
    return this.gitUpstream && this.sshKeyData
      ? {
          ssh: {
            ...this.gitUpstream,
            ...this.sshKeyData,
          },
        }
      : undefined;
  }

  public sshChanged(): void {
    this.sshChange.emit(this.data);
  }
}
