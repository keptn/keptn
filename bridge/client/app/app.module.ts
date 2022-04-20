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
import { MomentModule } from 'ngx-moment';
import {
  KtbExpandableTileComponent,
  KtbExpandableTileHeaderDirective,
} from './_components/ktb-expandable-tile/ktb-expandable-tile.component';
import {
  KtbHorizontalSeparatorComponent,
  KtbHorizontalSeparatorTitleDirective,
} from './_components/ktb-horizontal-separator/ktb-horizontal-separator.component';
import { KtbNotificationBarComponent } from './_components/ktb-notification-bar/ktb-notification-bar.component';
import { KtbProjectListComponent } from './_components/ktb-project-list/ktb-project-list.component';
import { KtbProjectTileComponent } from './_components/ktb-project-tile/ktb-project-tile.component';
import { KtbRootEventsListComponent } from './_components/ktb-root-events-list/ktb-root-events-list.component';
import {
  KtbSelectableTileComponent,
  KtbSelectableTileHeaderDirective,
} from './_components/ktb-selectable-tile/ktb-selectable-tile.component';
import { KtbSliBreakdownComponent } from './_components/ktb-sli-breakdown/ktb-sli-breakdown.component';
import { KtbHideHttpLoadingDirective } from './_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive';
import { KtbShowHttpLoadingDirective } from './_directives/ktb-show-http-loading/ktb-show-http-loading.directive';
import { KtbApprovalItemComponent } from './_components/ktb-approval-item/ktb-approval-item.component';
import { KtbCopyToClipboardComponent } from './_components/ktb-copy-to-clipboard/ktb-copy-to-clipboard.component';
import { KtbMarkdownComponent } from './_components/ktb-markdown/ktb-markdown.component';
import { KtbEvaluationDetailsComponent } from './_components/ktb-evaluation-details/ktb-evaluation-details.component';
import { KtbEvaluationInfoComponent } from './_components/ktb-evaluation-info/ktb-evaluation-info.component';
import {
  KtbEventItemComponent,
  KtbEventItemDetailDirective,
} from './_components/ktb-event-item/ktb-event-item.component';
import { KtbTaskItemComponent, KtbTaskItemDetailDirective } from './_components/ktb-task-item/ktb-task-item.component';
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
import { KtbDeleteConfirmationComponent } from './_components/_dialogs/ktb-delete-confirmation/ktb-delete-confirmation.component';
import { KtbModifyUniformSubscriptionComponent } from './_components/ktb-modify-uniform-subscription/ktb-modify-uniform-subscription.component';
import { DtThemingModule } from '@dynatrace/barista-components/theming';
import { KtbSubscriptionItemComponent } from './_components/ktb-subscription-item/ktb-subscription-item.component';
import { KtbConfirmationDialogComponent } from './_components/_dialogs/ktb-confirmation-dialog/ktb-confirmation-dialog.component';
import { POLLING_INTERVAL_MILLIS, RETRY_ON_HTTP_ERROR } from './_utils/app.utils';
import { KtbSequenceControlsComponent } from './_components/ktb-sequence-controls/ktb-sequence-controls.component';
import { environment } from '../environments/environment';
import { KtbProjectSettingsComponent } from './_components/ktb-project-settings/ktb-project-settings.component';
import { KtbWebhookSettingsComponent } from './_components/ktb-webhook-settings/ktb-webhook-settings.component';
import { KtbServiceSettingsComponent } from './_components/ktb-service-settings/ktb-service-settings.component';
import { KtbCreateServiceComponent } from './_components/ktb-create-service/ktb-create-service.component';
import { KtbServiceSettingsOverviewComponent } from './_components/ktb-service-settings-overview/ktb-service-settings-overview.component';
import { KtbServiceSettingsListComponent } from './_components/ktb-service-settings-list/ktb-service-settings-list.component';
import { KtbEditServiceComponent } from './_components/ktb-edit-service/ktb-edit-service.component';
import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { KtbEditServiceFileListComponent } from './_components/ktb-edit-service-file-list/ktb-edit-service-file-list.component';
import { DtTreeTableModule } from '@dynatrace/barista-components/tree-table';
import {
  KtbTreeListSelectComponent,
  KtbTreeListSelectDirective,
} from './_components/ktb-tree-list-select/ktb-tree-list-select.component';
import { OverlayModule } from '@angular/cdk/overlay';
import { KtbSequenceStateInfoComponent } from './_components/ktb-sequence-state-info/ktb-sequence-state-info.component';
import { KtbPayloadViewerComponent } from './_components/ktb-payload-viewer/ktb-payload-viewer.component';
import { DtRadioModule } from '@dynatrace/barista-components/radio';
import { NotFoundComponent } from './not-found/not-found.component';
import { KtbVariableSelectorComponent } from './_components/ktb-variable-selector/ktb-variable-selector.component';
import { KtbNotificationComponent } from './_components/ktb-notification/ktb-notification.component';
import { KtbProjectCreateMessageComponent } from './_components/_status-messages/ktb-project-create-message/ktb-project-create-message.component';
import { PendingChangesGuard } from './_guards/pending-changes.guard';
import { ArrayToStringPipe } from './_pipes/array-to-string';
import { KtbTriggerSequenceComponent } from './_components/ktb-trigger-sequence/ktb-trigger-sequence.component';
import { KtbTimeInputComponent } from './_components/ktb-time-input/ktb-time-input.component';
import {
  KtbDatetimePickerComponent,
  KtbDatetimePickerDirective,
} from './_components/ktb-datetime-picker/ktb-datetime-picker.component';
import { DtDatepickerModule } from '@dynatrace/barista-components/experimental/datepicker';
import { TruncateNumberPipe } from './_pipes/truncate-number';
import { KtbErrorViewComponent } from './_views/ktb-error-view/ktb-error-view.component';
import { KtbRootComponent } from './ktb-root/ktb-root.component';
import { KtbLogoutViewComponent } from './_views/ktb-logout-view/ktb-logout-view.component';
import { KtbProjectSettingsGitExtendedComponent } from './_components/ktb-project-settings-git-extended/ktb-project-settings-git-extended.component';
import { KtbProjectSettingsGitHttpsComponent } from './_components/ktb-project-settings-git-https/ktb-project-settings-git-https.component';
import { KtbProjectSettingsGitSshComponent } from './_components/ktb-project-settings-git-ssh/ktb-project-settings-git-ssh.component';
import { KtbProxyInputComponent } from './_components/ktb-proxy-input/ktb-proxy-input.component';
import { KtbIntegerInputDirective } from './_directives/ktb-integer-input/ktb-integer-input.directive';
import { KtbCertificateInputComponent } from './_components/ktb-certificate-input/ktb-certificate-input.component';
import { KtbSshKeyInputComponent } from './_components/ktb-ssh-key-input/ktb-ssh-key-input.component';
import { KtbProjectSettingsGitSshInputComponent } from './_components/ktb-project-settings-git-ssh-input/ktb-project-settings-git-ssh-input.component';
import { KtbLoadingDistractorComponent } from './_components/ktb-loading-distractor/ktb-loading-distractor.component';
import { KtbLoadingSpinnerComponent } from './_components/ktb-loading-spinner/ktb-loading-spinner.component';
import { SanitizeHtmlPipe } from './_pipes/sanitize-html.pipe';

registerLocaleData(localeEn, 'en');

export function init_app(appLoadService: AppInitService): () => Promise<unknown> {
  return (): Promise<string | null> => appLoadService.init();
}

@NgModule({
  declarations: [
    AppComponent,
    DashboardComponent,
    NotFoundComponent,
    AppHeaderComponent,
    ProjectBoardComponent,
    EvaluationBoardComponent,
    KtbSequenceViewComponent,
    KtbServiceViewComponent,
    KtbShowHttpLoadingDirective,
    KtbHideHttpLoadingDirective,
    KtbExpandableTileComponent,
    KtbExpandableTileHeaderDirective,
    KtbSelectableTileComponent,
    KtbSelectableTileHeaderDirective,
    KtbHorizontalSeparatorComponent,
    KtbHorizontalSeparatorTitleDirective,
    KtbRootEventsListComponent,
    KtbProjectTileComponent,
    KtbProjectListComponent,
    KtbEventItemComponent,
    KtbEventItemDetailDirective,
    KtbSequenceTasksListComponent,
    KtbTaskItemComponent,
    KtbTaskItemDetailDirective,
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
    KtbConfirmationDialogComponent,
    KtbSequenceControlsComponent,
    ToType,
    KtbUniformSubscriptionsComponent,
    ToDatePipe,
    KtbProjectSettingsComponent,
    KtbDeleteConfirmationComponent,
    KtbModifyUniformSubscriptionComponent,
    KtbWebhookSettingsComponent,
    KtbServiceSettingsComponent,
    KtbCreateServiceComponent,
    KtbServiceSettingsOverviewComponent,
    KtbServiceSettingsListComponent,
    KtbEditServiceComponent,
    KtbEditServiceFileListComponent,
    KtbTreeListSelectComponent,
    KtbTreeListSelectDirective,
    KtbSequenceStateInfoComponent,
    KtbPayloadViewerComponent,
    KtbVariableSelectorComponent,
    ArrayToStringPipe,
    KtbNotificationComponent,
    KtbProjectCreateMessageComponent,
    KtbTriggerSequenceComponent,
    KtbTimeInputComponent,
    KtbDatetimePickerComponent,
    KtbDatetimePickerDirective,
    TruncateNumberPipe,
    KtbErrorViewComponent,
    KtbRootComponent,
    KtbLogoutViewComponent,
    KtbProjectSettingsGitExtendedComponent,
    KtbProjectSettingsGitHttpsComponent,
    KtbProjectSettingsGitSshComponent,
    KtbProxyInputComponent,
    KtbIntegerInputDirective,
    KtbCertificateInputComponent,
    KtbSshKeyInputComponent,
    KtbProjectSettingsGitSshInputComponent,
    KtbLoadingDistractorComponent,
    KtbLoadingSpinnerComponent,
    SanitizeHtmlPipe,
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
    DtTopBarNavigationModule,
    DtCopyToClipboardModule,
    DtToggleButtonGroupModule,
    DtQuickFilterModule,
    DtRadioModule,
    MatDialogModule,
    DtIconModule.forRoot({
      svgIconLocation: `assets/icons/{{name}}.svg`,
    }),
    BrowserAnimationsModule,
    DtFilterFieldModule,
    ReactiveFormsModule,
    DtAlertModule,
    DtTreeTableModule,
    OverlayModule,
    DtDatepickerModule,
  ],
  entryComponents: [KtbDeletionDialogComponent, KtbConfirmationDialogComponent],
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
