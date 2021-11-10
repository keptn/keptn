import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbServiceDetailsComponent } from './ktb-service-details.component';
import { AppModule } from '../../app.module';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { ActivatedRoute, convertToParamMap } from '@angular/router';
import { Deployment } from '../../_models/deployment';
import { ServiceDeploymentMock } from '../../../../shared/fixtures/service-deployment-response.mock';
import { Location } from '@angular/common';
import { of } from 'rxjs';

describe('KtbServiceDetailsComponent', () => {
  let component: KtbServiceDetailsComponent;
  let fixture: ComponentFixture<KtbServiceDetailsComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      providers: [
        {
          provide: ActivatedRoute,
          useValue: {
            paramMap: of(convertToParamMap({ projectName: 'sockshop' })),
          },
        },
      ],
      imports: [AppModule, HttpClientTestingModule],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbServiceDetailsComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should select stage', () => {
    const location = TestBed.inject(Location);
    const locationSpy = jest.spyOn(location, 'go');
    component.deploymentInfo = {
      deployment: Deployment.fromJSON(ServiceDeploymentMock),
      stage: 'dev',
    };
    component.selectStage('production');
    expect(component.deploymentInfo.stage).toBe('production');
    expect(locationSpy).toHaveBeenCalledWith(
      `/project/sockshop/service/carts/context/${ServiceDeploymentMock.keptnContext}/stage/production`
    );
  });

  it('should be url', () => {
    expect(component.isUrl('http://keptn.sh')).toBe(true);
  });

  it('should be label', () => {
    expect(component.isUrl('myLabel')).toBe(false);
  });

  it('should show remediation dialog', () => {
    component.deploymentInfo = {
      deployment: Deployment.fromJSON(ServiceDeploymentMock),
      stage: 'production',
    };
    fixture.detectChanges();
    const remediationConfigButton = fixture.nativeElement.querySelector('[uitestid=show-remediation-config-button]');
    expect(remediationConfigButton).toBeTruthy();
    remediationConfigButton.click();
    fixture.detectChanges();
    const remediationConfigContent = document.querySelector('[uitestid=remediation-config-content]');
    expect(remediationConfigContent).toBeTruthy();
  });
});
