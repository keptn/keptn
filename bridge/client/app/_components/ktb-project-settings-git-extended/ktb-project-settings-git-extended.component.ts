import { Component, EventEmitter, Input, Output } from '@angular/core';
import { DtRadioChange } from '@dynatrace/barista-components/radio';
import { IGitDataExtended, IGitHttps, IGitSsh } from '../../_interfaces/git-upstream';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, ParamMap } from '@angular/router';
import { filter, map } from 'rxjs/operators';
import { isGitHTTPS } from '../../_utils/git-upstream.utils';

enum FormType {
  SSH,
  HTTPS,
}

@Component({
  selector: 'ktb-project-settings-git-extended',
  templateUrl: './ktb-project-settings-git-extended.component.html',
  styleUrls: ['./ktb-project-settings-git-extended.component.scss'],
})
export class KtbProjectSettingsGitExtendedComponent {
  private projectName = '';
  private _gitInputData?: IGitDataExtended;
  public FormType = FormType;
  public selectedForm: FormType = FormType.HTTPS;
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
    if (gitData) {
      this.selectedForm = isGitHTTPS(gitData) ? FormType.HTTPS : FormType.SSH;
    }
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
    return this.gitInputData && isGitHTTPS(this.gitInputData) ? undefined : this.gitInputData;
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

  public setSelectedForm($event: DtRadioChange<FormType>): void {
    this.selectedForm = $event.value ?? FormType.HTTPS;
    this.dataChanged(); // on change reset form
  }

  public updateUpstream(): void {
    if (this.gitData) {
      this.dataService.updateGitUpstream(this.projectName, this.gitData);
    }
  }

  public dataChanged(data?: IGitDataExtended): void {
    this.gitData = data;
    this.gitDataChange.emit(data);
  }
}
