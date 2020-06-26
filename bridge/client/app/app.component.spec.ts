import {TestBed, async, ComponentFixture, fakeAsync, tick} from '@angular/core/testing';
import { AppComponent } from './app.component';
import {BrowserModule, By} from "@angular/platform-browser";
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

  xit('should render project board for "sockshop"', fakeAsync(() => {
    router.navigate(['project', 'sockshop']);
    tick();
    expect(location.path()).toBe('/project/sockshop');
  }));
});
