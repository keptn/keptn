import { Component } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { forkJoin, Observable, of } from 'rxjs';
import { filter, map, switchMap, take, tap } from 'rxjs/operators';
import { UniformSubscription } from '../../_models/uniform-subscription';
import { DtFilterFieldDefaultDataSource } from '@dynatrace/barista-components/filter-field';
import { Project } from '../../_models/project';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { DtFilterFieldDefaultDataSourceAutocomplete } from '@dynatrace/barista-components/filter-field/src/filter-field-default-data-source';
import { EventTypes } from '../../../../shared/interfaces/event-types';

@Component({
  selector: 'ktb-modify-uniform-subscription',
  templateUrl: './ktb-modify-uniform-subscription.component.html',
})
export class KtbModifyUniformSubscriptionComponent {
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
    isGlobal: this.isGlobalControl
  });
  public readonly affixes: { value: string, displayValue: string }[] = [
    {
      value: '>',
      displayValue: '*'
    },
    {
      value: 'triggered',
      displayValue: 'triggered'
    },
    {
      value: 'started',
      displayValue: 'started'
    },
    {
      value: 'finished',
      displayValue: 'finished'
    }];

  constructor(private route: ActivatedRoute, private dataService: DataService, private router: Router) {
    const subscription$ = this.route.paramMap.pipe(
      map(paramMap => {
        return {
          integrationId: paramMap.get('integrationId'),
          subscriptionId: paramMap.get('subscriptionId'),
          projectName: paramMap.get('projectName')
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
      take(1)
    );

    const integrationId$ = this.route.paramMap
      .pipe(
        map(paramMap => paramMap.get('integrationId')),
        filter((integrationId: string | null): integrationId is string => !!integrationId),
        take(1)
      );

    const projectName$ = this.route.paramMap
      .pipe(
        map(paramMap => paramMap.get('projectName')),
        filter((projectName: string | null): projectName is string => !!projectName)
      );

    const taskNames$ = projectName$
      .pipe(
        switchMap(projectName => this.dataService.getTaskNames(projectName)),
        take(1)
      );
    const project$ = projectName$
      .pipe(
        switchMap(projectName => this.dataService.getProject(projectName)),
        filter((project?: Project): project is Project => !!project),
        tap(project => this.updateDataSource(project)),
        take(1)
      );

    this.data$ = forkJoin({
      taskNames: taskNames$,
      subscription: subscription$,
      project: project$,
      integrationId: integrationId$
    });
  }

  private updateDataSource(project: Project): void {
    this._dataSource.data = {
      autocomplete: [
        {
          name: 'Stage',
          autocomplete: project.stages.map(stage => {
            return {
              name: stage.stageName
            };
          })
        },
        {
          name: 'Service',
          autocomplete: project.getServices().map(service => {
            return {
              name: service.serviceName
            };
          })
        }
      ]
    } as DtFilterFieldDefaultDataSourceAutocomplete;
  }

  public updateSubscription(projectName: string, integrationId: string, subscription: UniformSubscription): void {
    this.updating = true;
    let update;
    subscription.event = `${EventTypes.PREFIX}${this.taskControl.value}.${this.taskSuffixControl.value}`;
    subscription.setIsGlobal(this.isGlobalControl.value, projectName);

    if (this.editMode) {
      update = this.dataService.updateUniformSubscription(integrationId, subscription);
    } else {
      update = this.dataService.createUniformSubscription(integrationId, subscription);
    }
    update.subscribe(() => {
      this.updating = false;
      this.router.navigate(['/', 'project', projectName, 'uniform', 'services', integrationId]);
    });
  }
}
