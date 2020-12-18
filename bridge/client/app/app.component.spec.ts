import {TestBed, async, ComponentFixture, fakeAsync, tick} from '@angular/core/testing';
import { AppComponent } from './app.component';
import {By} from "@angular/platform-browser";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {Location} from '@angular/common';
import {DataService} from "./_services/data.service";
import {MockDataService} from "./_services/mock-data.service";
import {Router} from "@angular/router";
import {AppModule} from "./app.module";
import {RouterTestingModule} from "@angular/router/testing";
import {routes} from "./app.routing";

describe('AppComponent', () => {
  let router: Router;
  let location: Location;
  let comp: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
      ],
      imports: [
        AppModule,
        HttpClientTestingModule,
        RouterTestingModule.withRoutes(routes),
      ],
      providers: [
        {provide: DataService, useClass: MockDataService}
      ]
    }).compileComponents().then(() => {
      router = TestBed.get(Router);
      location = TestBed.get(Location);
      fixture = TestBed.createComponent(AppComponent);
      comp = fixture.componentInstance;

      router.initialNavigation();
    });
  }));

  it('should create the app', () => {
    expect(comp).toBeTruthy();
  });

  it('should render title', async(() => {
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('.brand p').textContent).toContain('keptn');
  }));

  it('should render project "sockshop"', async(() => {
    fixture.detectChanges();
    const projectTileTitle = fixture.debugElement.query(By.css('#sockshop .dt-tile-title'));
    expect(projectTileTitle.nativeElement.textContent).toContain('sockshop');
  }));

  it('should set base href correctly', async(() => {
    fixture.detectChanges();

    // NOTE: function used in index.html, this is a duplicate only for testing
    function getBridgeBaseHref(origin, path) {
      if(path.indexOf('/bridge') != -1)
        return [origin, path.substring(0, path.indexOf('/bridge')), '/bridge/'].join('');
      else
        return origin;
    }

    // base = 'http://localhost:8000/'
    expect(getBridgeBaseHref('http://localhost:8000', '/dashboard')).toBe('http://localhost:8000');
    expect(getBridgeBaseHref('http://localhost:8000', '/project/sockshop')).toBe('http://localhost:8000');

    // base = 'http://localhost:8000/bridge/'
    expect(getBridgeBaseHref('http://localhost:8000', '/bridge/dashboard')).toBe('http://localhost:8000/bridge/');
    expect(getBridgeBaseHref('http://localhost:8000', '/bridge/project/sockshop')).toBe('http://localhost:8000/bridge/');

    // base 'http://0.0.0.1.xip.io/bridge/'
    expect(getBridgeBaseHref('http://0.0.0.1.xip.io', '/bridge/dashboard')).toBe('http://0.0.0.1.xip.io/bridge/');
    expect(getBridgeBaseHref('http://0.0.0.1.xip.io', '/bridge/project/sockshop')).toBe('http://0.0.0.1.xip.io/bridge/');

    // base = 'https://demo.keptn.sh/bridge/'
    expect(getBridgeBaseHref('https://demo.keptn.sh', '/bridge/dashboard')).toBe('https://demo.keptn.sh/bridge/');
    expect(getBridgeBaseHref('https://demo.keptn.sh', '/bridge/project/sockshop')).toBe('https://demo.keptn.sh/bridge/');

    // base = 'https://demo.io/keptn/bridge/'
    expect(getBridgeBaseHref('https://demo.io', '/keptn/bridge/dashboard')).toBe('https://demo.io/keptn/bridge/');
    expect(getBridgeBaseHref('https://demo.io', '/keptn/bridge/project/sockshop')).toBe('https://demo.io/keptn/bridge/');

    // base = 'https://bridge.demo.keptn.sh'
    expect(getBridgeBaseHref('https://bridge.demo.keptn.sh', '/dashboard')).toBe('https://bridge.demo.keptn.sh');
    expect(getBridgeBaseHref('https://bridge.demo.keptn.sh', '/project/sockshop')).toBe('https://bridge.demo.keptn.sh');

  }));

  xit('should render project board for "sockshop"', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);
    tick();
    expect(location.path()).toBe('/project/sockshop');
  }));
});
