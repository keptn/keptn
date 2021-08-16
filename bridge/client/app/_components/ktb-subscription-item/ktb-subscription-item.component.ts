import {ChangeDetectionStrategy, ChangeDetectorRef, Component, Input, OnDestroy, OnInit} from '@angular/core';
import {UniformSubscription} from '../../_models/uniform-subscription';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import {DataService} from '../../_services/data.service';
import {ActivatedRoute} from '@angular/router';
import {Project} from '../../_models/project';
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
      });
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
