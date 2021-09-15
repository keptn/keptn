import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { AppModule } from '../../app.module';
import { DataServiceMock } from '../../_services/data.service.mock';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BehaviorSubject, of, throwError } from 'rxjs';
import { EventService } from '../../_services/event.service';
import { DeleteResult, DeleteType } from '../../_interfaces/delete';
import { HttpErrorResponse } from '@angular/common/http';
import { ProjectMock } from '../../_models/project.mock';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FileTreeMock } from '../../_models/fileTree.mock';
import { By } from '@angular/platform-browser';

const paramMapSubject = new BehaviorSubject(convertToParamMap({
  serviceName: 'carts',
  projectName: 'sockshop',
}));

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
            paramMap: paramMapSubject.asObservable(),
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

  it('should get the file tree of all stages for project sockshop and service carts', (done) => {
    // given, when
    const expectedTree = FileTreeMock;

    const dataService = TestBed.inject(DataService);
    const spy = jest.spyOn(dataService, 'getFileTreeForService');

    // when
    paramMapSubject.next(convertToParamMap({
      serviceName: 'carts',
      projectName: 'sockshop',
    }));

    // then
    expect(spy).toHaveBeenCalledWith('sockshop', 'carts');
    component.fileTree$.subscribe((fileTree) => {
      expect(fileTree).toBeTruthy();
      expect(fileTree).toEqual(expectedTree);
      done();
    });
  });

  it('should show a message when file tree is empty', (done) => {
    // given, when
    const dataService = TestBed.inject(DataService);
    jest.spyOn(dataService, 'getFileTreeForService').mockReturnValue(of([]));

    // when
    paramMapSubject.next(convertToParamMap({
      serviceName: 'carts',
      projectName: 'sockshop',
    }));

    // then
    component.fileTree$.subscribe((fileTree) => {
      fixture.detectChanges();
      const section = fixture.debugElement.query(By.css('.settings-section:first-of-type > div'));
      expect(fileTree).toBeTruthy();
      expect(fileTree).toEqual([]);
      expect(section.nativeElement.textContent.trim()).toEqual('There are no files in the Git upstream repository');
      done();
    });
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
});
