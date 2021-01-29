import { registerLocaleData } from '@angular/common';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import localeEn from '@angular/common/locales/en';
import { NgModule } from '@angular/core';
import { FormsModule } from "@angular/forms";
import { FlexLayoutModule } from '@angular/flex-layout';
import { MatDialogModule } from '@angular/material/dialog';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { DtCardModule } from '@dynatrace/barista-components/card';
import { DtChartModule } from '@dynatrace/barista-components/chart';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtConsumptionModule } from '@dynatrace/barista-components/consumption';
import { DtContextDialogModule } from '@dynatrace/barista-components/context-dialog';
import { DtDrawerModule } from '@dynatrace/barista-components/drawer';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtExpandableTextModule } from '@dynatrace/barista-components/expandable-text';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtIndicatorModule } from '@dynatrace/barista-components/indicator';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtKeyValueListModule } from '@dynatrace/barista-components/key-value-list';
import { DtLoadingDistractorModule } from '@dynatrace/barista-components/loading-distractor';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtProgressBarModule } from '@dynatrace/barista-components/progress-bar';
import { DtProgressCircleModule } from '@dynatrace/barista-components/progress-circle';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtSwitchModule } from '@dynatrace/barista-components/switch';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtTopBarNavigationModule } from "@dynatrace/barista-components/top-bar-navigation";
import { DtCopyToClipboardModule } from "@dynatrace/barista-components/copy-to-clipboard";
import { DtToggleButtonGroupModule } from "@dynatrace/barista-components/toggle-button-group";
import { DtQuickFilterModule } from "@dynatrace/barista-components/experimental/quick-filter";

import { DtThemingModule } from '@dynatrace/barista-components/theming';
import { DtTileModule } from '@dynatrace/barista-components/tile';
import { DtToastModule } from '@dynatrace/barista-components/toast';

import { MomentModule } from 'ngx-moment';
import { KtbEvaluationDetailsComponent } from './_components/ktb-evaluation-details/ktb-evaluation-details.component';
import { KtbEventItemComponent, KtbEventItemDetail } from './_components/ktb-event-item/ktb-event-item.component';

import { KtbEventsListComponent } from './_components/ktb-events-list/ktb-events-list.component';
import { KtbExpandableTileComponent, KtbExpandableTileHeader } from './_components/ktb-expandable-tile/ktb-expandable-tile.component';
import { KtbHorizontalSeparatorComponent, KtbHorizontalSeparatorTitle } from './_components/ktb-horizontal-separator/ktb-horizontal-separator.component';
import { KtbHttpLoadingBarComponent } from './_components/ktb-http-loading-bar/ktb-http-loading-bar.component';
import { KtbNotificationBarComponent } from './_components/ktb-notification-bar/ktb-notification-bar.component';
import { KtbProjectListComponent } from './_components/ktb-project-list/ktb-project-list.component';
import { KtbProjectTileComponent } from './_components/ktb-project-tile/ktb-project-tile.component';
import { KtbRootEventsListComponent } from './_components/ktb-root-events-list/ktb-root-events-list.component';
import { KtbSelectableTileComponent } from './_components/ktb-selectable-tile/ktb-selectable-tile.component';
import { KtbSliBreakdownComponent } from './_components/ktb-sli-breakdown/ktb-sli-breakdown.component';
import { KtbHideHttpLoadingDirective } from './_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive';
import { KtbShowHttpLoadingDirective } from './_directives/ktb-show-http-loading/ktb-show-http-loading.directive';
import { KtbApprovalItemComponent } from "./_components/ktb-approval-item/ktb-approval-item.component";
import { KtbCopyToClipboardComponent } from "./_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.component";
import { KtbMarkdownComponent } from "./_components/ktb-markdown/ktb-markdown.component";

import { HttpErrorInterceptor } from './_interceptors/http-error-interceptor';
import { HttpLoadingInterceptor } from './_interceptors/http-loading-interceptor';
import { HttpDefaultInterceptor } from "./_interceptors/http-default-interceptor";

import { AtobPipe } from './_pipes/atob.pipe';
import { AppHeaderComponent } from './app-header/app-header.component';
import { AppComponent } from './app.component';

import { AppRouting } from './app.routing';

import { DashboardComponent } from './dashboard/dashboard.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { EvaluationBoardComponent } from "./evaluation-board/evaluation-board.component";
import { KtbSequenceTimelineComponent } from './_components/ktb-sequence-timeline/ktb-sequence-timeline.component';
import { KtbEnvironmentViewComponent } from './_views/ktb-environment-view/ktb-environment-view.component';

registerLocaleData(localeEn, 'en');

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    AppHeaderComponent,
    ProjectBoardComponent,
    EvaluationBoardComponent,
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
    KtbNotificationBarComponent,
    KtbApprovalItemComponent,
    KtbCopyToClipboardComponent,
    KtbMarkdownComponent,
    KtbSequenceTimelineComponent,
    KtbEnvironmentViewComponent,
  ],
  imports: [
    BrowserModule,
    FormsModule,
    BrowserAnimationsModule,
    HttpClientModule,
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
    DtOverlayModule,
    DtCheckboxModule,
    DtSwitchModule,
    DtConfirmationDialogModule,
    DtToastModule,
    DtTopBarNavigationModule,
    DtCopyToClipboardModule,
    DtToggleButtonGroupModule,
    DtQuickFilterModule,
    MatDialogModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    BrowserAnimationsModule,
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpDefaultInterceptor,
      multi: true,
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpErrorInterceptor,
      multi: true,
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpLoadingInterceptor,
      multi: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
}
