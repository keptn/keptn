import {ChangeDetectorRef, Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {Project} from '../../_models/project';
import {Stage} from '../../_models/stage';
import {DataService} from '../../_services/data.service';

@Component({
  selector: 'ktb-stage-overview',
  templateUrl: './ktb-stage-overview.component.html',
  styleUrls: ['./ktb-stage-overview.component.scss']
})
export class KtbStageOverviewComponent implements OnInit {
  public _project: Project;
  public selectedStage: Stage = null;

  @Output() selectedStageChange: EventEmitter<any> = new EventEmitter();

  @Input()
  get project() {
    return this._project;
  }

  set project(project: Project) {
    if (this._project !== project) {
      this._project = project;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private dataService: DataService, private _changeDetectorRef: ChangeDetectorRef) {
  }

  ngOnInit(): void {
  }

  trackStage(index: number, stage: Stage): string {
    return stage.stageName;
  }

  selectStage($event, stage: Stage, filterType?: string) {
    this.selectedStage = stage;
    $event.stopPropagation();
    this.selectedStageChange.emit({stage, filterType});
  }

}
