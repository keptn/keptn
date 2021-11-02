import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { FormControl, FormGroup, Validators } from '@angular/forms';

@Component({
  selector: 'ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss'],
})
export class KtbProjectSettingsGitComponent implements OnInit {
  private originalGitData: GitData | undefined;

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
  set gitData(gitData: GitData) {
    if (!this.originalGitData && gitData.remoteURI && gitData.gitUser) {
      this.originalGitData = {
        remoteURI: gitData.remoteURI,
        gitUser: gitData.gitUser,
      };
    }

    this.resetForm(gitData);
  }

  @Output()
  public gitUpstreamSubmit: EventEmitter<GitData> = new EventEmitter();

  @Output()
  public gitDataChanged: EventEmitter<GitData> = new EventEmitter();

  public gitUrlControl = new FormControl('');
  public gitUserControl = new FormControl('');
  public gitTokenControl = new FormControl('');
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl,
  });
  public _isLoading: boolean | undefined;

  ngOnInit(): void {
    if (!this.isCreateMode) {
      this.gitUrlControl.setValidators([Validators.required]);
      this.gitUserControl.setValidators([Validators.required]);
      this.gitTokenControl.setValidators([Validators.required]);
    }
  }

  public setGitUpstream(): void {
    this.gitUpstreamSubmit.emit({
      remoteURI: this.gitUrlControl.value,
      gitUser: this.gitUserControl.value,
      gitToken: this.gitTokenControl.value,
    });
    this.gitTokenControl.markAsUntouched();
    this.gitTokenControl.markAsPristine();
  }

  public onGitUpstreamFormChange(): void {
    this.gitDataChanged.emit({
      remoteURI: this.gitUrlControl.value,
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

  private resetForm(gitData: GitData | undefined): void {
    this.gitUrlControl.setValue(gitData?.remoteURI || '');
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

export interface GitData {
  remoteURI?: string;
  gitUser?: string;
  gitToken?: string;
  gitFormValid?: boolean;
}
