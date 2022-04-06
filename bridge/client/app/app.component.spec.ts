import { ComponentFixture, TestBed } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { By } from '@angular/platform-browser';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Location } from '@angular/common';
import { Router } from '@angular/router';
import { AppModule } from './app.module';
import { RouterTestingModule } from '@angular/router/testing';
import { routes } from './app.routing';
import { ApiService } from './_services/api.service';
import { ApiServiceMock } from './_services/api.service.mock';
import { DataService } from './_services/data.service';

describe('AppComponent', () => {
  let router: Router;
  let location: Location;
  let comp: AppComponent;
  let fixture: ComponentFixture<AppComponent>;
  enum MENU_ITEM {
    ENVIRONMENT,
    SERVICES,
    SEQUENCES,
    INTEGRATIONS,
    UNIFORM,
    SETTINGS,
  }

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      declarations: [],
      imports: [AppModule, HttpClientTestingModule, RouterTestingModule.withRoutes(routes)],
      providers: [{ provide: ApiService, useClass: ApiServiceMock }],
    }).compileComponents();

    router = TestBed.inject(Router);
    location = TestBed.inject(Location);
    fixture = TestBed.createComponent(AppComponent);
    comp = fixture.componentInstance;

    router.initialNavigation();
  });

  it('should create the app', () => {
    expect(comp).toBeTruthy();
  });

  it('should set base href correctly', () => {
    fixture.detectChanges();

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
    const dataService = TestBed.inject(DataService);
    const loadSpy = jest.spyOn(dataService, 'loadProjects');
    fixture.detectChanges();

    expect(loadSpy).toHaveBeenCalled();
  });

  function assertMenuItems(activeItem: MENU_ITEM): void {
    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toEqual(6);

    if (activeItem === MENU_ITEM.ENVIRONMENT) {
      expect(menuItems[0].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[0].nativeElement.getAttribute('class')).not.toContain('active');
    }

    if (activeItem === MENU_ITEM.SERVICES) {
      expect(menuItems[1].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[1].nativeElement.getAttribute('class')).not.toContain('active');
    }

    if (activeItem === MENU_ITEM.SEQUENCES) {
      expect(menuItems[2].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[2].nativeElement.getAttribute('class')).not.toContain('active');
    }

    if (activeItem === MENU_ITEM.INTEGRATIONS) {
      expect(menuItems[3].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[3].nativeElement.getAttribute('class')).not.toContain('active');
    }

    if (activeItem === MENU_ITEM.UNIFORM) {
      expect(menuItems[4].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[4].nativeElement.getAttribute('class')).not.toContain('active');
    }

    if (activeItem === MENU_ITEM.SETTINGS) {
      expect(menuItems[5].nativeElement.getAttribute('class')).toContain('active');
    } else {
      expect(menuItems[5].nativeElement.getAttribute('class')).not.toContain('active');
    }
  }
});
