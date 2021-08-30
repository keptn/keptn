import { Component, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { Project } from '../../_models/project';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { HttpErrorResponse } from '@angular/common/http';
import { FormUtils } from '../../_utils/form.utils';

@Component({
  selector: 'ktb-create-service',
  templateUrl: './ktb-create-service.component.html',
})
export class KtbCreateServiceComponent implements OnDestroy {
  public projectName?: string;
  public isCreating = false;
  public serviceNameControl: FormControl = new FormControl();
  private unsubscribe$: Subject<void> = new Subject<void>();
  public formGroup: FormGroup = new FormGroup({
    serviceName: this.serviceNameControl,
  });

  constructor(private route: ActivatedRoute, private dataService: DataService, private router: Router, private notificationsService: NotificationsService) {
    this.route.paramMap.pipe(
      map(params => params.get('projectName')),
      filter((projectName: string | null): projectName is string => !!projectName),
      switchMap(projectName => this.dataService.getProject(projectName)),
      takeUntil(this.unsubscribe$),
      filter((project?: Project): project is Project => !!project),
    ).subscribe(project => {
      this.projectName = project.projectName;
      const serviceNames = project.services?.map(service => service.serviceName) ?? [];
      this.serviceNameControl.setValidators([Validators.required, FormUtils.nameExistsValidator(serviceNames), Validators.pattern('[a-z]([a-z]|[0-9]|-)*')]);
    });
  }

  public createService(projectName: string): void {
    this.isCreating = true;
    this.dataService.createService(projectName, this.serviceNameControl.value).subscribe(async () => {
      this.isCreating = false;
      await this.router.navigate(['../'], {relativeTo: this.route});
      this.notificationsService.addNotification(NotificationType.Success, 'Service successfully created!', 5_000);
    }, (error: HttpErrorResponse) => {
      this.notificationsService.addNotification(NotificationType.Error, error.error, 5_000);
      this.isCreating = false;
    });
  }

  public ngOnDestroy(): void {
    this.notificationsService.clearNotifications();
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
