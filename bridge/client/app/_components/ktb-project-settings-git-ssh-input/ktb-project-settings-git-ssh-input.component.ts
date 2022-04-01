import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { IGitSshData } from '../../_interfaces/git-upstream';

@Component({
  selector: 'ktb-project-settings-git-ssh-input',
  templateUrl: './ktb-project-settings-git-ssh-input.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitSshInputComponent {
  public gitUrlControl = new FormControl('', [Validators.required, FormUtils.isSshValidator]);
  private gitUserControl = new FormControl('');
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
  });

  @Input()
  public isLoading = false;

  @Input()
  public set gitInputData(data: IGitSshData | undefined) {
    if (data) {
      this.gitUrlControl.setValue(data.gitRemoteURL);
      if (data.gitRemoteURL) {
        this.gitUrlControl.markAsDirty();
      }
      this.gitUserControl.setValue(data.gitUser);
      this.dataChanged();
    }
  }
  @Output()
  public gitDataChange = new EventEmitter<IGitSshData | undefined>();

  private get data(): IGitSshData | undefined {
    return this.gitUpstreamForm.valid
      ? {
          gitRemoteURL: this.gitUrlControl.value,
          gitUser: this.gitUserControl.value,
        }
      : undefined;
  }

  public dataChanged(): void {
    this.gitDataChange.emit(this.data);
  }
}
