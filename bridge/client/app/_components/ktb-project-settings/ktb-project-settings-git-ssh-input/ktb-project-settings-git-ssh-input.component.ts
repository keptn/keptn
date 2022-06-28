import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../../_utils/form.utils';
import { IGitBasicConfiguration } from '../../../../../shared/interfaces/project';

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
  public set gitInputData(data: IGitBasicConfiguration | undefined) {
    if (!data) {
      return;
    }
    this.gitUrlControl.setValue(data.remoteURL);
    this.gitUserControl.setValue(data.user);
  }
  @Output()
  public gitDataChange = new EventEmitter<IGitBasicConfiguration | undefined>();

  private get data(): IGitBasicConfiguration | undefined {
    return this.gitUpstreamForm.valid
      ? {
          remoteURL: this.gitUrlControl.value,
          user: this.gitUserControl.value,
        }
      : undefined;
  }

  public dataChanged(): void {
    this.gitDataChange.emit(this.data);
  }
}
