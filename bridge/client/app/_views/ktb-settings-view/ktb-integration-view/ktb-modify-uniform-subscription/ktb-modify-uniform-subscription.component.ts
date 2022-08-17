import { ChangeDetectorRef, Component, HostListener, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../../../_services/data.service';
import { combineLatest, EMPTY, mergeMap, Observable, of, shareReplay, Subject, throwError } from 'rxjs';
import { catchError, finalize, map, switchMap, takeUntil, tap } from 'rxjs/operators';
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
import { SecretScopeDefault } from '../../../../../../shared/interfaces/secret-scope';
import { EventTypes } from '../../../../../../shared/interfaces/event-types';
import { Trace } from '../../../../_models/trace';
import { PendingChangesComponent } from '../../../../_guards/pending-changes.guard';
import { DeleteResult, DeleteType, DeletionProgressEvent } from '../../../../_interfaces/delete';
import { EventService } from '../../../../_services/event.service';
import { handleDeletionError } from '../../../../_components/ktb-danger-zone/ktb-danger-zone.utils';

export interface SubscriptionState {
  taskNames: string[];
  subscription: UniformSubscription;
  filter: ISequencesFilter;
  isWebhookService: boolean;
  webhook?: IWebhookConfigClient;
  webhookSecrets?: IClientSecret[];
}

export type Params = { projectName: string; integrationId: string; subscriptionId?: string; editMode: boolean };
type Suffix = { value: string; displayValue: string };
type NotificationResult = { type: NotificationType; message: string };

@Component({
  selector: 'ktb-modify-uniform-subscription',
  templateUrl: './ktb-modify-uniform-subscription.component.html',
  styleUrls: ['./ktb-modify-uniform-subscription.component.scss'],
})
export class KtbModifyUniformSubscriptionComponent implements OnDestroy, PendingChangesComponent {
  private readonly unsubscribe$: Subject<void> = new Subject<void>();

  private taskControl = new FormControl('', [Validators.required]);
  public taskSuffixControl = new FormControl('', [Validators.required]);
  private isGlobalControl = new FormControl();
  public subscriptionForm = new FormGroup({
    taskPrefix: this.taskControl,
    taskSuffix: this.taskSuffixControl,
    isGlobal: this.isGlobalControl,
  });

  public data$: Observable<SubscriptionState>;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public updating = false;
  public eventPayload: Record<string, unknown> | undefined;
  private _previousFilter?: PreviousWebhookConfig;
  public isWebhookFormValid = true;
  public webhookFormDirty = false;
  public suffixes: Suffix[] = [
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

  public params$ = this.route.paramMap.pipe(
    mergeMap((paramMap) => {
      const projectName = paramMap.get('projectName');
      const integrationId = paramMap.get('integrationId');
      const subscriptionId = paramMap.get('subscriptionId');
      if (projectName && integrationId) {
        return of({
          projectName,
          integrationId,
          subscriptionId: subscriptionId ?? undefined,
          editMode: !!subscriptionId,
        } as Params);
      }
      return EMPTY;
    }),
    shareReplay(1)
  );

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private router: Router,
    private notificationsService: NotificationsService,
    private eventService: EventService,
    private _changeDetectorRef: ChangeDetectorRef
  ) {
    const subscription$ = this.params$.pipe(
      switchMap((params) => {
        if (params.subscriptionId) {
          return this.dataService.getUniformSubscription(params.integrationId, params.subscriptionId);
        }
        return of(new UniformSubscription(params.projectName));
      }),
      tap((subscription) => {
        if (subscription.id) {
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
      shareReplay(1)
    );

    const integrationInfo$ = this.params$.pipe(
      switchMap((params) => this.dataService.getUniformRegistrationInfo(params.integrationId)),
      tap((info) => {
        if (info.isControlPlane) {
          return;
        }
        this.suffixes = [
          {
            value: EventState.TRIGGERED,
            displayValue: EventState.TRIGGERED,
          },
        ];
      }),
      shareReplay(1)
    );

    const isWebhookService$ = integrationInfo$.pipe(map((info) => info.isWebhookService));

    const taskNames$ = this.params$.pipe(
      switchMap((params) => this.dataService.getTaskNames(params.projectName)),
      catchError((err: HttpErrorResponse) => {
        this.errorMessage = err.error;
        this.notificationsService.addNotification(NotificationType.ERROR, err.error);
        return throwError(() => err);
      })
    );

    const webhook$: Observable<IWebhookConfigClient | undefined> = combineLatest([
      this.params$,
      subscription$,
      integrationInfo$,
    ]).pipe(
      switchMap(([params, subscription, integrationInfo]) => {
        if (integrationInfo.isWebhookService && params.editMode && subscription.id) {
          const stage: string | undefined = subscription.filter?.stages?.[0];
          const services: string | undefined = subscription.filter?.services?.[0];
          this.updateEventPayload(
            params.projectName,
            subscription.filter?.stages ?? [],
            subscription.filter?.services ?? [],
            integrationInfo.isWebhookService
          );
          return this.dataService.getWebhookConfig(subscription.id, params.projectName, stage, services);
        }
        return of(undefined);
      })
    );

    const webhookSecrets$ = integrationInfo$.pipe(
      switchMap((info) => {
        if (info.isWebhookService) {
          return this.dataService.getSecretsForScope(SecretScopeDefault.WEBHOOK);
        }
        return of(undefined);
      })
    );

    const filter$ = this.params$.pipe(switchMap((params) => this.dataService.getSequenceFilter(params.projectName)));

    this.data$ = combineLatest([taskNames$, subscription$, filter$, isWebhookService$, webhook$, webhookSecrets$]).pipe(
      map(([taskNames, subscription, filterData, isWebhookService, webhook, webhookSecrets]) => {
        return {
          taskNames,
          subscription,
          filter: filterData,
          isWebhookService,
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
        const contextArray = data.context as unknown[];
        this.deleteSubscription(contextArray[0] as Params, contextArray[1] as SubscriptionState);
      }
    });
  }

  public deleteSubscription(params: Params, data: SubscriptionState): void {
    if (!data.subscription.id) return;
    this.eventService.deletionProgressEvent.next({ isInProgress: true });

    this.dataService
      .deleteSubscription(params.integrationId, data.subscription.id, data.isWebhookService)
      .pipe(
        map((): DeletionProgressEvent => ({ isInProgress: false, result: DeleteResult.SUCCESS })),
        catchError(handleDeletionError('Subscription'))
      )
      .subscribe((progressEvent) => {
        this.eventService.deletionProgressEvent.next(progressEvent);
        if (progressEvent.result === DeleteResult.SUCCESS) {
          this.resetForms();
          this.router.navigate([
            '/',
            'project',
            params.projectName,
            'settings',
            'uniform',
            'integrations',
            params.integrationId,
          ]);
        }
      });
  }

  private updateEventPayload(
    projectName: string,
    stages: string[],
    services: string[],
    isWebhookService: boolean
  ): void {
    const shouldUpdateEventPayload = isWebhookService && this.taskControl.value && this.taskSuffixControl.value;
    if (!shouldUpdateEventPayload) {
      return;
    }

    const event = `${EventTypes.PREFIX}${this.taskControl.value}`;
    const eventSuffix = this.taskSuffixControl.value;
    this.eventPayload = undefined;
    this.dataService
      .getIntersectedEvent(event, eventSuffix, projectName, stages, services)
      .subscribe((intersectedEvent: Record<string, unknown>) => {
        this.eventPayload = Object.keys(intersectedEvent).length ? intersectedEvent : Trace.defaultTrace;
      });
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
    editMode: boolean,
    projectName: string,
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: IWebhookConfigClient
  ): void {
    this.updating = true;
    const isShortPrefix = this.taskControl.value === EventTypes.PREFIX_SHORT && this.taskSuffixControl.value === '>';
    subscription.event = isShortPrefix
      ? `${this.taskControl.value}.${this.taskSuffixControl.value}`
      : `${EventTypes.PREFIX}${this.taskControl.value}.${this.taskSuffixControl.value}`;
    subscription.setIsGlobal(this.isGlobalControl.value, projectName);

    if (webhookConfig) {
      webhookConfig.type = subscription.event;
      webhookConfig.prevConfiguration = this._previousFilter;
    }

    const update$ = editMode
      ? this.dataService.updateUniformSubscription(integrationId, subscription, webhookConfig)
      : this.dataService.createUniformSubscription(integrationId, subscription, webhookConfig);

    update$
      .pipe(
        map(() => ({
          type: NotificationType.SUCCESS,
          message: 'Subscription successfully created!',
        })),
        catchError(() =>
          of({
            type: NotificationType.ERROR,
            message: 'The subscription could not be updated',
          })
        ),
        finalize(() => (this.updating = false))
      )
      .subscribe((notificationResult: NotificationResult) => {
        this.notificationsService.addNotification(notificationResult.type, notificationResult.message);
        if (notificationResult.type === NotificationType.SUCCESS) {
          this.resetForms();
          this.router.navigate(['/', 'project', projectName, 'settings', 'uniform', 'integrations', integrationId]);
        }
      });
  }

  public isFormValid(subscription: UniformSubscription): boolean {
    return (
      this.subscriptionForm.valid &&
      (!!subscription.filter.stages?.length || !subscription.filter.services?.length) &&
      this.isWebhookFormValid &&
      !this.updating
    );
  }

  public selectedTaskChanged(projectName: string, subscription: UniformSubscription, isWebhookService: boolean): void {
    this.updateEventPayload(
      projectName,
      subscription.filter.stages ?? [],
      subscription.filter.services ?? [],
      isWebhookService
    );
  }

  public subscriptionFilterChanged(
    subscription: UniformSubscription,
    projectName: string,
    isWebhookService: boolean
  ): void {
    this.isFilterDirty = !!subscription.filter.stages?.length || !!subscription.filter.services?.length;
    this.updateIsGlobalCheckbox(subscription);
    this.updateEventPayload(
      projectName,
      subscription.filter.stages ?? [],
      subscription.filter.services ?? [],
      isWebhookService
    );
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
    if (!this.taskControl.value) {
      return undefined;
    }
    const suffixDiffers = this.taskSuffixControl.value && this.taskSuffixControl.value !== this.suffixes[0].value;
    const suffix = suffixDiffers ? this.taskSuffixControl.value : 'triggered';
    return `${EventTypes.PREFIX}${this.taskControl.value}.${suffix}`;
  }

  public webhookFormValidityChanged(isValid: boolean): void {
    this.isWebhookFormValid = isValid;
    this._changeDetectorRef.detectChanges();
  }

  // @HostListener allows us to also guard against browser refresh, close, etc.
  @HostListener('window:beforeunload', ['$event'])
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
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
