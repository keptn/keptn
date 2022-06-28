import { ComponentFixture, TestBed } from '@angular/core/testing';
import { KtbAppHeaderComponent } from './ktb-app-header.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { RETRY_ON_HTTP_ERROR } from '../../_utils/app.utils';
import { KeptnInfo } from '../../_models/keptn-info';
import { KtbAppHeaderModule } from './ktb-app-header.module';
import { RouterTestingModule } from '@angular/router/testing';

describe('AppHeaderComponent', () => {
  let component: KtbAppHeaderComponent;
  let fixture: ComponentFixture<KtbAppHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [KtbAppHeaderModule, HttpClientTestingModule, RouterTestingModule],
      providers: [{ provide: RETRY_ON_HTTP_ERROR, useValue: false }],
    }).compileComponents();

    fixture = TestBed.createComponent(KtbAppHeaderComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should display the AUTH_MSG if set', () => {
    // given
    const getKeptnInfo = (authMsg: string | undefined): KeptnInfo => {
      return {
        authCommand: 'authCommand',
        bridgeInfo: {
          authMsg,
          authType: '',
          cliDownloadLink: '',
          enableVersionCheckFeature: false,
          featureFlags: { D3_HEATMAP_ENABLED: false, RESOURCE_SERVICE_ENABLED: false },
          showApiToken: false,
        },
      };
    };

    // when
    const value1 = component.getKeptnAuthCommand(getKeptnInfo('Hello there'));
    const value2 = component.getKeptnAuthCommand(getKeptnInfo(''));
    const value3 = component.getKeptnAuthCommand(getKeptnInfo(undefined));

    // then
    expect(value1).toEqual('Hello there');
    expect(value2).toEqual('authCommand');
    expect(value3).toEqual('authCommand');
  });
});
