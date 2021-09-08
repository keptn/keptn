import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { filter, map, takeUntil } from 'rxjs/operators';
import { Observable, Subject } from 'rxjs';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { EventService } from '../../_services/event.service';
import { DataService } from '../../_services/data.service';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { Project } from '../../_models/project';
import { ServiceResource } from '../../../../shared/interfaces/serviceResource';

@Component({
  selector: 'ktb-edit-service',
  templateUrl: './ktb-edit-service.component.html',
  styleUrls: ['./ktb-edit-service.component.scss'],
})
export class KtbEditServiceComponent implements OnDestroy {
  public serviceName?: string;
  public project$?: Observable<Project | undefined>;
  public resources$?: Observable<ServiceResource[] | undefined>;
  private projectName?: string;
  private unsubscribe$: Subject<void> = new Subject<void>();
  public fileTree: TreeEntry[] | undefined;

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

      this.project$ = this.dataService.getProject(this.projectName);
    });

    this.eventService.deletionTriggeredEvent.pipe(
      filter(event => event.type === DeleteType.SERVICE && event.name === this.serviceName),
      takeUntil(this.unsubscribe$),
    ).subscribe(() => {
      this.eventService.deletionProgressEvent.next({isInProgress: true});
      this.deleteService();
    });

    this.project$?.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(project => {
      if (!this.fileTree && project && this.serviceName) {
        this.dataService.getServiceResourceForAllStages(project.projectName, this.serviceName).subscribe((resources) => {
          this.fileTree = [];
          project.stages.forEach(stage => {
            const resourcesForStage = this.getResourcesForStage(resources, stage.stageName);
            this.fileTree?.push({stage: stage.stageName, files: this.getTree(resourcesForStage)});
          });
          console.log(this.fileTree);
        });
      }
    });
  }

  private getTree(resources: ServiceResource[]): TreeFileEntry[] {
    const tree = new Map();

    resources.forEach((resource) => {
      const uriParts = resource.resourceURI.split('/');
      const file = uriParts.pop();
      const folder = uriParts.join('/');

      const elem = tree.get(folder);
      if (!elem) {
        tree.set(folder, [file]);
      } else {
        elem.push(file);
        tree.set(folder, elem);
      }
    });

    const transformedTree: TreeFileEntry[] = [];

    tree.forEach((val: string[], key: string) => {
      transformedTree.push({folder: key, files: val});
    });

    return transformedTree;
  }

  public getLinkForStage(remoteUri: string | undefined, stageName: string): string {
    if (remoteUri) {
      if (remoteUri.includes('github.') || remoteUri.includes('gitlab.')) {
        return remoteUri + '/tree/' + stageName + '/' + this.serviceName;
      }
      if (remoteUri.includes('bitbucket.')) {
        return remoteUri + '/src/' + stageName + '/' + this.serviceName;
      }
      if (remoteUri.includes('azure.')) {
        return remoteUri + '?path=' + this.serviceName + '&version=GB' + stageName;
      }
      if (remoteUri.includes('git-codecommit.')) {
        const repoParts = remoteUri.split('/');
        const region = repoParts.find(part => part.includes('git-codecommit.'))?.split('.')[1];
        const repoName = repoParts[repoParts.length - 1];
        return 'https://' + region + '.console.aws.amazon.com/codesuite/codecommit/repositories/' + repoName + '/browse/refs/heads/' + stageName;
      }

      return remoteUri;
    }
    return '';
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

  public getResourcesForStage(resources: ServiceResource[], stageName: string): ServiceResource[] {
    return resources.filter(resource => resource.stageName === stageName);
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}

export interface TreeEntry {
  stage: string;
  files: TreeFileEntry[];
}

export interface TreeFileEntry {
  folder: string;
  files: string[];
}
