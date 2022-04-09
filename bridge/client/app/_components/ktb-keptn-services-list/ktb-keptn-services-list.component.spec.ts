import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbKeptnServicesListComponent } from './ktb-keptn-services-list.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { APIService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbKeptnServicesListComponent;
  let fixture: ComponentFixture<KtbKeptnServicesListComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [
        { provide: APIService, useClass: ApiServiceMock },
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
