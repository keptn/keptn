import { Component } from '@angular/core';
import { mergeMap, Observable, of } from 'rxjs';
import { ActivatedRoute, Router } from '@angular/router';
import { DataService } from '../../../_services/data.service';
import { catchError, filter, finalize, map, tap } from 'rxjs/operators';
import { FormControl, FormGroup, Validators } from '@angular/forms';
import { NotificationsService } from '../../../_services/notifications.service';
import { NotificationType } from '../../../_models/notification';
import { FormUtils } from '../../../_utils/form.utils';

@Component({
  selector: 'ktb-create-service',
  templateUrl: './ktb-create-service.component.html',
})
export class KtbCreateServiceComponent {
  public isLoading = true;
  public isCreating = false;

  public projectName$: Observable<string> = this.route.paramMap.pipe(
    map((params) => params.get('projectName')),
    filter((projectName): projectName is string => !!projectName)
  );

  public serviceNames$: Observable<string[]> = this.projectName$.pipe(
    mergeMap((projectName) =>
      this.dataService.getServiceNames(projectName).pipe(
        tap(this.setValidators),
        finalize(() => (this.isLoading = false))
      )
    )
  );

  public redirectTo$ = this.route.queryParamMap.pipe(
    map((params) => params.get('redirectTo')),
    map((value) => ({
      value,
    }))
  );

  public serviceNameControl: FormControl = new FormControl();
  public formGroup = new FormGroup({
    serviceName: this.serviceNameControl,
  });

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private dataService: DataService,
    private notificationsService: NotificationsService
  ) {}

  private setValidators = (services: string[]): void =>
    this.serviceNameControl.setValidators([
      Validators.required,
      FormUtils.nameExistsValidator(services),
      Validators.pattern('[a-z]([a-z]|[0-9]|-)*'),
    ]);

  public createService(projectName: string, redirectTo: string | null): void {
    this.isCreating = true;
    this.dataService
      .createService(projectName, this.serviceNameControl.value)
      .pipe(
        map(() => true),
        catchError(() => of(false)),
        filter((success) => success),
        finalize(() => (this.isCreating = false))
      )
      .subscribe(async () => {
        this.dataService.loadProject(projectName);
        await this.cancel(redirectTo);
        this.notificationsService.addNotification(NotificationType.SUCCESS, 'Service successfully created!');
      });
  }

  public async cancel(redirectTo: string | null): Promise<void> {
    if (redirectTo) {
      await this.router.navigateByUrl(redirectTo);
    } else {
      await this.router.navigate(['../'], { relativeTo: this.route });
    }
  }
}
