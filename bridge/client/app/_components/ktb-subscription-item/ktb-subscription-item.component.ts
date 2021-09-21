import { ChangeDetectionStrategy, ChangeDetectorRef, Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, Router } from '@angular/router';
import { Project } from '../../_models/project';
import { Subject } from 'rxjs';
import { DeleteDialogState } from '../_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';

@Component({
  selector: 'ktb-subscription-item[subscription][integrationId][isWebhookService]',
  templateUrl: './ktb-subscription-item.component.html',
  styleUrls: ['./ktb-subscription-item.component.scss'],
  changeDetection: ChangeDetectionStrategy.OnPush,
})
export class KtbSubscriptionItemComponent implements OnInit, OnDestroy {
  private _subscription?: UniformSubscription;
  public project?: Project;
  private readonly unsubscribe$ = new Subject<void>();
  private currentSubscription?: UniformSubscription;
  public deleteState: DeleteDialogState = null;

  @Output() subscriptionDeleted: EventEmitter<UniformSubscription> = new EventEmitter<UniformSubscription>();
  @Input() name?: string;
  @Input() integrationId?: string;
  @Input() isWebhookService = false;

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

  constructor(private changeDetectorRef: ChangeDetectorRef, private dataService: DataService, private route: ActivatedRoute, private router: Router) {
  }

  ngOnInit(): void {
    this.route.paramMap
      .pipe(
        map(params => params.get('projectName')),
        filter((projectName: string | null): projectName is string => !!projectName),
        switchMap(projectName => this.dataService.getProject(projectName)),
        filter((project: Project | undefined): project is Project => !!project),
        takeUntil(this.unsubscribe$),
      ).subscribe(project => {
      this.project = project;
    });
  }

  public editSubscription(subscription: UniformSubscription): void {
    this.router.navigate(['/', 'project', this.project?.projectName, 'uniform', 'services', this.integrationId, 'subscriptions', subscription.id, 'edit']);
  }

  public triggerDeleteSubscription(subscription: UniformSubscription): void {
    this.currentSubscription = subscription;
    this.deleteState = 'confirm';
  }

  public deleteSubscription(): void {
    if (this.integrationId && this.subscription?.id) {
      this.dataService.deleteSubscription(this.integrationId, this.subscription.id, this.isWebhookService).subscribe(() => {
        this.deleteState = 'success';
        this.subscriptionDeleted.emit(this.subscription);
      });
    }
  }


  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }

}
