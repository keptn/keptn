import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {FormControl, FormGroup, Validators} from "@angular/forms";

@Component({
  selector: 'app-ktb-project-settings-git',
  templateUrl: './ktb-project-settings-git.component.html',
  styleUrls: ['./ktb-project-settings-git.component.scss']
})
export class KtbProjectSettingsGitComponent implements OnInit {

  @Input()
  public isGitUpstreamInProgress: boolean;

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
  public gitUpstreamSet: EventEmitter<GitData> = new EventEmitter();

  public gitUrlControl = new FormControl('', [Validators.required]);
  public gitUserControl = new FormControl('', [Validators.required]);
  public gitTokenControl = new FormControl('', [Validators.required]);
  public gitUpstreamForm = new FormGroup({
    gitUrl: this.gitUrlControl,
    gitUser: this.gitUserControl,
    gitToken: this.gitTokenControl
  });

  constructor() { }

  ngOnInit(): void {
  }

  setGitUpstream() {
    this.gitUpstreamSet.emit({remoteURI: this.gitUrlControl.value, gitUser: this.gitUrlControl.value, gitToken: this.gitTokenControl.value});
  }

}

export interface GitData {
  remoteURI?: string;
  gitUser?: string;
  gitToken?: string;
}
