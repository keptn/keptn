import { OverlayModule } from '@angular/cdk/overlay';
import { APP_BASE_HREF, CommonModule, registerLocaleData } from '@angular/common';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
import localeEn from '@angular/common/locales/en';
import { APP_INITIALIZER, NgModule } from '@angular/core';
import { FlexLayoutModule } from '@angular/flex-layout';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { MatDialogModule } from '@angular/material/dialog';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtButtonGroupModule } from '@dynatrace/barista-components/button-group';
import { DtCardModule } from '@dynatrace/barista-components/card';
import { DtCheckboxModule } from '@dynatrace/barista-components/checkbox';
import { DtConfirmationDialogModule } from '@dynatrace/barista-components/confirmation-dialog';
import { DtContextDialogModule } from '@dynatrace/barista-components/context-dialog';
import { DtCopyToClipboardModule } from '@dynatrace/barista-components/copy-to-clipboard';
import { DtDrawerModule } from '@dynatrace/barista-components/drawer';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtExpandableSectionModule } from '@dynatrace/barista-components/expandable-section';
import { DtExpandableTextModule } from '@dynatrace/barista-components/expandable-text';
import { DtDatepickerModule } from '@dynatrace/barista-components/experimental/datepicker';
import { DtFilterFieldModule } from '@dynatrace/barista-components/filter-field';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtIndicatorModule } from '@dynatrace/barista-components/indicator';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtInputModule } from '@dynatrace/barista-components/input';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtProgressBarModule } from '@dynatrace/barista-components/progress-bar';
import { DtProgressCircleModule } from '@dynatrace/barista-components/progress-circle';
import { DtQuickFilterModule } from '@dynatrace/barista-components/quick-filter';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { DtSelectModule } from '@dynatrace/barista-components/select';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtSwitchModule } from '@dynatrace/barista-components/switch';
import { DtTableModule } from '@dynatrace/barista-components/table';
import { DtTagModule } from '@dynatrace/barista-components/tag';
import { DtThemingModule } from '@dynatrace/barista-components/theming';
import { DtTileModule } from '@dynatrace/barista-components/tile';
import { DtToggleButtonGroupModule } from '@dynatrace/barista-components/toggle-button-group';
import { DtTopBarNavigationModule } from '@dynatrace/barista-components/top-bar-navigation';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import { MomentModule } from 'ngx-moment';
import { environment } from '../environments/environment';
import { WindowConfig } from '../environments/environment.dynamic';
import { KtbConfirmationDialogModule } from './_components/_dialogs/ktb-confirmation-dialog/ktb-confirmation-dialog.module';
import { KtbDeleteConfirmationModule } from './_components/_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.module';
import { KtbDeletionDialogModule } from './_components/_dialogs/ktb-deletion-dialog/ktb-deletion-dialog.module';
import { KtbApprovalItemModule } from './_components/ktb-approval-item/ktb-approval-item.module';
import { KtbCertificateInputModule } from './_components/ktb-certificate-input/ktb-certificate-input.module';
import { KtbCopyToClipboardModule } from './_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.module';
import { KtbCreateSecretFormModule } from './_components/ktb-create-secret-form/ktb-create-secret-form.module';
import { KtbCreateServiceModule } from './_components/ktb-create-service/ktb-create-service.module';
import { KtbDangerZoneModule } from './_components/ktb-danger-zone/ktb-danger-zone.module';
import { KtbDateInputModule } from './_components/ktb-date-input/ktb-date-input.module';
import { KtbDeploymentListModule } from './_components/ktb-deployment-list/ktb-deployment-list.module';
import { KtbDeploymentStageTimelineModule } from './_components/ktb-deployment-stage-timeline/ktb-deployment-stage-timeline.module';
import { KtbEditServiceModule } from './_components/ktb-edit-service/ktb-edit-service.module';
import { KtbEvaluationDetailsModule } from './_components/ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbEvaluationInfoModule } from './_components/ktb-evaluation-info/ktb-evaluation-info.module';
import { KtbEventItemModule } from './_components/ktb-event-item/ktb-event-item.module';
import { KtbExpandableTileModule } from './_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { KtbHeatmapModule } from './_components/ktb-heatmap/ktb-heatmap.module';
import { KtbLoadingModule } from './_components/ktb-loading/ktb-loading.module';
import { KtbNotificationModule } from './_components/ktb-notification/ktb-notification.module';
import { KtbProjectListModule } from './_components/ktb-project-list/ktb-project-list.module';
import { KtbProjectSettingsModule } from './_components/ktb-project-settings/ktb-project-settings.module';
import { KtbProxyInputModule } from './_components/ktb-proxy-input/ktb-proxy-input.module';
import { KtbRootEventsListModule } from './_components/ktb-root-events-list/ktb-root-events-list.module';
import { KtbSelectableTileModule } from './_components/ktb-selectable-tile/ktb-selectable-tile.module';
import { KtbSequenceControlsModule } from './_components/ktb-sequence-controls/ktb-sequence-controls.module';
import { KtbSequenceStateInfoModule } from './_components/ktb-sequence-state-info/ktb-sequence-state-info.module';
import { KtbSequenceStateListModule } from './_components/ktb-sequence-state-list/ktb-sequence-state-list.module';
import { KtbSshKeyInputModule } from './_components/ktb-ssh-key-input/ktb-ssh-key-input.module';
import { KtbStageBadgeModule } from './_components/ktb-stage-badge/ktb-stage-badge.module';
import { KtbDragAndDropModule } from './_directives/ktb-drag-and-drop/ktb-drag-and-drop.module';
import { KtbHideHttpLoadingDirective } from './_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive';
import { KtbIntegerInputModule } from './_directives/ktb-integer-input/ktb-integer-input.module';
import { KtbShowHttpLoadingDirective } from './_directives/ktb-show-http-loading/ktb-show-http-loading.directive';
import { PendingChangesGuard } from './_guards/pending-changes.guard';
import { HttpDefaultInterceptor } from './_interceptors/http-default-interceptor';
import { HttpErrorInterceptor } from './_interceptors/http-error-interceptor';
import { HttpLoadingInterceptor } from './_interceptors/http-loading-interceptor';
import { KtbPipeModule } from './_pipes/ktb-pipe.module';
import { AppInitService } from './_services/app.init';
import { EventService } from './_services/event.service';
import { POLLING_INTERVAL_MILLIS, RETRY_ON_HTTP_ERROR } from './_utils/app.utils';
import { KtbEnvironmentViewComponent } from './_views/ktb-environment-view/ktb-environment-view.component';
import { KtbErrorViewModule } from './_views/ktb-error-view/ktb-error-view.module';
import { KtbIntegrationViewComponent } from './_views/ktb-integration-view/ktb-integration-view.component';
import { KtbLogoutViewComponent } from './_views/ktb-logout-view/ktb-logout-view.component';
import { KtbSequenceViewComponent } from './_views/ktb-sequence-view/ktb-sequence-view.component';
import { KtbServiceViewComponent } from './_views/ktb-service-view/ktb-service-view.component';
import { KtbSettingsViewComponent } from './_views/ktb-settings-view/ktb-settings-view.component';
import { AppComponent } from './app.component';
import { AppRouting } from './app.routing';
import { DashboardLegacyComponent } from './dashboard-legacy/dashboard-legacy.component';
import { EvaluationBoardComponent } from './evaluation-board/evaluation-board.component';
import { KtbRootComponent } from './ktb-root/ktb-root.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { KtbSubscriptionItemModule } from './_components/ktb-subscription-item/ktb-subscription-item.module';
import { KtbUniformSubscriptionsModule } from './_components/ktb-uniform-subscriptions/ktb-uniform-subscriptions.module';
import { KtbKeptnServicesListModule } from './_components/ktb-keptn-services-list/ktb-keptn-services-list.module';
import { KtbUniformRegistrationLogsModule } from './_components/ktb-uniform-registration-logs/ktb-uniform-registration-logs.module';
import { KtbPayloadViewerModule } from './_components/ktb-payload-viewer/ktb-payload-viewer.module';
import { KtbVariableSelectorModule } from './_components/ktb-variable-selector/ktb-variable-selector.module';
import { KtbTreeListSelectModule } from './_components/ktb-tree-list-select/ktb-tree-list-select.module';
import { KtbWebhookSettingsModule } from './_components/ktb-webhook-settings/ktb-webhook-settings.module';
import { KtbModifyUniformSubscriptionModule } from './_components/ktb-modify-uniform-subscription/ktb-modify-uniform-subscription.module';
import { KtbNoServiceInfoModule } from './_components/ktb-no-service-info/ktb-no-service-info.module';
import { KtbSequenceListModule } from './_components/ktb-sequence-list/ktb-sequence-list.module';
import { KtbSequenceTasksListModule } from './_components/ktb-sequence-tasks-list/ktb-sequence-tasks-list.module';
import { KtbSecretsListModule } from './_components/ktb-secrets-list/ktb-secrets-list.module';
import { KtbStageDetailsModule } from './_views/ktb-environment-view/ktb-stage-details/ktb-stage-details.module';
import { KtbStageOverviewModule } from './_views/ktb-environment-view/ktb-stage-overview/ktb-stage-overview.module';
import { KtbAppHeaderModule } from './_components/ktb-app-header/ktb-app-header.module';
import { KtbSequenceTimelineModule } from './_components/ktb-sequence-timeline/ktb-sequence-timeline.module';
import { KtbServiceDetailsModule } from './_components/ktb-service-details/ktb-service-details.module';
import { KtbServiceSettingsModule } from './_components/ktb-service-settings/ktb-service-settings.module';

registerLocaleData(localeEn, 'en');

export function init_app(appLoadService: AppInitService): () => Promise<unknown> {
  return (): Promise<WindowConfig | null> => appLoadService.init();
}

const angularModules = [
  BrowserModule,
  FormsModule,
  BrowserAnimationsModule,
  HttpClientModule,
  FlexLayoutModule,
  MatDialogModule,
  BrowserAnimationsModule,
  ReactiveFormsModule,
  CommonModule,
  OverlayModule,
];

const dtModules = [
  DtFilterFieldModule,
  DtIconModule.forRoot({
    svgIconLocation: `assets/icons/{{name}}.svg`,
  }),
  DtAlertModule,
  DtTreeTableModule,
  DtDatepickerModule,
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
  DtTableModule,
  DtTagModule,
  DtExpandableTextModule,
  DtExpandablePanelModule,
  DtExpandableSectionModule,
  DtShowMoreModule,
  DtIndicatorModule,
  DtProgressCircleModule,
  DtOverlayModule,
  DtCheckboxModule,
  DtSwitchModule,
  DtConfirmationDialogModule,
  DtTopBarNavigationModule,
  DtCopyToClipboardModule,
  DtToggleButtonGroupModule,
  DtQuickFilterModule,
  DtRadioModule,
];

const ktbModules = [
  KtbAppHeaderModule,
  KtbApprovalItemModule,
  KtbCertificateInputModule,
  KtbConfirmationDialogModule,
  KtbCopyToClipboardModule,
  KtbCreateSecretFormModule,
  KtbCreateServiceModule,
  KtbDangerZoneModule,
  KtbDateInputModule,
  KtbDeleteConfirmationModule,
  KtbDeletionDialogModule,
  KtbDeploymentListModule,
  KtbDeploymentStageTimelineModule,
  KtbDragAndDropModule,
  KtbEditServiceModule,
  KtbErrorViewModule,
  KtbEvaluationDetailsModule,
  KtbEvaluationInfoModule,
  KtbEventItemModule,
  KtbExpandableTileModule,
  KtbHeatmapModule,
  KtbIntegerInputModule,
  KtbKeptnServicesListModule,
  KtbLoadingModule,
  KtbModifyUniformSubscriptionModule,
  KtbNoServiceInfoModule,
  KtbNotificationModule,
  KtbPayloadViewerModule,
  KtbPipeModule,
  KtbProjectListModule,
  KtbProjectSettingsModule,
  KtbProxyInputModule,
  KtbRootEventsListModule,
  KtbSecretsListModule,
  KtbSelectableTileModule,
  KtbSequenceControlsModule,
  KtbSequenceListModule,
  KtbSequenceStateInfoModule,
  KtbSequenceStateListModule,
  KtbSequenceTasksListModule,
  KtbSequenceTimelineModule,
  KtbServiceDetailsModule,
  KtbServiceSettingsModule,
  KtbSshKeyInputModule,
  KtbStageBadgeModule,
  KtbStageDetailsModule,
  KtbStageOverviewModule,
  KtbSubscriptionItemModule,
  KtbTreeListSelectModule,
  KtbUniformRegistrationLogsModule,
  KtbUniformSubscriptionsModule,
  KtbVariableSelectorModule,
  KtbWebhookSettingsModule,
];

@NgModule({
  declarations: [
    AppComponent,
    DashboardLegacyComponent,
    NotFoundComponent,
    ProjectBoardComponent,
    EvaluationBoardComponent,
    KtbSequenceViewComponent,
    KtbServiceViewComponent,
    KtbShowHttpLoadingDirective,
    KtbHideHttpLoadingDirective,
    KtbEnvironmentViewComponent,
    KtbIntegrationViewComponent,
    KtbSettingsViewComponent,
    KtbRootComponent,
    KtbLogoutViewComponent,
  ],
  imports: [...angularModules, ...dtModules, ...ktbModules, AppRouting, MomentModule],
  entryComponents: [],
  providers: [
    EventService,
    AppInitService,
    PendingChangesGuard,
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
      useValue: environment.pollingIntervalMillis ?? 30_000,
    },
    {
      provide: RETRY_ON_HTTP_ERROR,
      useValue: true,
    },
  ],
  bootstrap: [KtbRootComponent],
})
export class AppModule {}
