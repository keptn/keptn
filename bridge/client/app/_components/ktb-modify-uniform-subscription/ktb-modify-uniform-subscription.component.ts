import { Component, OnDestroy, ViewChild } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { combineLatest, forkJoin, Observable, of, Subject } from 'rxjs';
import { filter, map, switchMap, take, takeUntil, tap } from 'rxjs/operators';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { Project } from '../../_models/project';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { EventTypes } from '../../../../shared/interfaces/event-types';
import { UniformRegistration } from '../../_models/uniform-registration';
import { KeptnService } from '../../../../shared/models/keptn-service';
import { KtbWebhookSettingsComponent } from '../ktb-webhook-settings/ktb-webhook-settings.component';
import { WebhookConfig } from '../../_models/webhook-config';

@Component({
  selector: 'ktb-modify-uniform-subscription',
  templateUrl: './ktb-modify-uniform-subscription.component.html',
})
export class KtbModifyUniformSubscriptionComponent implements OnDestroy {
  private readonly unsubscribe$ = new Subject<void>();
  private taskControl = new FormControl('', [Validators.required]);
  private taskSuffixControl = new FormControl('', [Validators.required]);
  private isGlobalControl = new FormControl();
  public data$: Observable<{ taskNames: string[], subscription: UniformSubscription, project: Project, integrationId: string }>;
  public _dataSource = new DtFilterFieldDefaultDataSource();
  public editMode = false;
  public updating = false;
  public subscriptionForm = new FormGroup({
    taskPrefix: this.taskControl,
    taskSuffix: this.taskSuffixControl,
    isGlobal: this.isGlobalControl,
  });
  private webhookSettings?: KtbWebhookSettingsComponent;
  public uniformRegistration?: UniformRegistration;
  public suffixes: { value: string, displayValue: string }[] = [
    {
      value: '>',
      displayValue: '*',
    },
    {
      value: 'triggered',
      displayValue: 'triggered',
    },
    {
      value: 'started',
      displayValue: 'started',
    },
    {
      value: 'finished',
      displayValue: 'finished',
    }];

  @ViewChild('webhookSettings', { static: false }) set webhookSettingsElement(webhookSettings: KtbWebhookSettingsComponent) {
    if (webhookSettings) { // initially setter gets called with undefined
      this.webhookSettings = webhookSettings;
    }
  }

  constructor(private route: ActivatedRoute, private dataService: DataService, private router: Router) {
    const subscription$ = this.route.paramMap.pipe(
      map(paramMap => {
        return {
          integrationId: paramMap.get('integrationId'),
          subscriptionId: paramMap.get('subscriptionId'),
          projectName: paramMap.get('projectName'),
        };
      }),
      filter((params): params is  { integrationId: string, subscriptionId: string | null, projectName: string } => !!(params.integrationId && params.projectName)),
      switchMap(params => {
        this.editMode = !!params.subscriptionId;
        if (params.subscriptionId) {
          return this.dataService.getUniformSubscription(params.integrationId, params.subscriptionId);
        } else {
          return of(new UniformSubscription(params.projectName));
        }
      }),
      tap(subscription => {
        this.taskControl.setValue(subscription.prefix);
        this.taskSuffixControl.setValue(subscription.suffix);
        this.isGlobalControl.setValue(subscription.isGlobal);
      }),
      take(1),
    );

    const integrationId$ = this.route.paramMap
      .pipe(
        map(paramMap => paramMap.get('integrationId')),
        filter((integrationId: string | null): integrationId is string => !!integrationId),
        take(1),
      );

    integrationId$.pipe(
      switchMap(integrationId => this.dataService.getIsUniformRegistrationControlPlane(integrationId)),
    ).subscribe(isControlPlane => {
      if (!isControlPlane) {
        this.suffixes = [
          {
            value: 'triggered',
            displayValue: 'triggered',
          },
        ];
      }
    });

    const registrations$ = this.dataService.getUniformRegistrations();
    combineLatest([registrations$, integrationId$])
      .pipe(takeUntil(this.unsubscribe$))
      .subscribe(([uniformRegistrations, integrationId]) => {
        const uniformRegistration = uniformRegistrations.find(uR => uR.id === integrationId);
        this.setUniformRegistration(uniformRegistration);
      });

    const projectName$ = this.route.paramMap
      .pipe(
        map(paramMap => paramMap.get('projectName')),
        filter((projectName: string | null): projectName is string => !!projectName),
      );

    const taskNames$ = projectName$
      .pipe(
        switchMap(projectName => this.dataService.getTaskNames(projectName)),
        take(1),
      );
    const project$ = projectName$
      .pipe(
        switchMap(projectName => this.dataService.getProject(projectName)),
        filter((project?: Project): project is Project => !!project),
        tap(project => this.updateDataSource(project)),
        take(1),
      );

    this.data$ = forkJoin({
      taskNames: taskNames$,
      subscription: subscription$,
      project: project$,
      integrationId: integrationId$,
    });
  }

  private setUniformRegistration(uniformRegistration: UniformRegistration | undefined): void {
    if(this.uniformRegistration !== uniformRegistration) {
      this.uniformRegistration = uniformRegistration;
    }
  }

  private updateDataSource(project: Project): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: project.stages.map(stage => {
            return {
              name: stage.stageName,
            };
          }),
        },
        {
          name: 'Service',
          autocomplete: project.getServices().map(service => {
            return {
              name: service.serviceName,
            };
          }),
        },
      ],
    } as DtFilterFieldDefaultDataSourceAutocomplete;
  }

  public updateSubscription(projectName: string, integrationId: string, subscription: UniformSubscription): void {
    this.updating = true;
    const updates = [];
    subscription.event = `${EventTypes.PREFIX}${this.taskControl.value}.${this.taskSuffixControl.value}`;
    subscription.setIsGlobal(this.isGlobalControl.value, projectName);

    if (this.editMode) {
      updates.push(this.dataService.updateUniformSubscription(integrationId, subscription));
    } else {
      updates.push(this.dataService.createUniformSubscription(integrationId, subscription));
    }

    if (this.isWebhookService()) {
      const webhookSettingsForm = this.webhookSettings?.webhookConfigForm;
      if (webhookSettingsForm && webhookSettingsForm.valid) {
        const webhookConfig: WebhookConfig = new WebhookConfig();
        webhookConfig.type = subscription.event;
        webhookConfig.filter = subscription.filter;
        webhookConfig.prevFilter = this.webhookSettings?.prevFilter;
        webhookConfig.method = webhookSettingsForm.get('method')?.value;
        webhookConfig.url = webhookSettingsForm.get('url')?.value;
        webhookConfig.payload = webhookSettingsForm.get('payload')?.value;
        webhookConfig.proxy = webhookSettingsForm.get('proxy')?.value;
        for (const header of this.webhookSettings?.headerControls || []) {
          webhookConfig.header?.push({
            name: header.get('name')?.value,
            value: header.get('value')?.value,
          });
        }
        updates.push(this.dataService.saveWebhookConfig(webhookConfig));
      }
    }

    forkJoin(
      ...updates
    ).subscribe(() => {
      this.updating = false;
      this.router.navigate(['/', 'project', projectName, 'uniform', 'services', integrationId]);
    }, err => {
      this.updating = false;
    });
  }

  public isWebhookService(): boolean {
    return this?.uniformRegistration?.name === KeptnService.WEBHOOK_SERVICE;
  }

  ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
