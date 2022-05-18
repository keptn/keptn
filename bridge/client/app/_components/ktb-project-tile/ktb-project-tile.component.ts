import { Component, Input, OnDestroy } from '@angular/core';
import { Project } from '../../_models/project';
import { filter, takeUntil } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { IMetadata } from '../../_interfaces/metadata';
import { Subject } from 'rxjs';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss'],
})
export class KtbProjectTileComponent implements OnDestroy {
  public _project?: Project;
  public supportedShipyardVersion?: string | null;
  private unsubscribe$ = new Subject<void>();

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
    this.dataService.keptnMetadata
      .pipe(
        takeUntil(this.unsubscribe$),
        filter((metadata): metadata is IMetadata | null => metadata !== undefined)
      )
      .subscribe((metadata) => {
        this.supportedShipyardVersion = metadata === null ? null : metadata.shipyardversion;
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
