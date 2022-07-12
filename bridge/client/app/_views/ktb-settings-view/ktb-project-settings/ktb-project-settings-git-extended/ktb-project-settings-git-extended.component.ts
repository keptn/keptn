import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { isGitHttps, isGitSsh } from '../../../../_utils/git-upstream.utils';
import {
  IGitDataExtended,
  IGitHttpsConfiguration,
  IGitSshConfiguration,
} from '../../../../../../shared/interfaces/project';

export enum GitFormType {
  SSH,
  HTTPS,
  NO_UPSTREAM,
}

export type IGitDataExtendedWithNoUpstream = IGitDataExtended | null; // null => no-upstream

@Component({
  selector: 'ktb-project-settings-git-extended',
  templateUrl: './ktb-project-settings-git-extended.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitExtendedComponent implements OnInit {
  public selectedForm: GitFormType = GitFormType.NO_UPSTREAM;
  public gitInputDataHttps?: IGitHttpsConfiguration;
  public gitInputDataSsh?: IGitSshConfiguration;
  public FormType = GitFormType;
  public gitDataHttps?: IGitHttpsConfiguration;
  public gitDataSsh?: IGitSshConfiguration;

  @Input()
  public isCreateMode = false;

  @Input()
  public gitUpstreamRequired = true;

  @Input()
  public gitInputData: IGitDataExtended | undefined;

  @Output()
  public gitDataChange = new EventEmitter<IGitDataExtendedWithNoUpstream | undefined>();

  public get gitData(): IGitDataExtendedWithNoUpstream | undefined {
    switch (this.selectedForm) {
      case GitFormType.HTTPS:
        return this.gitDataHttps;
      case GitFormType.SSH:
        return this.gitDataSsh;
      case GitFormType.NO_UPSTREAM:
        return null;
      default:
        return undefined;
    }
  }

  public ngOnInit(): void {
    if (!this.gitUpstreamRequired || !this.gitInputData) {
      this.selectedForm = this.gitUpstreamRequired ? GitFormType.HTTPS : GitFormType.NO_UPSTREAM;

      if (this.selectedForm === GitFormType.NO_UPSTREAM) {
        this.dataChanged(this.selectedForm, this.gitData);
      }

      return;
    }

    if (isGitHttps(this.gitInputData)) {
      this.gitInputDataHttps = this.gitInputData;
      this.selectedForm = GitFormType.HTTPS;
      return;
    }

    if (isGitSsh(this.gitInputData)) {
      this.gitInputDataSsh = this.gitInputData;
      this.selectedForm = GitFormType.SSH;
    }
  }

  public setSelectedForm($event: DtRadioChange<GitFormType>): void {
    this.selectedForm = $event.value ?? GitFormType.HTTPS;
    this.dataChanged(this.selectedForm, this.gitData);
  }

  public dataChanged(type: GitFormType, data?: IGitDataExtendedWithNoUpstream): void {
    // the data should be split into two in order to update the parent form correctly if the selected form is switched.
    // On switch the child component does not emit new data and therefore the selected data is not updated
    if (data) {
      switch (type) {
        case GitFormType.HTTPS:
          this.gitDataHttps = data as IGitHttpsConfiguration;
          break;
        case GitFormType.SSH:
          this.gitDataSsh = data as IGitSshConfiguration;
          break;
      }
    } else {
      if (this.selectedForm === GitFormType.HTTPS) {
        this.gitDataHttps = undefined;
      } else if (this.selectedForm === GitFormType.SSH) {
        this.gitDataSsh = undefined;
      }
    }
    this.gitDataChange.emit(data);
  }
}
