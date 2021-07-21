import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from '../../_models/subscription';
import {map, switchMap, takeUntil} from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../_models/project';
import { DtFilterFieldChangeEvent, DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import {Subject} from 'rxjs';
import {ProjectMock} from '../../_models/project.mock';
import { DtAutoComplete } from '../../_models/dt-filter';

@Component({
  selector: 'ktb-subscription-item[subscription]',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionItemComponent implements OnInit, OnDestroy {
  private _subscription!: Subscription;
  public tasks: string[] = [];
  public _dataSource = new DtFilterFieldDefaultDataSource();
  private readonly unsubscribe$ = new Subject<void>();

  @Input()
  get subscription(): Subscription {
    return this._subscription;
  }
  set subscription(subscription: Subscription) {
    if (this._subscription !== subscription) {
      this._subscription = subscription;
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.dataService.taskNamesTriggered
      .pipe(
        takeUntil(this.unsubscribe$),
      ).subscribe( tasks => {
        this.tasks = ['all', ...tasks];
      });

    this.route.params
      .pipe(
        map(params => params.projectName),
        switchMap(projectName => this.dataService.getProject(projectName)),
        takeUntil(this.unsubscribe$)
      ).subscribe(project => {
        if (project) {
          this.updateDataSource(project);
        }
      });
  }

  public updateDataSource(project: Project) {
    project = ProjectMock;
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Stages',
          autocomplete: project.stages.map(stage => {
            return {
              name: stage.stageName
            };
          })
        },
        {
          name: 'Services',
          autocomplete: project.getServices().map(service => {
            return {
              name: service.serviceName
            };
          })
        }
      ]
    };
  }

  // tslint:disable-next-line:no-any
  filterChanged(subscription: Subscription, event: DtFilterFieldChangeEvent<any>) { // can't set another type because of "is not assignable to..."
    event = event as DtFilterFieldChangeEvent<DtAutoComplete>;
    const result = event.filters.reduce((filters: {Stages: string[], Services: string[]}, filter) => {
      filters[filter[0].name as 'Stages' | 'Services'].push(filter[1].name);
      return filters;
    }, {Stages: [], Services: []});
    subscription.services = result.Services;
    subscription.stages = result.Stages;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
