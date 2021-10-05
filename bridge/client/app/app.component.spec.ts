import { ComponentFixture, discardPeriodicTasks, fakeAsync, TestBed, tick } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { By } from '@angular/platform-browser';
import { HttpClientTestingModule } from '@angular/common/http/testing';
import { Location } from '@angular/common';
import { DataService } from './_services/data.service';
import { DataServiceMock } from './_services/data.service.mock';
import { Router } from '@angular/router';
import { AppModule } from './app.module';
import { RouterTestingModule } from '@angular/router/testing';
import { routes } from './app.routing';

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
      providers: [{ provide: DataService, useClass: DataServiceMock }],
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

  it('should render title', fakeAsync(() => {
    fixture.detectChanges();
    const compiled = fixture.debugElement.nativeElement;
    expect(compiled.querySelector('.brand p').textContent).toContain('keptn');

    tick(3000);
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?
  }));

  it('should render project "sockshop"', fakeAsync(() => {
    fixture.detectChanges();
    const projectTileTitle = fixture.debugElement.query(By.css('#sockshop .dt-tile-title'));
    expect(projectTileTitle.nativeElement.textContent).toContain('sockshop');

    tick(3000);
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?
  }));

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

  it('should render project board for "sockshop"', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);

    tick(3000);
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop');
    assertMenuItems(MENU_ITEM.ENVIRONMENT);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/service', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/service');

    assertMenuItems(MENU_ITEM.SERVICES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service', 'carts']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/service/carts');

    assertMenuItems(MENU_ITEM.SERVICES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service', 'carts', 'context', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/service/carts/context/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432');

    assertMenuItems(MENU_ITEM.SERVICES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext/stage/:stage', fakeAsync(() => {
    router.navigate([
      'project',
      'sockshop',
      'service',
      'carts',
      'context',
      '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432',
      'stage',
      'staging',
    ]);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual(
      '/project/sockshop/service/carts/context/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging'
    );

    assertMenuItems(MENU_ITEM.SERVICES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/sequence');

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging');

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext/stage/:stage', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'stage', 'dev']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/dev');

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext/event/:eventId', fakeAsync(() => {
    router.navigate([
      'project',
      'sockshop',
      'sequence',
      '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432',
      'event',
      'e8f12220-b0f7-4e2f-898a-b6b7e699f12a',
    ]);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual(
      '/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/event/e8f12220-b0f7-4e2f-898a-b6b7e699f12a'
    );

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432');

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging');

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext/:stage', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'dev']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/dev');

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/dev');

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext/:eventtype', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'sh.keptn.event.evaluation.triggered']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/sh.keptn.event.evaluation.triggered');

    tick();
    fixture.detectChanges();

    expect(location.path()).toEqual(
      '/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/event/7c105021-3a50-47c7-aaa9-2e6286b17d89'
    );

    assertMenuItems(MENU_ITEM.SEQUENCES);

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toEqual('/dashboard');
  }));

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
