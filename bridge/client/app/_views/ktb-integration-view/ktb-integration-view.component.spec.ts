import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbIntegrationViewComponent } from './ktb-integration-view.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { of } from 'rxjs';
import { ApiService } from '../../_services/api.service';
import { ApiServiceMock } from '../../_services/api.service.mock';
import { RouterTestingModule } from '@angular/router/testing';
import { KtbIntegrationViewModule } from './ktb-integration-view.module';

describe('KtbKeptnServicesListComponent', () => {
  let component: KtbIntegrationViewComponent;
  let fixture: ComponentFixture<KtbIntegrationViewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbIntegrationViewModule, HttpClientTestingModule, RouterTestingModule],
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
    fixture = TestBed.createComponent(KtbIntegrationViewComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
