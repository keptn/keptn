import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { IGitDataExtended, IGitHttps, IGitSsh } from '../../_interfaces/git-upstream';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { filter, map } from 'rxjs/operators';
import { isGitHTTPS, isGitSSH } from '../../_utils/git-upstream.utils';

export enum GitFormType {
  SSH,
  HTTPS,
  NO_UPSTREAM,
}

@Component({
  selector: 'ktb-project-settings-git-extended',
  templateUrl: './ktb-project-settings-git-extended.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitExtendedComponent implements OnInit {
  private projectName?: string;
  public selectedForm: GitFormType = GitFormType.NO_UPSTREAM;
  public gitInputDataHttps?: IGitHttps;
  public gitInputDataSsh?: IGitSsh;
  public FormType = GitFormType;
  public gitDataHttps?: IGitDataExtended;
  public gitDataSsh?: IGitDataExtended;
  public upstreamConfigured = false;

  @Input()
  public isCreateMode = false;

  @Input()
  public isGitUpstreamInProgress = false;

  @Input()
  public required = true;

  @Input()
  public gitInputData: IGitDataExtended | undefined;

  @Output()
  public gitDataChange = new EventEmitter<IGitDataExtended | undefined>();

  @Output()
  public resetTouched = new EventEmitter<void>();

  public get gitData(): IGitDataExtended | undefined {
    switch (this.selectedForm) {
      case GitFormType.HTTPS:
        return this.gitDataHttps;
      case GitFormType.SSH:
        return this.gitDataSsh;
      case GitFormType.NO_UPSTREAM:
        return { noupstream: '' };
      default:
        return undefined;
    }
  }

  public ngOnInit(): void {
    if (!this.gitInputData || this.isRemoteUrlEmpty(this.gitInputData)) {
      this.selectedForm = this.required ? GitFormType.HTTPS : GitFormType.NO_UPSTREAM;

      if (this.selectedForm === GitFormType.NO_UPSTREAM) {
        this.dataChanged(this.selectedForm, this.gitData);
      }

      return;
    }

    if (isGitHTTPS(this.gitInputData)) {
      this.gitInputDataHttps = this.gitInputData;
      this.selectedForm = GitFormType.HTTPS;
      this.upstreamConfigured = true;
      return;
    }

    if (isGitSSH(this.gitInputData)) {
      this.gitInputDataSsh = this.gitInputData;
      this.selectedForm = GitFormType.SSH;
      this.upstreamConfigured = true;
    }
  }

  constructor(private readonly dataService: DataService, readonly routes: ActivatedRoute) {
    this.routes.paramMap
      .pipe(
        map((params: ParamMap) => params.get('projectName')),
        filter((projectName: string | null): projectName is string => !!projectName)
      )
      .subscribe((projectName: string) => {
        this.projectName = projectName;
      });
  }

  private isRemoteUrlEmpty(gitInputData: IGitDataExtended): boolean {
    return (
      (isGitHTTPS(gitInputData) && !gitInputData.https.gitRemoteURL) ||
      (isGitSSH(gitInputData) && !gitInputData.ssh.gitRemoteURL)
    );
  }

  public setSelectedForm($event: DtRadioChange<GitFormType>): void {
    this.selectedForm = $event.value ?? GitFormType.HTTPS;
    this.dataChanged(this.selectedForm, this.gitData);
  }

  public updateUpstream(): void {
    if (this.gitData && this.projectName) {
      this.isGitUpstreamInProgress = true;
      this.dataService.updateGitUpstream(this.projectName, this.gitData).subscribe(
        () => {
          this.isGitUpstreamInProgress = false;
          this.resetTouched.emit();
        },
        () => {
          this.isGitUpstreamInProgress = false;
        }
      );
    }
  }

  public dataChanged(type: GitFormType, data?: IGitDataExtended): void {
    // the data should be split into two in order to update the parent form correctly if the selected form is switched.
    // On switch the child component does not emit new data and therefore the selected data is not updated
    if (data) {
      switch (type) {
        case GitFormType.HTTPS:
          this.gitDataHttps = data;
          break;
        case GitFormType.SSH:
          this.gitDataSsh = data;
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
