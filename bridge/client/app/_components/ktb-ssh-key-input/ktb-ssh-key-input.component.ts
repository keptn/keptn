import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';
import { IGitSshData } from '../../../../shared/interfaces/project';

@Component({
  selector: 'ktb-ssh-key-input',
  templateUrl: './ktb-ssh-key-input.component.html',
  styleUrls: [],
})
export class KtbSshKeyInputComponent {
  public privateKeyControl = new FormControl('', [Validators.required, FormUtils.isSshKeyValidator]);
  private privateKeyPasswordControl = new FormControl('');
  public sshKeyForm = new FormGroup({
    privateKey: this.privateKeyControl,
    privateKeyPassword: this.privateKeyPasswordControl,
  });
  public dropError?: string;

  @Input()
  public set sshInput(data: IGitSshData | undefined) {
    if (data) {
      this.privateKeyControl.setValue(atob(data.privateKey));
      this.privateKeyPasswordControl.setValue(data.privateKeyPass);
    }
  }

  @Output()
  public sshDataChange = new EventEmitter<IGitSshData | undefined>();

  private get data(): IGitSshData | undefined {
    return this.sshKeyForm.valid
      ? {
          privateKey: btoa(this.privateKeyControl.value),
          privateKeyPass: this.privateKeyPasswordControl.value,
        }
      : undefined;
  }

  public async validateSshPrivateKey(files: FileList | null): Promise<void> {
    const file = files?.[0];
    if (file) {
      this.privateKeyControl.setValue((await file.text()).trim());
      this.privateKeyControl.markAsDirty();
      this.sshDataChanged();
    }
  }

  public sshDataChanged(): void {
    this.sshDataChange.emit(this.data);
  }
}
