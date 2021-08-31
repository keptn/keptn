import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { filter, map, takeUntil } from 'rxjs/operators';
import { Subject } from 'rxjs';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { EventService } from '../../_services/event.service';
import { DataService } from '../../_services/data.service';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';

@Component({
  selector: 'ktb-edit-service',
  templateUrl: './ktb-edit-service.component.html',
  styleUrls: ['./ktb-edit-service.component.scss'],
})
export class KtbEditServiceComponent implements OnDestroy {
  public serviceName?: string;
  private projectName?: string;
  private unsubscribe$: Subject<void> = new Subject<void>();

  constructor(private route: ActivatedRoute, private eventService: EventService, private dataService: DataService, private router: Router, private notificationsService: NotificationsService) {
    this.route.paramMap.pipe(
      map(params => {
        return {
          serviceName: params.get('serviceName'),
          projectName: params.get('projectName'),
        };
      }),
      filter((params): params is { serviceName: string, projectName: string } => !!params.serviceName && !!params.projectName),
    ).subscribe(params => {
      this.serviceName = params.serviceName;
      this.projectName = params.projectName;
    });

    this.eventService.deletionTriggeredEvent.pipe(
      filter(event => event.type === DeleteType.SERVICE && event.name === this.serviceName),
      takeUntil(this.unsubscribe$),
    ).subscribe(() => {
      this.eventService.deletionProgressEvent.next({isInProgress: true});
      this.deleteService();
    });
  }

  private deleteService(): void {
    const projectName = this.projectName;
    if (this.serviceName && projectName) {
      this.dataService.deleteService(projectName, this.serviceName).subscribe(async () => {
        this.eventService.deletionProgressEvent.next({isInProgress: false, result: DeleteResult.SUCCESS});
        this.dataService.loadProject(projectName);
        await this.router.navigate(['../../'], {relativeTo: this.route});
        this.notificationsService.addNotification(NotificationType.Success, 'Service deleted', 5_000);
      }, (error: HttpErrorResponse) => {
        this.eventService.deletionProgressEvent.next({isInProgress: false, result: DeleteResult.ERROR, error: error.error});
      });
    }
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
