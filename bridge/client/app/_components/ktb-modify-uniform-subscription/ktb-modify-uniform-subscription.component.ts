import { ChangeDetectorRef, Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { forkJoin, Observable, of, Subject, throwError } from 'rxjs';
import { filter, map, switchMap, take, takeUntil, tap, catchError } from 'rxjs/operators';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { Project } from '../../_models/project';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { EventTypes } from '../../../../shared/interfaces/event-types';
import { UniformRegistration } from '../../_models/uniform-registration';
import { WebhookConfig } from '../../../../shared/models/webhook-config';
import { AppUtils } from '../../_utils/app.utils';
import { PreviousWebhookConfig } from '../../../../shared/interfaces/webhook-config';
import { NotificationsService } from '../../_services/notifications.service';
import { UniformRegistrationInfo } from '../../../../shared/interfaces/uniform-registration-info';
import { NotificationType } from '../../_models/notification';
import { Secret } from '../../_models/secret';
import { SecretScope } from '../../../../shared/interfaces/secret-scope';
import { EventState } from '../../../../shared/models/event-state';
import { Trace } from '../../_models/trace';
import { HttpErrorResponse } from '@angular/common/http';

@Component({
  selector: 'ktb-modify-uniform-subscription',
  templateUrl: './ktb-modify-uniform-subscription.component.html',
  providers: [NotificationsService],
  styleUrls: ['./ktb-modify-uniform-subscription.component.scss'],
})
export class KtbModifyUniformSubscriptionComponent implements OnDestroy {
  private readonly unsubscribe$: Subject<void> = new Subject<void>();
  private taskControl = new FormControl('', [Validators.required]);
  public eventPayload: Record<string, unknown> | undefined;
  public taskSuffixControl = new FormControl('', [Validators.required]);
  private isGlobalControl = new FormControl();
  public data$: Observable<{
    taskNames: string[];
    subscription: UniformSubscription;
    project: Project;
    integrationId: string;
    webhook?: WebhookConfig;
    webhookSecrets?: Secret[];
  }>;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public editMode = false;
  public updating = false;
  public subscriptionForm = new FormGroup({
    taskPrefix: this.taskControl,
    taskSuffix: this.taskSuffixControl,
    isGlobal: this.isGlobalControl,
  });
  private _previousFilter?: PreviousWebhookConfig;
  public uniformRegistration?: UniformRegistration;
  public isWebhookFormValid = true;
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

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private router: Router,
    private notificationsService: NotificationsService,
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
    const project$ = projectName$.pipe(
      switchMap((projectName) => this.dataService.getProject(projectName)),
      filter((project?: Project): project is Project => !!project),
      tap((project) => this.updateDataSource(project)),
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
          let webhook: Observable<WebhookConfig | undefined>;
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
          return this.dataService.getSecretsForScope(SecretScope.WEBHOOK);
        }
        return of(undefined);
      })
    );

    this.data$ = forkJoin({
      taskNames: taskNames$,
      subscription: subscription$,
      project: project$,
      integrationId: integrationId$,
      webhook: webhook$,
      webhookSecrets: webhookSecrets$,
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

  private updateDataSource(project: Project): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: project.stages.map((stage) => ({
            name: stage.stageName,
          })),
        },
        {
          name: 'Service',
          autocomplete: project.getServices().map((service) => ({
            name: service.serviceName,
          })),
        },
      ],
    } as DtFilterFieldDefaultDataSourceAutocomplete;
  }

  public updateSubscription(
    projectName: string,
    integrationId: string,
    subscription: UniformSubscription,
    webhookConfig?: WebhookConfig
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
      webhookConfig.filter = subscription.filter;
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
        this.router.navigate(['/', 'project', projectName, 'uniform', 'services', integrationId]);
      },
      () => {
        this.notificationsService.addNotification(
          NotificationType.ERROR,
          'The subscription could not be updated',
          5_000
        );
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

  public ngOnDestroy(): void {
    this.notificationsService.clearNotifications();
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
