import { APP_BASE_HREF, registerLocaleData } from '@angular/common';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import localeEn from '@angular/common/locales/en';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
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
import { DtExpandableSectionModule } from '@dynatrace/barista-components/expandable-section';
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
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtTopBarNavigationModule } from '@dynatrace/barista-components/top-bar-navigation';
import { DtCopyToClipboardModule } from '@dynatrace/barista-components/copy-to-clipboard';
import { DtToggleButtonGroupModule } from '@dynatrace/barista-components/toggle-button-group';
import { DtQuickFilterModule } from '@dynatrace/barista-components/quick-filter';

import { DtTileModule } from '@dynatrace/barista-components/tile';
import { DtToastModule } from '@dynatrace/barista-components/toast';

import { MomentModule } from 'ngx-moment';

import { KtbEventsListComponent } from './_components/ktb-events-list/ktb-events-list.component';
import { KtbExpandableTileComponent, KtbExpandableTileHeader } from './_components/ktb-expandable-tile/ktb-expandable-tile.component';
import { KtbHorizontalSeparatorComponent, KtbHorizontalSeparatorTitle } from './_components/ktb-horizontal-separator/ktb-horizontal-separator.component';
import { KtbHttpLoadingBarComponent } from './_components/ktb-http-loading-bar/ktb-http-loading-bar.component';
import { KtbNotificationBarComponent } from './_components/ktb-notification-bar/ktb-notification-bar.component';
import { KtbProjectListComponent } from './_components/ktb-project-list/ktb-project-list.component';
import { KtbProjectTileComponent } from './_components/ktb-project-tile/ktb-project-tile.component';
import { KtbRootEventsListComponent } from './_components/ktb-root-events-list/ktb-root-events-list.component';
import { KtbSelectableTileComponent, KtbSelectableTileHeaderDirective } from './_components/ktb-selectable-tile/ktb-selectable-tile.component';
import { KtbSliBreakdownComponent } from './_components/ktb-sli-breakdown/ktb-sli-breakdown.component';
import { KtbHideHttpLoadingDirective } from './_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive';
import { KtbShowHttpLoadingDirective } from './_directives/ktb-show-http-loading/ktb-show-http-loading.directive';
import { KtbApprovalItemComponent } from './_components/ktb-approval-item/ktb-approval-item.component';
import { KtbCopyToClipboardComponent } from './_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.component';
import { KtbMarkdownComponent } from './_components/ktb-markdown/ktb-markdown.component';
import { KtbEvaluationDetailsComponent } from './_components/ktb-evaluation-details/ktb-evaluation-details.component';
import { KtbEvaluationInfoComponent } from './_components/ktb-evaluation-info/ktb-evaluation-info.component';
import { KtbEventItemComponent, KtbEventItemDetail } from './_components/ktb-event-item/ktb-event-item.component';
import { KtbTaskItemComponent, KtbTaskItemDetail } from './_components/ktb-task-item/ktb-task-item.component';
import { KtbSequenceTasksListComponent } from './_components/ktb-sequence-tasks-list/ktb-sequence-tasks-list.component';

import { HttpErrorInterceptor } from './_interceptors/http-error-interceptor';
import { HttpLoadingInterceptor } from './_interceptors/http-loading-interceptor';
import { HttpDefaultInterceptor } from './_interceptors/http-default-interceptor';

import { AppComponent } from './app.component';
import { AppRouting } from './app.routing';
import { AppHeaderComponent } from './app-header/app-header.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { EvaluationBoardComponent } from './evaluation-board/evaluation-board.component';
import { KtbSequenceTimelineComponent } from './_components/ktb-sequence-timeline/ktb-sequence-timeline.component';
import { KtbEnvironmentViewComponent } from './_views/ktb-environment-view/ktb-environment-view.component';
import { KtbIntegrationViewComponent } from './_views/ktb-integration-view/ktb-integration-view.component';
import { KtbStageOverviewComponent } from './_components/ktb-stage-overview/ktb-stage-overview.component';
import { KtbStageDetailsComponent } from './_components/ktb-stage-details/ktb-stage-details.component';
import { KtbSequenceViewComponent } from './_views/ktb-sequence-view/ktb-sequence-view.component';
import { KtbServiceViewComponent } from './_views/ktb-service-view/ktb-service-view.component';
import { KeptnUrlPipe } from './_pipes/keptn-url.pipe';
import { KtbSliBreakdownCriteriaItemComponent } from './_components/ktb-sli-breakdown-criteria-item/ktb-sli-breakdown-criteria-item.component';
import { KtbServicesListComponent } from './_components/ktb-services-list/ktb-services-list.component';
import { KtbStageBadgeComponent } from './_components/ktb-stage-badge/ktb-stage-badge.component';
import { KtbUniformViewComponent } from './_views/ktb-uniform-view/ktb-uniform-view.component';
import { KtbKeptnServicesListComponent } from './_components/ktb-keptn-services-list/ktb-keptn-services-list.component';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { KtbDeploymentListComponent } from './_components/ktb-deployment-list/ktb-deployment-list.component';
import { KtbUserComponent } from './_components/ktb-user/ktb-user.component';
import { KtbServiceDetailsComponent } from './_components/ktb-service-details/ktb-service-details.component';
import { KtbSettingsViewComponent } from './_views/ktb-settings-view/ktb-settings-view.component';
import { KtbDeploymentStageTimelineComponent } from './_components/ktb-deployment-stage-timeline/ktb-deployment-stage-timeline.component';
import { KtbSequenceListComponent } from './_components/ktb-sequence-list/ktb-sequence-list.component';
import { KtbUniformRegistrationLogsComponent } from './_components/ktb-uniform-registration-logs/ktb-uniform-registration-logs.component';
import { AppInitService } from './_services/app.init';
import { KtbSecretsListComponent } from './_components/ktb-secrets-list/ktb-secrets-list.component';
import { KtbCreateSecretFormComponent } from './_components/ktb-create-secret-form/ktb-create-secret-form.component';
import { KtbNoServiceInfoComponent } from './_components/ktb-no-service-info/ktb-no-service-info.component';
import { KtbSequenceStateListComponent } from './_components/ktb-sequence-state-list/ktb-sequence-state-list.component';
import { KtbProjectSettingsGitComponent } from './_components/ktb-project-settings-git/ktb-project-settings-git.component';
import { KtbProjectSettingsShipyardComponent } from './_components/ktb-project-settings-shipyard/ktb-project-settings-shipyard.component';
import { KtbDragAndDropDirective } from './_directives/ktb-drag-and-drop/ktb-drag-and-drop.directive';
import { KtbDangerZoneComponent } from './_components/ktb-danger-zone/ktb-danger-zone.component';
import { KtbDeletionDialogComponent } from './_components/_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.component';
import { EventService } from './_services/event.service';
import { ToType } from './_pipes/to-type';
import { KtbUniformSubscriptionsComponent } from './_components/ktb-uniform-subscriptions/ktb-uniform-subscriptions.component';
import { ToDatePipe } from './_pipes/to-date.pipe';
import { DtThemingModule } from '@dynatrace/barista-components/theming';
import { KtbSubscriptionItemComponent } from './_components/ktb-subscription-item/ktb-subscription-item.component';
import { POLLING_INTERVAL_MILLIS, RETRY_ON_HTTP_ERROR } from './_utils/app.utils';
import { environment } from '../environments/environment';

registerLocaleData(localeEn, 'en');

export function init_app(appLoadService: AppInitService) {
  return () => appLoadService.init();
}

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    AppHeaderComponent,
    ProjectBoardComponent,
    EvaluationBoardComponent,
    KtbSequenceViewComponent,
    KtbServiceViewComponent,
    KtbHttpLoadingBarComponent,
    KtbShowHttpLoadingDirective,
    KtbHideHttpLoadingDirective,
    KtbExpandableTileComponent,
    KtbExpandableTileHeader,
    KtbSelectableTileComponent,
    KtbSelectableTileHeaderDirective,
    KtbHorizontalSeparatorComponent,
    KtbHorizontalSeparatorTitle,
    KtbRootEventsListComponent,
    KtbProjectTileComponent,
    KtbProjectListComponent,
    KtbEventsListComponent,
    KtbEventItemComponent,
    KtbEventItemDetail,
    KtbSequenceTasksListComponent,
    KtbTaskItemComponent,
    KtbTaskItemDetail,
    KtbEvaluationDetailsComponent,
    KtbEvaluationInfoComponent,
    KtbStageBadgeComponent,
    KtbSliBreakdownComponent,
    KtbNotificationBarComponent,
    KtbApprovalItemComponent,
    KtbCopyToClipboardComponent,
    KtbMarkdownComponent,
    KtbSequenceTimelineComponent,
    KtbEnvironmentViewComponent,
    KtbStageOverviewComponent,
    KtbIntegrationViewComponent,
    KtbStageDetailsComponent,
    KeptnUrlPipe,
    KtbSliBreakdownCriteriaItemComponent,
    KtbServicesListComponent,
    KtbSequenceStateListComponent,
    KtbUserComponent,
    KtbUniformViewComponent,
    KtbKeptnServicesListComponent,
    KtbSubscriptionItemComponent,
    KtbDeploymentListComponent,
    KtbServiceDetailsComponent,
    KtbSettingsViewComponent,
    KtbDeploymentStageTimelineComponent,
    KtbSequenceListComponent,
    KtbUniformRegistrationLogsComponent,
    KtbSecretsListComponent,
    KtbCreateSecretFormComponent,
    KtbNoServiceInfoComponent,
    KtbProjectSettingsGitComponent,
    KtbProjectSettingsShipyardComponent,
    KtbDragAndDropDirective,
    KtbDangerZoneComponent,
    KtbDeletionDialogComponent,
    ToType,
    KtbUniformSubscriptionsComponent,
    ToDatePipe,
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
    DtTableModule,
    DtTagModule,
    DtExpandableTextModule,
    DtExpandablePanelModule,
    DtExpandableSectionModule,
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
    DtFilterFieldModule,
    ReactiveFormsModule,
  ],
  entryComponents: [
    KtbDeletionDialogComponent,
  ],
  providers: [
    EventService,
    AppInitService,
    {
      provide: APP_BASE_HREF,
      useValue: environment.baseUrl,
    },
    {
      provide: APP_INITIALIZER,
      useFactory: init_app,
      deps: [AppInitService],
      multi: true,
    },
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
    {
      provide: POLLING_INTERVAL_MILLIS,
      useValue: 30_000,
    },
    {
      provide: RETRY_ON_HTTP_ERROR,
      useValue: true,
    },
  ],
  bootstrap: [AppComponent],
})
export class AppModule {
}
