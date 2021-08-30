import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbCreateServiceComponent } from './ktb-create-service.component';
import { DataService } from '../../_services/data.service';
import { DataServiceMock } from '../../_services/data.service.mock';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { AppModule } from '../../app.module';

describe('KtbCreateServiceComponent', () => {
  let component: KtbCreateServiceComponent;
  let fixture: ComponentFixture<KtbCreateServiceComponent>;
  const projectName = 'sockshop';

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        AppModule,
      ],
      providers: [
        {provide: DataService, useClass: DataServiceMock},
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({
              projectName,
            })),
          },
        },
      ],
    }).compileComponents();
    fixture = TestBed.createComponent(KtbCreateServiceComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should show duplicate error', () => {
    const serviceNames = ['carts', 'carts-db'];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
      expect(component.formGroup.hasError('duplicate'));
    }
  });

  it('should show pattern error', () => {
    const serviceNames = ['Service', '1service', '-service', '$service', 'serVice', 'ser_ice'];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
      expect(component.formGroup.hasError('pattern'));
    }
  });

  it('should show required error', () => {
    const serviceNames = ['service', ''];
    for (const serviceName of serviceNames) {
      component.serviceNameControl.setValue(serviceName);
      component.formGroup.updateValueAndValidity();
    }
    expect(component.formGroup.hasError('required'));
    checkCreateButton(false);
  });

  it('should create service', () => {
    // given
    const serviceName = 'service-1';
    component.serviceNameControl.setValue(serviceName);
    component.formGroup.updateValueAndValidity();
    fixture.detectChanges();
    const dataService = TestBed.inject(DataService);
    const createServiceSpy = jest.spyOn(dataService, 'createService');
    const createButton = getCreateButton();

    expect(component.formGroup.errors).toBeNull();
    checkCreateButton(true);

    // when
    createButton.click();
    fixture.detectChanges();

    // then
    expect(createServiceSpy).toHaveBeenCalledWith(projectName, serviceName);
  });

  function checkCreateButton(isEnabled: boolean): void {
    const disabled: null | string = getCreateButton().getAttribute('disabled');
    if (isEnabled) {
      expect(disabled).toBeNull();
    } else {
      expect(disabled).not.toBeNull();
    }
  }

  function getCreateButton(): HTMLElement {
    return fixture.nativeElement.querySelector('button[uitestid=createServiceButton]');
  }
});
