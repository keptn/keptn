import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";

@Component({
  selector: 'ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss']
})
export class KtbProjectSettingsGitComponent implements OnInit {

  @Input()
  public isGitUpstreamInProgress: boolean;

  @Input()
  public isCreateMode: boolean;

  @Input()
  set gitData(gitData: GitData) {
    this.gitUrlControl.setValue(gitData.remoteURI || '');
    this.gitUserControl.setValue(gitData.gitUser || '');

    if (gitData.remoteURI && gitData.gitUser) {
      this.gitTokenControl.setValue('***********************');
    } else {
      this.gitTokenControl.setValue('');
    }
  }

  @Output()
  public onGitUpstreamSubmit: EventEmitter<GitData> = new EventEmitter();

  @Output()
  private onGitDataChanged: EventEmitter<GitData> = new EventEmitter();

  public gitUrlControl = new FormControl('');
  public gitUserControl = new FormControl('');
  public gitTokenControl = new FormControl('');
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl
  });

  constructor() { }

  ngOnInit(): void {
    if(!this.isCreateMode) {
      this.gitUrlControl.setValidators([Validators.required]);
      this.gitUserControl.setValidators([Validators.required]);
      this.gitTokenControl.setValidators([Validators.required]);
    }
  }

  public setGitUpstream() {
    this.onGitUpstreamSubmit.emit({remoteURI: this.gitUrlControl.value, gitUser: this.gitUserControl.value, gitToken: this.gitTokenControl.value});
  }

  public onGitUpstreamFormChange() {
    this.onGitDataChanged.emit({remoteURI: this.gitUrlControl.value, gitUser: this.gitUserControl.value, gitToken: this.gitTokenControl.value});
  }

}

export interface GitData {
  remoteURI?: string;
  gitUser?: string;
  gitToken?: string;
}
