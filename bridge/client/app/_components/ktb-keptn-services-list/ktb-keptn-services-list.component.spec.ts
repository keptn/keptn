import { ComponentFixture, fakeAsync, TestBed, waitForAsync } from '@angular/core/testing';
import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { UniformRegistrationsMock } from '../../_models/uniform-registrations.mock';
import { of } from 'rxjs';
import { UniformRegistrationLogsMock } from '../../_models/uniform-registrations-logs.mock';

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

  beforeEach(waitForAsync(() => {
    TestBed.configureTestingModule({
      declarations: [],
      imports: [
        AppModule,
        HttpClientTestingModule,
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({
              projectName: 'sockshop'
            }))
          }
        }
      ]
    })
      .compileComponents()
      .then(() => {
        localStorage.setItem('keptn_integration_dates', '');
        fixture = TestBed.createComponent(KtbKeptnServicesListComponent);
        component = fixture.componentInstance;
        fixture.detectChanges();
      });
  }));

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show 10 registrations', () => {
    const registrations = fixture.nativeElement.querySelectorAll('dt-row');
    expect(registrations.length).toEqual(10);
  });

  it('should show error event indicator', () => {
    // given
    const firstCell = fixture.nativeElement.querySelector('dt-cell');
    const indicator = firstCell.querySelector('.notification-indicator');

    // then
    expect(indicator).toBeTruthy();
    expect(indicator.innerText).toEqual('10');
  });

  it('should not show error event indicator', () => {
    // given
    const firstColumn = fixture.nativeElement.querySelector('dt-row dt-cell:nth-child(1)');

    // then
    for (let i = 1; i < firstColumn.length; ++i) {
      const indicator = firstColumn[i].querySelector('.notification-indicator');
      expect(indicator).toBeFalsy();
    }
  });

  it('should remove error event indicator on selection change', () => {
    // given
    const dataService = TestBed.inject(DataService);
    const spySave = spyOn(dataService, 'setUniformDate');
    const firstRow = fixture.nativeElement.querySelector('dt-row');
    const secondRow = fixture.nativeElement.querySelector('dt-row:nth-of-type(2)');
    const firstCell = firstRow.querySelector('dt-cell');

    // when
    firstRow.click();
    fixture.detectChanges();
    let indicator = firstCell.querySelector('.notification-indicator');
    const registration = component.selectedUniformRegistration;
    expect(indicator).toBeTruthy();
    expect(registration?.unreadEventsCount).toEqual(10);
    expect(spySave).toHaveBeenCalledOnceWith(UniformRegistrationsMock[0].id, UniformRegistrationLogsMock[0].time);

    secondRow.click();
    fixture.detectChanges();

    // then
    indicator = firstCell.querySelector('.notification-indicator');
    expect(indicator).toBeFalsy();
    expect(registration?.unreadEventsCount).toEqual(0);
  });

  it('should show error events list', () => {
    fixture.nativeElement.querySelector('dt-row').click();
    fixture.detectChanges();

    const logs = fixture.nativeElement.querySelector('ktb-uniform-registration-logs');
    expect(logs).toBeTruthy();
    expect(fixture.nativeElement.querySelector('h3').innerText).toEqual('ansible-service');
  });

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));
});
