import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { catchError, filter, map, mergeMap, takeUntil, withLatestFrom } from 'rxjs/operators';
import { of, Subject } from 'rxjs';
import { DeleteData, DeleteResult, DeleteType, DeletionProgressEvent } from '../../../_interfaces/delete';
import { EventService } from '../../../_services/event.service';
import { DataService } from '../../../_services/data.service';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationsService } from '../../../_services/notifications.service';
import { NotificationType } from '../../../_models/notification';

@Component({
  selector: 'ktb-edit-service',
  templateUrl: './ktb-edit-service.component.html',
  styleUrls: ['./ktb-edit-service.component.scss'],
})
export class KtbEditServiceComponent implements OnDestroy {
  public params$ = this.route.paramMap.pipe(
    map((params) => ({
      serviceName: params.get('serviceName'),
      projectName: params.get('projectName'),
    })),
    filter(
      (params): params is { serviceName: string; projectName: string } => !!params.serviceName && !!params.projectName
    )
  );

  public project$ = this.params$.pipe(mergeMap((params) => this.dataService.loadPlainProject(params.projectName)));

  public fileTree$ = this.params$.pipe(
    mergeMap((params) => this.dataService.getFileTreeForService(params.projectName, params.serviceName))
  );

  private unsubscribe$: Subject<void> = new Subject<void>();

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private eventService: EventService,
    private dataService: DataService,
    private notificationsService: NotificationsService
  ) {
    this.eventService.deletionTriggeredEvent
      .pipe(
        withLatestFrom(this.params$),
        filter(([event, params]) => event.type === DeleteType.SERVICE && event.name === params.serviceName),
        takeUntil(this.unsubscribe$)
      )
      .subscribe(([, params]) => {
        this.eventService.deletionProgressEvent.next({ isInProgress: true });
        this.deleteService(params.projectName, params.serviceName);
      });
  }

  private deleteService(projectName: string, serviceName: string): void {
    this.dataService
      .deleteService(projectName, serviceName)
      .pipe(
        map(() => ({ isInProgress: false, result: DeleteResult.SUCCESS })),
        catchError((error: HttpErrorResponse) => {
          return of({
            isInProgress: false,
            result: DeleteResult.ERROR,
            error: error.error,
          });
        })
      )
      .subscribe(async (event: DeletionProgressEvent) => {
        this.eventService.deletionProgressEvent.next(event);
        if (event.result === DeleteResult.SUCCESS) {
          this.dataService.loadProject(projectName);
          await this.router.navigate(['../../'], { relativeTo: this.route });
          this.notificationsService.addNotification(NotificationType.SUCCESS, 'Service deleted');
        }
      });
  }

  public getServiceDeletionData(serviceName: string): DeleteData {
    return {
      type: DeleteType.SERVICE,
      name: serviceName,
    };
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}
