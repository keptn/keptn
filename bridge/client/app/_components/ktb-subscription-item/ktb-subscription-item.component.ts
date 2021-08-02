import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {UniformSubscription} from '../../_models/uniformSubscription';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../_models/project';
import { DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import {Subject} from 'rxjs';

@Component({
  selector: 'ktb-subscription-item[subscription]',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class KtbSubscriptionItemComponent implements OnInit, OnDestroy {
  private _subscription?: UniformSubscription;
  public project?: Project;
  public tasks: string[] = [];
  public _dataSource = new DtFilterFieldDefaultDataSource();
  private readonly unsubscribe$ = new Subject<void>();

  @Input() name?: string;

  @Input()
  get subscription(): UniformSubscription | undefined {
    return this._subscription;
  }
  set subscription(subscription: UniformSubscription | undefined) {
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
        filter((project: Project | undefined): project is Project => !!project),
        takeUntil(this.unsubscribe$)
      ).subscribe(project => {
        this.project = project;
        this.updateDataSource(project);
      });
  }

  public updateDataSource(project: Project) {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Service',
          autocomplete: project.getServices().map(service => {
            return {
              name: service.serviceName
            };
          })
        }
      ]
    };
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
