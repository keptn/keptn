import {ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Project} from "../../_models/project";
import {filter, takeUntil} from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {Metadata} from '../../_models/metadata';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-project-tile',
  templateUrl: './ktb-project-tile.component.html',
  styleUrls: ['./ktb-project-tile.component.scss']
})
export class KtbProjectTileComponent implements OnInit, OnDestroy {

  public _project: Project;
  public supportedShipyardVersion: string;
  private unsubscribe$ = new Subject();

  @Input()
  get project(): Project {
    return this._project;
  }
  set project(value: Project) {
    if (this._project !== value) {
      this._project = value;
      this._changeDetectorRef.markForCheck();
    }
  }

  constructor(private _changeDetectorRef: ChangeDetectorRef, private dataService: DataService) {
    this.dataService.keptnInfo
      .pipe(filter(keptnInfo => !!keptnInfo))
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(keptnInfo => {
        this.supportedShipyardVersion = (keptnInfo.metadata as Metadata)?.shipyardversion;
      });
  }

  ngOnInit() {
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }
}
