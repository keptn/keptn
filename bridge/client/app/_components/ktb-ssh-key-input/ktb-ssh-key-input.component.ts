import { Component, EventEmitter, Input, Output } from '@angular/core';
import { ISshKeyData } from '../../_interfaces/git-upstream';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../_utils/form.utils';

@Component({
  selector: 'ktb-ssh-key-input',
  templateUrl: './ktb-ssh-key-input.component.html',
  styleUrls: ['./ktb-ssh-key-input.component.scss'],
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
  public set sshInput(data: ISshKeyData | undefined) {
    if (data) {
      this.privateKeyControl.setValue(atob(data.gitPrivateKey));
      if (data.gitPrivateKey) {
        this.privateKeyControl.markAsDirty();
      }
      this.privateKeyPasswordControl.setValue(data.gitPrivateKeyPass);
      this.sshDataChanged();
    }
  }

  @Output()
  public sshDataChange = new EventEmitter<ISshKeyData | undefined>();

  private get data(): ISshKeyData | undefined {
    return this.sshKeyForm.valid
      ? {
          gitPrivateKey: btoa(this.privateKeyControl.value),
          gitPrivateKeyPass: this.privateKeyPasswordControl.value,
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
