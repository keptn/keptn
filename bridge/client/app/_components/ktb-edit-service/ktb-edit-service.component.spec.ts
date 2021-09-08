import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { AppModule } from '../../app.module';
import { DataServiceMock } from '../../_services/data.service.mock';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of, throwError } from 'rxjs';
import { EventService } from '../../_services/event.service';
import { DeleteResult, DeleteType } from '../../_interfaces/delete';
import { HttpErrorResponse } from '@angular/common/http';
import { ServiceResourceMock } from '../../_models/serviceResource.mock';
import { ProjectMock } from '../../_models/project.mock';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('KtbEditServiceComponent', () => {
  let component: KtbEditServiceComponent;
  let fixture: ComponentFixture<KtbEditServiceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute, useValue: {
            paramMap: of(convertToParamMap({
              serviceName: 'carts',
              projectName: 'sockshop',
            })),
            snapshot: {},
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbEditServiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should delete service', () => {
    // given
    const eventService = TestBed.inject(EventService);
    const deleteProgressSpy = jest.spyOn(eventService.deletionProgressEvent, 'next');

    // when
    eventService.deletionTriggeredEvent.next({type: DeleteType.SERVICE, name: 'carts'});
    fixture.detectChanges();

    // then
    expect(deleteProgressSpy).toHaveBeenCalledWith({isInProgress: true});
    expect(deleteProgressSpy).toHaveBeenCalledWith({isInProgress: false, result: DeleteResult.SUCCESS});
  });

  it('should show error', () => {
    // given
    const eventService = TestBed.inject(EventService);
    const dataService = TestBed.inject(DataService);
    const deleteProgressSpy = jest.spyOn(eventService.deletionProgressEvent, 'next');
    dataService.deleteService = jest.fn().mockReturnValue(throwError(new HttpErrorResponse({error: 'service could not be deleted'})));

    // when
    eventService.deletionTriggeredEvent.next({type: DeleteType.SERVICE, name: 'carts'});
    fixture.detectChanges();

    // then
    expect(deleteProgressSpy).toHaveBeenCalledWith({isInProgress: false, result: DeleteResult.ERROR, error: 'service could not be deleted'});
  });

  it('should get service resources for project sockshop and service carts', () => {
    // given, when
    const project = ProjectMock;
    const dataService = TestBed.inject(DataService);
    const spy = jest.spyOn(dataService, 'getServiceResourceForAllStages');

    // when
    component.getResourcesAndTransform(project, 'carts');

    // then
    expect(spy).toHaveBeenCalledWith('sockshop', 'carts');
    expect(component.fileTree).toBeTruthy();
    expect(component.fileTree?.length).toBeGreaterThan(0);
  });

  it('should get resources for a given stage', () => {
    // given, when
    const resourcesForDev = component.getResourcesForStage(ServiceResourceMock, 'dev');

    // then
    expect(resourcesForDev).toBeTruthy();
    expect(resourcesForDev.length).toEqual(5);
    resourcesForDev.forEach((resource) => {
      expect(resource.stageName).toEqual('dev');
    });
  });

  it('should return a transformed fileTree for a given stage', () => {
    // given
    const resources = ServiceResourceMock;
    const expectedTree = [{
      files: [
        'Chart.yaml',
        'values.yaml',
      ],
      folder:
        '/helm/carts',
    }, {
      files: [
        'deployment.yaml',
        'service.yaml',
      ],
      folder:
        '/helm/carts/templates',
    }, {
      files: [
        'metadata.yaml',
      ],
      folder:
        '',
    }];

    // when
    const fileTree = component.getFileTreeForStage(resources, 'dev');

    // then
    expect(fileTree).toBeTruthy();
    expect(fileTree).toEqual(expectedTree);
  });

  it('should get the link for github for a given stage', () => {
    // given, when
    const link = component.getLinkForStage('https://github.com/keptn/sockshop-upstream', 'dev');

    // then
    expect(link).toEqual('https://github.com/keptn/sockshop-upstream/tree/dev/carts');
  });

  it('should get the link for bitbucket for a stage', () => {
    // given, when
    const link = component.getLinkForStage('https://bitbucket.org/keptn/sockshop-upstream', 'dev');

    // then
    expect(link).toEqual('https://bitbucket.org/keptn/sockshop-upstream/src/dev/carts');
  });

  it('should get the link for azure for a stage', () => {
    // given, when
    const link = component.getLinkForStage('https://dev.azure.com/keptn/_git/sockshop-upstream', 'dev');

    // then
    expect(link).toEqual('https://dev.azure.com/keptn/_git/sockshop-upstream?path=carts&version=GBdev');
  });

  it('should get the link for codeCommit for a stage', () => {
    // given, when
    const link = component.getLinkForStage('https://git-codecommit.eu-central-1.amazonaws.com/v1/repos/sockshop-upstream', 'dev');

    // then
    expect(link).toEqual('https://eu-central-1.console.aws.amazon.com/codesuite/codecommit/repositories/sockshop-upstream/browse/refs/heads/dev');
  });

  it('should return the repository url when not github, bitbucket, azure or codeCommit', () => {
    // given, when
    const link = component.getLinkForStage('https://some-other-git-provider.com/keptn/keptn-upstream', 'dev');

    // then
    expect(link).toEqual('https://some-other-git-provider.com/keptn/keptn-upstream');
  });

  it('should return an empty string if no remote URI is set', () => {
    // given, when
    const link = component.getLinkForStage('', 'dev');

    // then
    expect(link).toEqual('');
  });

  it('should a note that the Git upstream has to be set if the remoteURI is not set', () => {
    // given, when
    const projectMock = ProjectMock;
    projectMock.gitRemoteURI = '';
    projectMock.gitUser = '';
    component.project$ = of(projectMock);
    component.serviceName = 'carts';
    fixture.detectChanges();

    const elem = fixture.nativeElement.querySelector('.settings-section:first-of-type span');

    // then
    expect(elem).toBeTruthy();
    expect(elem.textContent).toContain('There is no Git upstream repository set.');
  });

  it('should show the files for all stages when the remoteURI is set', () => {
    // given
    const projectMock = ProjectMock;
    projectMock.gitRemoteURI = 'https://github.com/keptn/sockshop-upstream';
    projectMock.gitUser = 'keptn';
    component.project$ = of(projectMock);

    // when
    component.getResourcesAndTransform(projectMock, 'carts');
    fixture.detectChanges();

    // then
    const stages = fixture.nativeElement.querySelectorAll('.settings-section:first-of-type > div > div.bold');
    const folders = fixture.nativeElement.querySelectorAll('.settings-section:first-of-type .dt-expandable-section-header-content');

    expect(stages).toBeTruthy();
    expect(stages.length).toEqual(3);
    expect(stages[0].textContent).toEqual('development');
    expect(stages[1].textContent).toEqual('staging');
    expect(stages[2].textContent).toEqual('production');
    expect(folders).toBeTruthy();
    expect(folders.length).toEqual(4);
    expect(folders[0].textContent).toEqual('/helm/carts');
    expect(folders[1].textContent).toEqual('/helm/carts/templates');
    expect(folders[2].textContent).toEqual('/helm/carts');
    expect(folders[3].textContent).toEqual('/helm/carts/templates');
  });
});
