import { Component, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { filter, map, switchMap, takeUntil } from 'rxjs/operators';
import { Observable, Subject } from 'rxjs';
import { DeleteData, DeleteResult, DeleteType } from '../../_interfaces/delete';
import { EventService } from '../../_services/event.service';
import { DataService } from '../../_services/data.service';
import { HttpErrorResponse } from '@angular/common/http';
import { NotificationsService } from '../../_services/notifications.service';
import { NotificationType } from '../../_models/notification';
import { Project } from '../../_models/project';
import { ServiceResource } from '../../../../shared/interfaces/serviceResource';
import { TreeEntry } from '../../ktb-edit-service-file-list/ktb-edit-service-file-list.component';

@Component({
  selector: 'ktb-edit-service',
  templateUrl: './ktb-edit-service.component.html',
  styleUrls: ['./ktb-edit-service.component.scss'],
})
export class KtbEditServiceComponent implements OnDestroy {
  public serviceName?: string;
  public project$?: Observable<Project | undefined>;
  private projectName?: string;
  private unsubscribe$: Subject<void> = new Subject<void>();
  public treeLoading = false;
  public fileTree: FileTree[] = [];

  constructor(private route: ActivatedRoute, private eventService: EventService, private dataService: DataService, private router: Router, private notificationsService: NotificationsService) {
    const params$ = this.route.paramMap.pipe(
      map(params => {
        return {
          serviceName: params.get('serviceName'),
          projectName: params.get('projectName'),
        };
      }),
      filter((params): params is { serviceName: string, projectName: string } => !!params.serviceName && !!params.projectName),
    );

    params$.subscribe(params => {
      this.serviceName = params.serviceName;
      this.projectName = params.projectName;
    });

    this.project$ = params$.pipe(
      switchMap(params => this.dataService.getProject(params.projectName)),
    );

    this.project$.pipe(
      takeUntil(this.unsubscribe$),
    ).subscribe(project => {
      if (this.fileTree.length === 0 && project && project.gitRemoteURI && this.serviceName) {
        this.treeLoading = true;
        this.getResourcesAndTransform(project, this.serviceName);
      }
    });

    this.eventService.deletionTriggeredEvent.pipe(
      filter(event => event.type === DeleteType.SERVICE && event.name === this.serviceName),
      takeUntil(this.unsubscribe$),
    ).subscribe(() => {
      this.eventService.deletionProgressEvent.next({isInProgress: true});
      this.deleteService();
    });
  }

  public getResourcesAndTransform(project: Project, serviceName: string): void {
    this.dataService.getServiceResourceForAllStages(project.projectName, serviceName).subscribe((resources) => {
      project.stages.forEach(stage => {
        const fileTreeObj: FileTree = {
          stageName: stage.stageName,
          tree: [],
        };
        fileTreeObj.tree = this.processFileTreeForStage(resources, stage.stageName);
        this.fileTree.push(fileTreeObj);
      });
      this.treeLoading = false;
    });
  }

  public processFileTreeForStage(resources: ServiceResource[], stageName: string): TreeEntry[] {
    const resourcesForStage = this.getResourcesForStage(resources, stageName);
    return this._getTreeForStage(resourcesForStage);
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

  private _getTreeForStage(resources: ServiceResource[]): TreeEntry[] {
    let tree: TreeEntry[] = [];
    const root = {fileName: '', children: []};

    for (const resource of resources) {
      const uriParts = resource.resourceURI.substring(1, resource.resourceURI.length).split('/');
      const subTree = this._getSubTree([...uriParts]);

      if (tree.length === 0) {
        tree = [...subTree];
      } else {
        this._mergeTrees(tree, subTree[0], root);
        if (root.children.length !== 0) {
          tree = [...tree, ...root.children];
        }
      }
    }
    return tree;
  }

  private _getSubTree(uriParts: string[]): TreeEntry[] {
    const tree = [];
    if (uriParts.length === 1) {
      tree.push({fileName: uriParts[0]});
    } else {
      const entry: TreeEntry = {
        fileName: uriParts[0],
        children: undefined,
      };
      if (uriParts.length > 1) {
        uriParts.shift();
      }
      entry.children = this._getSubTree(uriParts);
      tree.push(entry);
    }

    return tree;
  }

  private _mergeTrees(t1: TreeEntry[], t2: TreeEntry, parent: TreeEntry): void {
    if (t1.length === 0) {
      parent?.children?.push(t2);
      return;
    }
    let found = false;
    for (const e1 of t1) {
      if (e1.fileName === t2.fileName) {
        if (e1.children && t2.children) {
          this._mergeTrees(e1.children, t2.children[0], e1);
        }
        found = true;
        break;
      }
    }
    if (!found) {
      const children = parent.children;

      const folders = children?.filter(child => child.children) || [];
      const files = children?.filter(child => !child.children) || [];

      if (t2.children) {
        folders.push(t2);
      } else {
        files.push(t2);
      }

      folders.sort((a, b) => this._compareStrings(a, b));
      files.sort((a, b) => this._compareStrings(a, b));

      parent.children = [...folders, ...files];
    }
  }

  private _compareStrings(a: TreeEntry, b: TreeEntry): number {
    if (a.fileName === b.fileName) {
      return 0;
    }
    return a.fileName < b.fileName ? -1 : 1;
  }

  public ngOnDestroy(): void {
    this.unsubscribe$.next();
    this.unsubscribe$.complete();
  }
}

export interface FileTree {
  stageName: string;
  tree: TreeEntry[];
}
