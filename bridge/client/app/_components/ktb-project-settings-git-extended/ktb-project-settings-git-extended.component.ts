import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { IGitDataExtended, IGitHttps, IGitSsh } from '../../_interfaces/git-upstream';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { filter, map } from 'rxjs/operators';
import { isGitHTTPS } from '../../_utils/git-upstream.utils';
import { NotificationType } from '../../_models/notification';
import { NotificationsService } from '../../_services/notifications.service';

export enum GitFormType {
  SSH,
  HTTPS,
}

@Component({
  selector: 'ktb-project-settings-git-extended',
  templateUrl: './ktb-project-settings-git-extended.component.html',
  styleUrls: [],
})
export class KtbProjectSettingsGitExtendedComponent {
  private projectName?: string;
  public gitInputDataHttps?: IGitHttps;
  public gitInputDataSsh?: IGitSsh;
  public FormType = GitFormType;
  public selectedForm: GitFormType = GitFormType.HTTPS;
  public gitDataHttps?: IGitHttps;
  public gitDataSsh?: IGitSsh;

  @Input()
  public isLoading = false;

  @Input()
  public isCreateMode = false;

  @Input()
  public isGitUpstreamInProgress = false;

  @Input()
  public set gitInputData(gitData: IGitDataExtended | undefined) {
    if (gitData) {
      if (isGitHTTPS(gitData)) {
        this.gitInputDataHttps = gitData;
        this.selectedForm = GitFormType.HTTPS;
      } else {
        this.gitInputDataSsh = gitData;
        this.selectedForm = GitFormType.SSH;
      }
    } else {
      this.selectedForm = GitFormType.HTTPS;
    }
  }

  @Output()
  public gitDataChange = new EventEmitter<IGitDataExtended | undefined>();

  @Output()
  public resetTouched = new EventEmitter<void>();

  public get gitData(): IGitDataExtended | undefined {
    return this.selectedForm === GitFormType.HTTPS ? this.gitDataHttps : this.gitDataSsh;
  }

  constructor(
    private readonly dataService: DataService,
    readonly routes: ActivatedRoute,
    private notificationsService: NotificationsService
  ) {
    this.routes.paramMap
      .pipe(
        map((params: ParamMap) => params.get('projectName')),
        filter((projectName: string | null): projectName is string => !!projectName)
      )
      .subscribe((projectName: string) => {
        this.projectName = projectName;
      });
  }

  public setSelectedForm($event: DtRadioChange<GitFormType>): void {
    this.selectedForm = $event.value ?? GitFormType.HTTPS;
    this.dataChanged(this.gitData);
  }

  public updateUpstream(): void {
    if (this.gitData && this.projectName) {
      this.isGitUpstreamInProgress = true;
      this.dataService.updateGitUpstream(this.projectName, this.gitData).subscribe(
        () => {
          this.isGitUpstreamInProgress = false;
          this.notificationsService.addNotification(
            NotificationType.SUCCESS,
            'The Git upstream was changed successfully.'
          );
          this.resetTouched.emit();
        },
        () => {
          this.isGitUpstreamInProgress = false;
        }
      );
    }
  }

  public dataChanged(data?: IGitDataExtended): void {
    // the data should be split into two in order to update the parent form correctly if the selected form is switched.
    // On switch the child component does not emit new data and therefore the selected data is not updated
    if (data) {
      if (isGitHTTPS(data)) {
        this.gitDataHttps = data;
      } else {
        this.gitDataSsh = data;
      }
    } else {
      if (this.selectedForm === GitFormType.HTTPS) {
        this.gitDataHttps = undefined;
      } else {
        this.gitDataSsh = undefined;
      }
    }
    this.gitDataChange.emit(data);
  }
}
