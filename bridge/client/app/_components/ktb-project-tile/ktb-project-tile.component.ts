import { Component, Input, OnDestroy } from '@angular/core';
import { Project } from '../../_models/project';
import { filter, takeUntil } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { Metadata } from '../../_models/metadata';
import { Subject } from 'rxjs';
import { KeptnInfo } from '../../_models/keptn-info';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
})
export class KtbProjectTileComponent implements OnDestroy {

  public _project?: Project;
  public supportedShipyardVersion?: string;
  private unsubscribe$ = new Subject();

  @Input()
  get project(): Project | undefined {
    return this._project;
  }

  set project(value: Project | undefined) {
    if (this._project !== value) {
      this._project = value;
    }
  }

  constructor(private dataService: DataService) {
    this.dataService.keptnInfo
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((keptnInfo: KeptnInfo | undefined): keptnInfo is KeptnInfo => !!keptnInfo),
      ).subscribe(keptnInfo => {
      this.supportedShipyardVersion = (keptnInfo.metadata as Metadata)?.shipyardversion;
    });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
