import { Component, OnDestroy } from '@angular/core';
import { Subject } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../_services/data.service';
import { filter, map, takeUntil } from 'rxjs/operators';
import { FormControl, FormGroup, Validators } from '@angular/forms';
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
  public isLoading = true;
  private unsubscribe$: Subject<void> = new Subject<void>();
  private redirectTo?: string;
  public formGroup: FormGroup = new FormGroup({
    serviceName: this.serviceNameControl,
  });

  constructor(
    private route: ActivatedRoute,
    private dataService: DataService,
    private router: Router,
    private notificationsService: NotificationsService
  ) {
    this.route.queryParamMap
      .pipe(
        map((params) => params.get('redirectTo')),
        takeUntil(this.unsubscribe$),
        filter((redirectTo): redirectTo is string => !!redirectTo)
      )
      .subscribe((redirectTo) => {
        this.redirectTo = redirectTo;
      });

    this.route.paramMap
      .pipe(
        map((params) => params.get('projectName')),
        filter((projectName): projectName is string => !!projectName),
        takeUntil(this.unsubscribe$)
      )
      .subscribe((projectName) => {
        this.projectName = projectName;
        this.isLoading = true;

        this.dataService.getServiceNames(this.projectName).subscribe((serviceNames) => {
          this.isLoading = false;

          this.serviceNameControl.setValidators([
            Validators.required,
            FormUtils.nameExistsValidator(serviceNames),
            Validators.pattern('[a-z]([a-z]|[0-9]|-)*'),
          ]);
        });
      });
  }

  public createService(projectName: string): void {
    this.isCreating = true;
    this.dataService.createService(projectName, this.serviceNameControl.value).subscribe(
      async () => {
        this.isCreating = false;
        if (this.projectName) {
          this.dataService.loadProject(this.projectName);
        }
        await this.cancel();
        this.notificationsService.addNotification(NotificationType.SUCCESS, 'Service successfully created!');
      },
      () => {
        this.isCreating = false;
      }
    );
  }

  public async cancel(): Promise<void> {
    if (this.redirectTo) {
      await this.router.navigateByUrl(this.redirectTo);
    } else {
      await this.router.navigate(['../'], { relativeTo: this.route });
    }
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
