import {TestBed, async, ComponentFixture, fakeAsync, tick} from '@angular/core/testing';
import { AppComponent } from './app.component';
import {AppHeaderComponent} from "./app-header/app-header.component";
import {KtbHttpLoadingBarComponent} from "./_components/ktb-http-loading-bar/ktb-http-loading-bar.component";
import {DashboardComponent} from "./dashboard/dashboard.component";
import {ProjectBoardComponent} from "./project-board/project-board.component";
import {KtbHttpLoadingSpinnerComponent} from "./_components/ktb-http-loading-spinner/ktb-http-loading-spinner.component";
import {KtbShowHttpLoadingDirective} from "./_directives/ktb-show-http-loading/ktb-show-http-loading.directive";
import {KtbHideHttpLoadingDirective} from "./_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive";
import {
  KtbExpandableTileComponent,
  KtbExpandableTileHeader
} from "./_components/ktb-expandable-tile/ktb-expandable-tile.component";
import {KtbSelectableTileComponent} from "./_components/ktb-selectable-tile/ktb-selectable-tile.component";
import {
  KtbHorizontalSeparatorComponent,
  KtbHorizontalSeparatorTitle
} from "./_components/ktb-horizontal-separator/ktb-horizontal-separator.component";
import {KtbRootEventsListComponent} from "./_components/ktb-root-events-list/ktb-root-events-list.component";
import {KtbProjectTileComponent} from "./_components/ktb-project-tile/ktb-project-tile.component";
import {KtbProjectListComponent} from "./_components/ktb-project-list/ktb-project-list.component";
import {KtbEventsListComponent} from "./_components/ktb-events-list/ktb-events-list.component";
import {AtobPipe} from "./_pipes/atob.pipe";
import {KtbEventItemComponent, KtbEventItemDetail} from "./_components/ktb-event-item/ktb-event-item.component";
import {KtbEvaluationDetailsComponent} from "./_components/ktb-evaluation-details/ktb-evaluation-details.component";
import {KtbSliBreakdownComponent} from "./_components/ktb-sli-breakdown/ktb-sli-breakdown.component";
import {BrowserModule, By} from "@angular/platform-browser";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {Location} from '@angular/common';
import {AppRouting, routes} from "./app.routing";
import {FlexLayoutModule} from "@angular/flex-layout";
import {MomentModule} from "ngx-moment";
import {DtThemingModule} from "@dynatrace/barista-components/theming";
import {DtButtonModule} from "@dynatrace/barista-components/button";
import {DtButtonGroupModule} from "@dynatrace/barista-components/button-group";
import {DtSelectModule} from "@dynatrace/barista-components/select";
import {DtMenuModule} from "@dynatrace/barista-components/menu";
import {DtDrawerModule} from "@dynatrace/barista-components/drawer";
import {DtContextDialogModule} from "@dynatrace/barista-components/context-dialog";
import {DtInputModule} from "@dynatrace/barista-components/input";
import {DtEmptyStateModule} from "@dynatrace/barista-components/empty-state";
import {DtCardModule} from "@dynatrace/barista-components/card";
import {DtTileModule} from "@dynatrace/barista-components/tile";
import {DtInfoGroupModule} from "@dynatrace/barista-components/info-group";
import {DtProgressBarModule} from "@dynatrace/barista-components/progress-bar";
import {DtLoadingDistractorModule} from "@dynatrace/barista-components/loading-distractor";
import {DtTagModule} from "@dynatrace/barista-components/tag";
import {DtExpandableTextModule} from "@dynatrace/barista-components/expandable-text";
import {DtExpandablePanelModule} from "@dynatrace/barista-components/expandable-panel";
import {DtShowMoreModule} from "@dynatrace/barista-components/show-more";
import {DtIndicatorModule} from "@dynatrace/barista-components/core";
import {DtProgressCircleModule} from "@dynatrace/barista-components/progress-circle";
import {DtConsumptionModule} from "@dynatrace/barista-components/consumption";
import {DtKeyValueListModule} from "@dynatrace/barista-components/key-value-list";
import {DtChartModule} from "@dynatrace/barista-components/chart";
import {DtIconModule} from "@dynatrace/barista-components/icon";
import {DataService} from "./_services/data.service";
import {MockDataService} from "./_services/mock-data.service";
import {Router} from "@angular/router";
import {RouterTestingModule} from "@angular/router/testing";

describe('AppComponent', () => {
  let router: Router;
  let location: Location;
  let comp: AppComponent;
  let fixture: ComponentFixture<AppComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [
        AppComponent,
        DashboardComponent,
        AppHeaderComponent,
        ProjectBoardComponent,
        KtbHttpLoadingSpinnerComponent,
        KtbHttpLoadingBarComponent,
        KtbShowHttpLoadingDirective,
        KtbHideHttpLoadingDirective,
        KtbExpandableTileComponent,
        KtbExpandableTileHeader,
        KtbSelectableTileComponent,
        KtbHorizontalSeparatorComponent,
        KtbHorizontalSeparatorTitle,
        KtbRootEventsListComponent,
        KtbProjectTileComponent,
        KtbProjectListComponent,
        KtbEventsListComponent,
        AtobPipe,
        KtbEventItemComponent,
        KtbEventItemDetail,
        KtbEvaluationDetailsComponent,
        KtbSliBreakdownComponent,
      ],
      imports: [
        BrowserModule,
        BrowserAnimationsModule,
        HttpClientTestingModule,
        AppRouting,
        FlexLayoutModule,
        MomentModule,
        DtThemingModule,
        DtButtonModule,
        DtButtonGroupModule,
        DtSelectModule,
        DtMenuModule,
        DtDrawerModule,
        DtContextDialogModule,
        DtInputModule,
        DtEmptyStateModule,
        DtCardModule,
        DtTileModule,
        DtInfoGroupModule,
        DtProgressBarModule,
        DtLoadingDistractorModule,
        DtTagModule,
        DtExpandableTextModule,
        DtExpandablePanelModule,
        DtShowMoreModule,
        DtIndicatorModule,
        DtProgressCircleModule,
        DtConsumptionModule,
        DtKeyValueListModule,
        DtChartModule,
        DtIconModule.forRoot({
          svgIconLocation: `/assets/icons/{{name}}.svg`,
        }),
        RouterTestingModule.withRoutes(routes)
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
