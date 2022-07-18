import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceSettingsOverviewComponent } from './ktb-service-settings-overview.component';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { firstValueFrom, of } from 'rxjs';
import { By } from '@angular/platform-browser';
import { KtbServiceSettingsModule } from './ktb-service-settings.module';
import { RouterTestingModule } from '@angular/router/testing';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ApiServiceMock } from '../../../_services/api.service.mock';
import { ApiService } from '../../../_services/api.service';
import { take } from 'rxjs/operators';

describe('KtbServiceSettingsOverviewComponent', () => {
  let component: KtbServiceSettingsOverviewComponent;
  let fixture: ComponentFixture<KtbServiceSettingsOverviewComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [KtbServiceSettingsModule, RouterTestingModule, HttpClientTestingModule],
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
    fixture = TestBed.createComponent(KtbServiceSettingsOverviewComponent);
    component = fixture.componentInstance;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set the projectName', (done) => {
    component.projectName$.pipe(take(1)).subscribe((projectName) => {
      expect(projectName).toBe('sockshop');
      done();
    });
  });

  it('should load serviceNames for the project', (done) => {
    component.serviceNames$.pipe(take(1)).subscribe((serviceNames) => {
      expect(serviceNames).toEqual(['carts', 'carts-db']);
      done();
    });
  });

  it('should show create button', async () => {
    // given
    await firstValueFrom(component.serviceNames$);
    fixture.detectChanges();

    // when, then
    const createButton = fixture.debugElement.query(By.css('button'));
    expect(createButton).toBeTruthy();
    expect(createButton.nativeElement.textContent.trim()).toEqual('Create service');
  });
});
