import { TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { ApiService } from './_services/api.service';
import { ApiServiceMock } from './_services/api.service.mock';
import { DataService } from './_services/data.service';
import { HttpClient } from '@angular/common/http';
import { HttpClientTestingModule } from '@angular/common/http/testing';

describe('AppComponent', () => {
  let component: AppComponent;
  let dataService: DataService;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [HttpClientTestingModule],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    dataService = TestBed.inject(DataService);
  });

  it('should create the app', () => {
    createComponent();
    expect(component).toBeTruthy();
  });

  it('should set base href correctly', () => {
    // NOTE: function used in index.html, this is a duplicate only for testing
    function getBridgeBaseHref(origin: string, path: string): string {
      if (path.indexOf('/bridge') !== -1) {
        return [origin, path.substring(0, path.indexOf('/bridge')), '/bridge/'].join('');
      } else {
        return origin;
      }
    }

    // base = 'http://localhost:8000/'
    expect(getBridgeBaseHref('http://localhost:8000', '/dashboard')).toEqual('http://localhost:8000');
    expect(getBridgeBaseHref('http://localhost:8000', '/project/sockshop')).toEqual('http://localhost:8000');

    // base = 'http://localhost:8000/bridge/'
    expect(getBridgeBaseHref('http://localhost:8000', '/bridge/dashboard')).toEqual('http://localhost:8000/bridge/');
    expect(getBridgeBaseHref('http://localhost:8000', '/bridge/project/sockshop')).toEqual(
      'http://localhost:8000/bridge/'
    );

    // base 'http://0.0.0.1.xip.io/bridge/'
    expect(getBridgeBaseHref('http://0.0.0.1.xip.io', '/bridge/dashboard')).toEqual('http://0.0.0.1.xip.io/bridge/');
    expect(getBridgeBaseHref('http://0.0.0.1.xip.io', '/bridge/project/sockshop')).toEqual(
      'http://0.0.0.1.xip.io/bridge/'
    );

    // base = 'https://demo.keptn.sh/bridge/'
    expect(getBridgeBaseHref('https://demo.keptn.sh', '/bridge/dashboard')).toEqual('https://demo.keptn.sh/bridge/');
    expect(getBridgeBaseHref('https://demo.keptn.sh', '/bridge/project/sockshop')).toEqual(
      'https://demo.keptn.sh/bridge/'
    );

    // base = 'https://demo.io/keptn/bridge/'
    expect(getBridgeBaseHref('https://demo.io', '/keptn/bridge/dashboard')).toEqual('https://demo.io/keptn/bridge/');
    expect(getBridgeBaseHref('https://demo.io', '/keptn/bridge/project/sockshop')).toEqual(
      'https://demo.io/keptn/bridge/'
    );

    // base = 'https://bridge.demo.keptn.sh'
    expect(getBridgeBaseHref('https://bridge.demo.keptn.sh', '/dashboard')).toEqual('https://bridge.demo.keptn.sh');
    expect(getBridgeBaseHref('https://bridge.demo.keptn.sh', '/project/sockshop')).toEqual(
      'https://bridge.demo.keptn.sh'
    );
  });

  it('should load projects after info is loaded', () => {
    // given, when
    const loadSpy = jest.spyOn(dataService, 'loadProjects');
    createComponent();

    // then
    expect(loadSpy).toHaveBeenCalled();
  });

  function createComponent(): void {
    component = new AppComponent(TestBed.inject(HttpClient), dataService);
  }
});
