import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { IGitData } from '../../_interfaces/git-upstream';

@Component({
  selector: 'ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss'],
})
export class KtbProjectSettingsGitComponent implements OnInit {
  private originalGitData: IGitData | undefined;

  @Input()
  public isGitUpstreamInProgress = false;

  @Input() isInGitExtended = false;

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
    if (!this.originalGitData && gitData.gitRemoteURL && gitData.gitUser) {
      this.originalGitData = {
        gitRemoteURL: gitData.gitRemoteURL,
        gitUser: gitData.gitUser,
      };
    }
    this.resetForm(gitData);
  }

  @Output()
  public gitUpstreamSubmit: EventEmitter<IGitData> = new EventEmitter();

  @Output()
  public gitDataChanged: EventEmitter<IGitData> = new EventEmitter();

  public gitUrlControl = new FormControl('');
  public gitUserControl = new FormControl('');
  public gitTokenControl = new FormControl('');
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl,
  });
  private _isLoading: boolean | undefined;
  public get required(): boolean {
    return !this.isCreateMode || this.isInGitExtended;
  }

  public ngOnInit(): void {
    if (this.required) {
      this.gitUrlControl.setValidators([Validators.required]);
      this.gitTokenControl.setValidators([Validators.required]);
    }
  }

  public setGitUpstream(): void {
    this.gitUpstreamSubmit.emit({
      gitRemoteURL: this.gitUrlControl.value,
      gitUser: this.gitUserControl.value,
      gitToken: this.gitTokenControl.value,
    });
    this.gitTokenControl.markAsUntouched();
    this.gitTokenControl.markAsPristine();
  }

  public onGitUpstreamFormChange(): void {
    this.gitDataChanged.emit({
      gitRemoteURL: this.gitUrlControl.value,
      gitUser: this.gitUserControl.value,
      gitToken: this.gitTokenControl.value,
      gitFormValid: !this.isButtonDisabled(),
    });
  }

  public isButtonDisabled(): boolean {
    return this.gitUpstreamForm.invalid || !this.gitUpstreamForm.dirty || this.isGitUpstreamInProgress;
  }

  public reset(): void {
    this.resetForm(this.originalGitData);
  }

  private resetForm(gitData: IGitData | undefined): void {
    this.gitUrlControl.setValue(gitData?.gitRemoteURL || '');
    this.gitUrlControl.markAsUntouched();
    this.gitUrlControl.markAsPristine();
    this.gitUserControl.setValue(gitData?.gitUser || '');
    this.gitUserControl.markAsUntouched();
    this.gitUserControl.markAsPristine();
    this.gitTokenControl.setValue('');
    this.gitTokenControl.markAsUntouched();
    this.gitTokenControl.markAsPristine();
  }
}
