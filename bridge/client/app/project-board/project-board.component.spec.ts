import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ProjectBoardComponent } from './project-board.component';
import {DtTag, DtTagList, DtTagModule} from '@dynatrace/barista-components/tag';
import {DtLoadingDistractorModule, DtLoadingSpinner} from "@dynatrace/barista-components/loading-distractor";
import {DtInfoGroup, DtInfoGroupModule, DtInfoGroupTitle} from "@dynatrace/barista-components/info-group";
import {
  DtEmptyState,
  DtEmptyStateItem,
  DtEmptyStateItemImage,
  DtEmptyStateItemTitle, DtEmptyStateModule
} from "@dynatrace/barista-components/empty-state";
import {RouterModule} from "@angular/router";
import {AppComponent} from "../app.component";
import {DashboardComponent} from "../dashboard/dashboard.component";
import {AppHeaderComponent} from "../app-header/app-header.component";
import {KtbHttpLoadingSpinnerComponent} from "../_components/ktb-http-loading-spinner/ktb-http-loading-spinner.component";
import {KtbHttpLoadingBarComponent} from "../_components/ktb-http-loading-bar/ktb-http-loading-bar.component";
import {KtbShowHttpLoadingDirective} from "../_directives/ktb-show-http-loading/ktb-show-http-loading.directive";
import {KtbHideHttpLoadingDirective} from "../_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive";
import {
  KtbExpandableTileComponent,
  KtbExpandableTileHeader
} from "../_components/ktb-expandable-tile/ktb-expandable-tile.component";
import {KtbSelectableTileComponent} from "../_components/ktb-selectable-tile/ktb-selectable-tile.component";
import {
  KtbHorizontalSeparatorComponent,
  KtbHorizontalSeparatorTitle
} from "../_components/ktb-horizontal-separator/ktb-horizontal-separator.component";
import {KtbRootEventsListComponent} from "../_components/ktb-root-events-list/ktb-root-events-list.component";
import {KtbProjectTileComponent} from "../_components/ktb-project-tile/ktb-project-tile.component";
import {KtbProjectListComponent} from "../_components/ktb-project-list/ktb-project-list.component";
import {KtbEventsListComponent} from "../_components/ktb-events-list/ktb-events-list.component";
import {AtobPipe} from "../_pipes/atob.pipe";
import {KtbEventItemComponent, KtbEventItemDetail} from "../_components/ktb-event-item/ktb-event-item.component";
import {KtbEvaluationDetailsComponent} from "../_components/ktb-evaluation-details/ktb-evaluation-details.component";
import {KtbSliBreakdownComponent} from "../_components/ktb-sli-breakdown/ktb-sli-breakdown.component";
import {BrowserModule} from "@angular/platform-browser";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";
import {HttpClientTestingModule} from "@angular/common/http/testing";
import {AppRouting} from "../app.routing";
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
import {DtCardModule} from "@dynatrace/barista-components/card";
import {DtTileModule} from "@dynatrace/barista-components/tile";
import {DtProgressBarModule} from "@dynatrace/barista-components/progress-bar";
import {DtExpandableTextModule} from "@dynatrace/barista-components/expandable-text";
import {DtExpandablePanelModule} from "@dynatrace/barista-components/expandable-panel";
import {DtShowMoreModule} from "@dynatrace/barista-components/show-more";
import {DtIndicatorModule} from "@dynatrace/barista-components/core";
import {DtProgressCircleModule} from "@dynatrace/barista-components/progress-circle";
import {DtConsumptionModule} from "@dynatrace/barista-components/consumption";
import {DtKeyValueListModule} from "@dynatrace/barista-components/key-value-list";
import {DtChartModule} from "@dynatrace/barista-components/chart";
import {DtIconModule} from "@dynatrace/barista-components/icon";

describe('ProjectBoardComponent', () => {
  let component: ProjectBoardComponent;
  let fixture: ComponentFixture<ProjectBoardComponent>;

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
        BrowserAnimationsModule
      ],
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ProjectBoardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
