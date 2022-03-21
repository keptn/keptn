import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { IGitDataExtended, IGitHttps, IGitSsh } from '../../_interfaces/git-upstream';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { filter, map } from 'rxjs/operators';
import { isGitHTTPS } from '../../_utils/git-upstream.utils';

export enum GitFormType {
  SSH,
  HTTPS,
}

@Component({
  selector: 'ktb-project-settings-git-extended',
  templateUrl: './ktb-project-settings-git-extended.component.html',
  styleUrls: ['./ktb-project-settings-git-extended.component.scss'],
})
export class KtbProjectSettingsGitExtendedComponent {
  // TODO: on https/ssh change, should the data be discarded or not? If not, the invalid data needs to be temporarily saved
  //  solution: on change of FormControl in component change the gitInputData. Should be a reference. e.g. https component: get and set for proxy in order to adjust the parent gitDataInput
  private projectName?: string;
  private _gitInputData?: IGitDataExtended;
  public FormType = GitFormType;
  public selectedForm: GitFormType = GitFormType.HTTPS;
  public gitData?: IGitDataExtended;

  @Input()
  public isLoading = false;

  @Input()
  public isCreateMode = false;

  @Input()
  public isGitUpstreamInProgress = false;

  @Input()
  public set gitInputData(gitData: IGitDataExtended | undefined) {
    this._gitInputData = gitData;
    this.selectedForm = !gitData || isGitHTTPS(gitData) ? GitFormType.HTTPS : GitFormType.SSH;
  }
  public get gitInputData(): IGitDataExtended | undefined {
    return this._gitInputData;
  }

  @Output()
  public gitDataChange = new EventEmitter<IGitDataExtended | undefined>();

  public get gitInputDataHTTPS(): IGitHttps | undefined {
    return this.gitInputData && isGitHTTPS(this.gitInputData) ? this.gitInputData : undefined;
  }

  public get gitInputDataSSH(): IGitSsh | undefined {
    return this.gitInputData && !isGitHTTPS(this.gitInputData) ? this.gitInputData : undefined;
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

  public setSelectedForm($event: DtRadioChange<GitFormType>): void {
    this.selectedForm = $event.value ?? GitFormType.HTTPS;
    this.dataChanged(); // on change reset form because the data is invalid then
  }

  public updateUpstream(): void {
    if (this.gitData && this.projectName) {
      this.dataService.updateGitUpstream(this.projectName, this.gitData).subscribe();
    }
  }

  public dataChanged(data?: IGitDataExtended): void {
    this.gitData = data;
    this.gitDataChange.emit(data);
  }
}
