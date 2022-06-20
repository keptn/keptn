import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { KtbKeptnServicesListModule } from './ktb-keptn-services-list.module';

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbKeptnServicesListModule, HttpClientTestingModule],
      providers: [
        { provide: ApiService, useClass: ApiServiceMock },
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(
              convertToParamMap({
                projectName: 'sockshop',
              })
            ),
          },
        },
      ],
    }).compileComponents();

    localStorage.setItem('keptn_integration_dates', '');
    fixture = TestBed.createComponent(KtbKeptnServicesListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
