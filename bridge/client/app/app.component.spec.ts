import {
  TestBed,
  async,
  ComponentFixture,
  fakeAsync,
  tick,
  discardPeriodicTasks,
  flush,
  flushMicrotasks
} from '@angular/core/testing';
import {AppComponent} from './app.component';
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
  let mockDataService: MockDataService;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [],
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
      mockDataService = TestBed.inject(MockDataService);
      comp = fixture.componentInstance;

      router.initialNavigation();
    });
  }));

  afterEach(fakeAsync(() => {
    fixture.destroy();
    TestBed.resetTestingModule();
  }));

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

  it('should set base href correctly', async(() => {
    fixture.detectChanges();

    // NOTE: function used in index.html, this is a duplicate only for testing
    function getBridgeBaseHref(origin, path) {
      if (path.indexOf('/bridge') != -1)
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

  it('should render project board for "sockshop"', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/service', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/service');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service', 'carts']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/service/carts');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service', 'carts', 'context', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/service/carts/context/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/service/:serviceName/context/:shkeptncontext/stage/:stage', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'service', 'carts', 'context', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'stage', 'staging']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/service/carts/context/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext/stage/:stage', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'stage', 'dev']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/dev');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink project/:projectName/sequence/:shkeptncontext/event/:eventId', fakeAsync(() => {
    router.navigate(['project', 'sockshop', 'sequence', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'event', 'e8f12220-b0f7-4e2f-898a-b6b7e699f12a']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/event/e8f12220-b0f7-4e2f-898a-b6b7e699f12a');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432');

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/staging');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext/:stage', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'dev']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/dev');

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/stage/dev');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

  xit('deepLink trace/:shkeptncontext/:eventtype', fakeAsync(() => {
    router.navigate(['trace', '6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432', 'sh.keptn.event.evaluation.triggered']);

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/trace/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/sh.keptn.event.evaluation.triggered');

    tick();
    fixture.detectChanges();

    expect(location.path()).toBe('/project/sockshop/sequence/6f1327d2-ded2-48ab-a1c6-e4f3d0ebe432/event/7c105021-3a50-47c7-aaa9-2e6286b17d89');

    const menuItems = fixture.debugElement.queryAll(By.css('.dt-menu .dt-menu-item'));
    expect(menuItems.length).toBe(4);
    expect(menuItems[0].nativeElement.textContent).toContain('Environment');
    expect(menuItems[0].nativeElement).not.toHaveClass('active');
    expect(menuItems[1].nativeElement.textContent).toContain('Services');
    expect(menuItems[1].nativeElement).not.toHaveClass('active');
    expect(menuItems[2].nativeElement.textContent).toContain('Sequences');
    expect(menuItems[2].nativeElement).toHaveClass('active');

    router.navigate(['/']);
    tick();
    fixture.detectChanges();
    discardPeriodicTasks(); // fixes "x timer(s) still in the queue"; TODO: check if that message means that subscriptions are not correctly unsubscribed?

    expect(location.path()).toBe('/dashboard');
  }));

});
