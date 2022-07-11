import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbEditServiceComponent } from './ktb-edit-service.component';
import { DataService } from '../../_services/data.service';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { BehaviorSubject, firstValueFrom, throwError } from 'rxjs';
import { EventService } from '../../_services/event.service';
import { DeleteResult, DeleteType } from '../../_interfaces/delete';
import { HttpErrorResponse } from '@angular/common/http';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { FileTreeMock } from '../../_services/_mockData/fileTree.mock';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbEditServiceModule } from './ktb-edit-service.module';
import { RouterTestingModule } from '@angular/router/testing';

const paramMapSubject = new BehaviorSubject(
  convertToParamMap({
    serviceName: 'carts',
    projectName: 'sockshop',
  })
);

describe('KtbEditServiceComponent', () => {
  let component: KtbEditServiceComponent;
  let fixture: ComponentFixture<KtbEditServiceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbEditServiceModule, HttpClientTestingModule, RouterTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
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
    eventService.deletionTriggeredEvent.next({ type: DeleteType.SERVICE, name: 'carts' });
    fixture.detectChanges();

    // then
    expect(deleteProgressSpy).toHaveBeenCalledWith({ isInProgress: true });
    expect(deleteProgressSpy).toHaveBeenCalledWith({ isInProgress: false, result: DeleteResult.SUCCESS });
  });

  it('should show error', () => {
    // given
    const eventService = TestBed.inject(EventService);
    const dataService = TestBed.inject(DataService);
    const deleteProgressSpy = jest.spyOn(eventService.deletionProgressEvent, 'next');
    dataService.deleteService = jest
      .fn()
      .mockReturnValue(throwError(() => new HttpErrorResponse({ error: 'service could not be deleted' })));

    // when
    eventService.deletionTriggeredEvent.next({ type: DeleteType.SERVICE, name: 'carts' });
    fixture.detectChanges();

    // then
    expect(deleteProgressSpy).toHaveBeenCalledWith({
      isInProgress: false,
      result: DeleteResult.ERROR,
      error: 'service could not be deleted',
    });
  });

  it('should get the file tree of all stages for project sockshop and service carts', async () => {
    // given, when
    const expectedTree = FileTreeMock;
    const fileTree = await firstValueFrom(component.fileTree$);
    const dataService = TestBed.inject(DataService);
    const spy = jest.spyOn(dataService, 'getFileTreeForService');

    // when
    paramMapSubject.next(
      convertToParamMap({
        serviceName: 'carts',
        projectName: 'sockshop',
      })
    );

    // then
    expect(spy).toHaveBeenCalledWith('sockshop', 'carts');
    expect(fileTree).toBeTruthy();
    expect(fileTree).toEqual(expectedTree);
  });
});
