import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppHeaderComponent } from './app-header.component';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { AppModule } from '../app.module';
import { RETRY_ON_HTTP_ERROR } from '../_utils/app.utils';
import { KeptnInfo } from '../_models/keptn-info';

describe('AppHeaderComponent', () => {
  let component: AppHeaderComponent;
  let fixture: ComponentFixture<AppHeaderComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule],
      providers: [{ provide: RETRY_ON_HTTP_ERROR, useValue: false }],
    }).compileComponents();

    fixture = TestBed.createComponent(AppHeaderComponent);
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
