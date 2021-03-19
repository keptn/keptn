import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {Subscription} from '../../_models/subscription';
import {map, takeUntil} from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../_models/project';
import {DtFilterFieldDefaultDataSource} from '@dynatrace/barista-components/filter-field';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-subscription-item',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionItemComponent implements OnInit, OnDestroy {
  private _subscription: Subscription;
  public tasks: string[] = [];
  private _isReadOnly: boolean;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  private _data: any;
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

  @Input()
  get isReadOnly(): boolean {
    return this._isReadOnly;
  }
  set isReadOnly(isReadOnly: boolean) {
    if (this._isReadOnly !== isReadOnly) {
      this._isReadOnly = isReadOnly;
      this._dataSource = new DtFilterFieldDefaultDataSource(this._data);
      this.changeDetectorRef.markForCheck();
    }
  }

  constructor(private changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute) { }

  ngOnInit(): void {
    this.route.params
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(params => {
        this.dataService.taskNames.pipe(
          takeUntil(this.unsubscribe$),
          map(tasks => tasks.map(task => task + '.triggered'))
        ).subscribe( tasks => {
            this.tasks = tasks;
          }
        );
        this.dataService.getProject(params.projectName)
          .pipe(takeUntil(this.unsubscribe$))
          .subscribe(project => {
            this.updateDataSource(project);
          });
      });
  }

  public updateDataSource(project: Project) {
    this._data = {
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
          autocomplete: project.services.map(service => {
            return {
              name: service.serviceName
            };
          })
        }
      ]
    };
    this._dataSource = new DtFilterFieldDefaultDataSource(this._data);
  }

  filterChanged(subscription: Subscription, event) {
    const result = event.filters.reduce((filters, filter) => {
      filters[filter[0].name].push(filter[1].name);
      return filters;
    }, {Stages: [], Services: []});
    subscription.services = result.Services;
    subscription.stages = result.Stages;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
  }

}
