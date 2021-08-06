import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup, Validators} from '@angular/forms';

@Component({
  selector: 'ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss']
})
export class KtbProjectSettingsGitComponent implements OnInit {

  @Input()
  public isGitUpstreamInProgress = false;

  @Input()
  public isCreateMode = false;

  @Input()
  set gitData(gitData: GitData) {
    this.gitUrlControl.setValue(gitData.remoteURI || '');
    this.gitUserControl.setValue(gitData.gitUser || '');
    this.gitTokenControl.setValue('');
    this.gitTokenControl.markAsUntouched();
    this.gitTokenControl.markAsPristine();
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
    gitToken: this.gitTokenControl
  });

  ngOnInit(): void {
    if (!this.isCreateMode) {
      this.gitUrlControl.setValidators([Validators.required]);
      this.gitUserControl.setValidators([Validators.required]);
      this.gitTokenControl.setValidators([Validators.required]);
    }
  }

  public setGitUpstream() {
    this.gitUpstreamSubmit.emit({remoteURI: this.gitUrlControl.value, gitUser: this.gitUserControl.value, gitToken: this.gitTokenControl.value});
  }

  public onGitUpstreamFormChange() {
    this.gitDataChanged.emit({remoteURI: this.gitUrlControl.value, gitUser: this.gitUserControl.value, gitToken: this.gitTokenControl.value});
  }

}

export interface GitData {
  remoteURI?: string;
  gitUser?: string;
  gitToken?: string;
}
