import { ChangeDetectorRef, Component, HostListener, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../../../_services/data.service';
import { combineLatest, forkJoin, Observable, of, Subject, throwError } from 'rxjs';
import { catchError, filter, map, switchMap, take, takeUntil, tap } from 'rxjs/operators';
import { UniformSubscription } from '../../../../_models/uniform-subscription';
import { DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { IWebhookConfigClient, PreviousWebhookConfig } from 'shared/interfaces/webhook-config';
import { ISequencesFilter } from '../../../../../../shared/interfaces/sequencesFilter';
import { IClientSecret } from '../../../../../../shared/interfaces/secret';
import { EventState } from '../../../../../../shared/models/event-state';
import { AppUtils } from '../../../../_utils/app.utils';
import { NotificationsService } from '../../../../_services/notifications.service';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationType } from '../../../../_models/notification';
import { UniformRegistrationInfo } from '../../../../../../shared/interfaces/uniform-registration-info';
import { SecretScopeDefault } from '../../../../../../shared/interfaces/secret-scope';
import { EventTypes } from '../../../../../../shared/interfaces/event-types';
import { Trace } from '../../../../_models/trace';
import { PendingChangesComponent } from '../../../../_guards/pending-changes.guard';
import { DeleteResult, DeleteType, DeletionProgressEvent } from '../../../../_interfaces/delete';
import { EventService } from '../../../../_services/event.service';
import { handleDeletionError } from '../../../../_components/ktb-danger-zone/ktb-danger-zone.utils';

export interface SubscriptionState {
  projectName: string;
  taskNames: string[];
  subscription: UniformSubscription;
  filter: ISequencesFilter;
  integrationId: string;
  webhook?: IWebhookConfigClient;
  webhookSecrets?: IClientSecret[];
}

@Component({
  selector: 'ktb-modify-uniform-subscription',
  templateUrl: './ktb-modify-uniform-subscription.component.html',
  styleUrls: ['./ktb-modify-uniform-subscription.component.scss'],
})
export class KtbModifyUniformSubscriptionComponent implements OnDestroy, PendingChangesComponent {
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  private taskControl = new FormControl('', [Validators.required]);
  public eventPayload: Record<string, unknown> | undefined;
  public taskSuffixControl = new FormControl('', [Validators.required]);
  private isGlobalControl = new FormControl();
  public data$: Observable<SubscriptionState>;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public editMode = false;
  public updating = false;
  public subscriptionForm = new FormGroup({
    taskPrefix: this.taskControl,
    taskSuffix: this.taskSuffixControl,
    isGlobal: this.isGlobalControl,
  });
  private _previousFilter?: PreviousWebhookConfig;
  public isWebhookFormValid = true;
  public webhookFormDirty = false;
  public isWebhookService = false;
  public suffixes: { value: string; displayValue: string }[] = [
    {
      value: '>',
      displayValue: '*',
    },
    {
      value: EventState.TRIGGERED,
      displayValue: EventState.TRIGGERED,
    },
    {
      value: EventState.STARTED,
      displayValue: EventState.STARTED,
    },
    {
      value: EventState.FINISHED,
      displayValue: EventState.FINISHED,
    },
  ];
  public errorMessage?: string;

  private pendingChangesSubject = new Subject<boolean>();
  public dialogLabel = 'Pending Changes dialog';
  public message = 'You have pending changes. Are you sure you want to leave this page?';
  public unsavedDialogState: null | 'unsaved' = null;
  private isFilterDirty = false;

  public deleteType = DeleteType.SUBSCRIPTION;

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private router: Router,
    private notificationsService: NotificationsService,
    private eventService: EventService,
    private _changeDetectorRef: ChangeDetectorRef
  ) {
    const subscription$ = this.route.paramMap.pipe(
      map((paramMap) => ({
        integrationId: paramMap.get('integrationId'),
        subscriptionId: paramMap.get('subscriptionId'),
        projectName: paramMap.get('projectName'),
      })),
      filter(
        (params): params is { integrationId: string; subscriptionId: string | null; projectName: string } =>
          !!(params.integrationId && params.projectName)
      ),
      switchMap((params) => {
        this.editMode = !!params.subscriptionId;
        if (params.subscriptionId) {
          return this.dataService.getUniformSubscription(params.integrationId, params.subscriptionId);
        } else {
          return of(new UniformSubscription(params.projectName));
        }
      }),
      tap((subscription) => {
        if (this.editMode) {
          this._previousFilter = {
            filter: AppUtils.copyObject(subscription.filter),
            type: subscription.event,
          };
        }
        this.taskControl.setValue(subscription.prefix);
        this.taskSuffixControl.setValue(subscription.suffix);
        this.isGlobalControl.setValue(subscription.isGlobal);

        this.updateIsGlobalCheckbox(subscription);
      }),
      take(1)
    );

    const integrationId$ = this.route.paramMap.pipe(
      map((paramMap) => paramMap.get('integrationId')),
      filter((integrationId: string | null): integrationId is string => !!integrationId),
      take(1)
    );

    const integrationInfo$ = integrationId$.pipe(
      switchMap((integrationId) => this.dataService.getUniformRegistrationInfo(integrationId)),
      take(1),
      takeUntil(this.unsubscribe$)
    );

    integrationInfo$.subscribe((info) => {
      if (!info.isControlPlane) {
        this.suffixes = [
          {
            value: EventState.TRIGGERED,
            displayValue: EventState.TRIGGERED,
          },
        ];
      }
      this.isWebhookService = info.isWebhookService;
    });

    const projectName$ = this.route.paramMap.pipe(
      map((paramMap) => paramMap.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName),
      take(1)
    );

    const taskNames$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getTaskNames(projectName)),
      catchError((err: HttpErrorResponse) => {
        this.errorMessage = err.error;
        this.notificationsService.addNotification(NotificationType.ERROR, err.error);
        return throwError(err);
      }),
      take(1)
    );

    const webhook$ = forkJoin({
      subscription: subscription$,
      projectName: projectName$,
      integrationInfo: integrationInfo$,
    }).pipe(
      switchMap(
        (data: {
          subscription: UniformSubscription;
          projectName: string;
          integrationInfo: UniformRegistrationInfo;
        }) => {
          let webhook: Observable<IWebhookConfigClient | undefined>;
          if (data.integrationInfo.isWebhookService && this.editMode && data.subscription.id) {
            const stage: string | undefined = data.subscription.filter?.stages?.[0];
            const services: string | undefined = data.subscription.filter?.services?.[0];
            webhook = this.dataService.getWebhookConfig(data.subscription.id, data.projectName, stage, services);
            this.updateEventPayload(
              data.projectName,
              data.subscription.filter?.stages ?? [],
              data.subscription.filter?.services ?? []
            );
          } else {
            webhook = of(undefined);
          }
          return webhook;
        }
      ),
      take(1)
    );

    const webhookSecrets$ = integrationInfo$.pipe(
      switchMap((info) => {
        if (info.isWebhookService) {
          return this.dataService.getSecretsForScope(SecretScopeDefault.WEBHOOK);
        }
        return of(undefined);
      })
    );

    const filter$ = projectName$.pipe(switchMap((projectName) => this.dataService.getSequenceFilter(projectName)));

    this.data$ = combineLatest([
      projectName$,
      taskNames$,
      subscription$,
      filter$,
      integrationId$,
      webhook$,
      webhookSecrets$,
    ]).pipe(
      map(([projectName, taskNames, subscription, filterData, integrationId, webhook, webhookSecrets]) => {
        return {
          taskNames,
          subscription,
          filter: filterData,
          projectName,
          integrationId,
          webhook,
          webhookSecrets,
        };
      }),
      tap((data) => {
        this.updateDataSource(data.filter.stages, data.filter.services, data.subscription);
      })
    );

    this.eventService.deletionTriggeredEvent.pipe(takeUntil(this.unsubscribe$)).subscribe((data) => {
      if (data.type === DeleteType.SUBSCRIPTION && data.context) {
        this.deleteSubscription(data.context as SubscriptionState);
      }
    });
  }

  public deleteSubscription(data: SubscriptionState): void {
    if (!data.subscription.id) return;

    this.dataService
      .deleteSubscription(data.integrationId, data.subscription.id, this.isWebhookService)
      .pipe(
        map((): DeletionProgressEvent => {
          return { isInProgress: false, result: DeleteResult.SUCCESS };
        }),
        catchError(handleDeletionError('Subscription'))
      )
      .subscribe((progressEvent) => {
        this.eventService.deletionProgressEvent.next(progressEvent);
        if (progressEvent.result === DeleteResult.SUCCESS) {
          this.resetForms();
          this.router.navigate([
            '/',
            'project',
            data.projectName,
            'settings',
            'uniform',
            'integrations',
            data.integrationId,
          ]);
        }
      });
  }

  private updateEventPayload(projectName: string, stages: string[], services: string[]): void {
    if (this.isWebhookService && this.taskControl.value && this.taskSuffixControl.value) {
      this.eventPayload = undefined;
      this.dataService
        .getIntersectedEvent(
          `${EventTypes.PREFIX}${this.taskControl.value}`,
          this.taskSuffixControl.value,
          projectName,
          stages,
          services
        )
        .subscribe((event: Record<string, unknown>) => {
          this.eventPayload = Object.keys(event).length ? event : Trace.defaultTrace;
        });
    }
  }

  public reloadPage(): void {
    window.location.reload();
  }

  private updateDataSource(stages: string[], services: string[], subscription: UniformSubscription): void {
    const availableServices = subscription.filter.services?.filter((service) => services.some((s) => s === service));

    // check if services have been deleted
    if (availableServices && availableServices?.length !== subscription.filter.services?.length) {
      subscription.filter.services = availableServices;
    }
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: stages.map((name) => ({ name })),
        },
        {
          name: 'Service',
          autocomplete: services.map((name) => ({ name })),
        },
      ],
    } as DtFilterFieldDefaultDataSourceAutocomplete;
  }

  public updateSubscription(
    projectName: string,
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: IWebhookConfigClient
  ): void {
    this.updating = true;
    let update;

    if (this.taskControl.value === EventTypes.PREFIX_SHORT && this.taskSuffixControl.value === '>') {
      subscription.event = `${this.taskControl.value}.${this.taskSuffixControl.value}`;
    } else {
      subscription.event = `${EventTypes.PREFIX}${this.taskControl.value}.${this.taskSuffixControl.value}`;
    }
    subscription.setIsGlobal(this.isGlobalControl.value, projectName);

    if (webhookConfig) {
      webhookConfig.type = subscription.event;
      webhookConfig.prevConfiguration = this._previousFilter;
    }

    if (this.editMode) {
      update = this.dataService.updateUniformSubscription(integrationId, subscription, webhookConfig);
    } else {
      update = this.dataService.createUniformSubscription(integrationId, subscription, webhookConfig);
    }

    update.subscribe(
      () => {
        this.updating = false;
        this.notificationsService.addNotification(NotificationType.SUCCESS, 'Subscription successfully created!');
        this.resetForms();
        this.router.navigate(['/', 'project', projectName, 'settings', 'uniform', 'integrations', integrationId]);
      },
      () => {
        this.notificationsService.addNotification(NotificationType.ERROR, 'The subscription could not be updated');
        this.updating = false;
      }
    );
  }

  public isFormValid(subscription: UniformSubscription): boolean {
    return (
      this.subscriptionForm.valid &&
      (!!subscription.filter.stages?.length || !subscription.filter.services?.length) &&
      this.isWebhookFormValid &&
      !this.updating
    );
  }

  public selectedTaskChanged(projectName: string, subscription: UniformSubscription): void {
    this.updateEventPayload(projectName, subscription.filter.stages ?? [], subscription.filter.services ?? []);
  }

  public subscriptionFilterChanged(subscription: UniformSubscription, projectName: string): void {
    this.isFilterDirty = !!subscription.filter.stages?.length || !!subscription.filter.services?.length;
    this.updateIsGlobalCheckbox(subscription);
    this.updateEventPayload(projectName, subscription.filter.stages ?? [], subscription.filter.services ?? []);
  }

  public updateIsGlobalCheckbox(subscription: UniformSubscription): void {
    if (subscription.hasFilter()) {
      this.isGlobalControl.disable({ onlySelf: true, emitEvent: false });
      this.isGlobalControl.setValue(false);
    } else {
      this.isGlobalControl.enable({ onlySelf: true, emitEvent: false });
    }
  }

  public getSelectedTask(): string | undefined {
    let type: string | undefined;
    if (this.taskControl.value) {
      if (this.taskSuffixControl.value && this.taskSuffixControl.value !== this.suffixes[0].value) {
        type = `${EventTypes.PREFIX}${this.taskControl.value}.${this.taskSuffixControl.value}`;
      } else {
        type = `${EventTypes.PREFIX}${this.taskControl.value}.triggered`;
      }
    }

    return type;
  }

  public webhookFormValidityChanged(isValid: boolean): void {
    this.isWebhookFormValid = isValid;
    this._changeDetectorRef.detectChanges();
  }

  // @HostListener allows us to also guard against browser refresh, close, etc.
  @HostListener('window:beforeunload', ['$event'])
  public canDeactivate(_$event?: BeforeUnloadEvent): Observable<boolean> {
    if (this.subscriptionForm.dirty || this.webhookFormDirty || this.isFilterDirty) {
      this.showNotification();
      return this.pendingChangesSubject.asObservable();
    }
    return of(true);
  }

  public resetForms(): void {
    this.subscriptionForm.reset();
    this.webhookFormDirty = false;
    this.isFilterDirty = false;
  }

  public reject(): void {
    this.pendingChangesSubject.next(false);
    this.hideNotification();
  }

  public reset(): void {
    this.pendingChangesSubject.next(true);
    this.hideNotification();
  }

  public showNotification(): void {
    this.unsavedDialogState = 'unsaved';

    const dialog = document.querySelector(`div[aria-label="${this.dialogLabel}"]`);

    if (!dialog) return;

    dialog.classList.add('shake');
    setTimeout(() => dialog.classList.remove('shake'), 500);
  }

  public hideNotification(): void {
    this.unsavedDialogState = null;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
