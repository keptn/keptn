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

describe('KtbEditServiceComponent', () => {
  let component: KtbEditServiceComponent;
  let fixture: ComponentFixture<KtbEditServiceComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
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
});
