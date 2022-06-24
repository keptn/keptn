import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { FormUtils } from '../../../_utils/form.utils';
import { IGitData } from './ktb-project-settings-git.utils';

@Component({
  selector: 'ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss'],
})
export class KtbProjectSettingsGitComponent {
  private originalGitData: Omit<IGitData, 'valid'> | undefined;

  @Input()
  public isGitUpstreamInProgress = false;

  @Input()
  public isCreateMode = false;

  @Input()
  set isLoading(isLoading: boolean | undefined) {
    if (!this.isCreateMode && !!isLoading) {
      this.gitUrlControl.disable();
      this.gitUserControl.disable();
      this.gitTokenControl.disable();
    } else {
      this.gitUrlControl.enable();
      this.gitUserControl.enable();
      this.gitTokenControl.enable();
    }
    this._isLoading = isLoading;
  }

  get isLoading(): boolean | undefined {
    return this._isLoading;
  }

  @Input()
  set gitData(gitData: IGitData) {
    if (!this.originalGitData) {
      this.originalGitData = {
        remoteURL: gitData.remoteURL,
        user: gitData.user,
      };
    }
    this.resetForm(gitData);
  }

  @Output()
  public gitDataChanged: EventEmitter<IGitData> = new EventEmitter();

  public gitUrlControl = new FormControl('', [Validators.required, FormUtils.isUrlValidator]);
  public gitUserControl = new FormControl('');
  public gitTokenControl = new FormControl('', [Validators.required]);
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl,
  });
  private _isLoading: boolean | undefined;

  public onGitUpstreamFormChange(): void {
    this.gitDataChanged.emit({
      remoteURL: this.gitUrlControl.value,
      user: this.gitUserControl.value,
      token: this.gitTokenControl.value,
      valid: !this.isButtonDisabled(),
    });
  }

  public isButtonDisabled(): boolean {
    return this.gitUpstreamForm.invalid || !this.gitUpstreamForm.dirty || this.isGitUpstreamInProgress;
  }

  public reset(): void {
    this.resetForm(this.originalGitData);
  }

  private resetForm(gitData: Omit<IGitData, 'valid'> | undefined): void {
    this.resetControl(this.gitUrlControl, gitData?.remoteURL || '');
    this.resetControl(this.gitUserControl, gitData?.user || '');
    this.resetControl(this.gitTokenControl, '');
  }

  private resetControl(control: FormControl, value: string): void {
    control.setValue(value);
    control.markAsUntouched();
    control.markAsPristine();
  }
}
