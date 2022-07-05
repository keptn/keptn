import { APP_INITIALIZER, NgModule } from '@angular/core';
import { APP_BASE_HREF, CommonModule, registerLocaleData } from '@angular/common';
import { HTTP_INTERCEPTORS } from '@angular/common/http';
import localeEn from '@angular/common/locales/en';
import { FlexModule } from '@angular/flex-layout';

import { MomentModule } from 'ngx-moment';

import { environment } from '../environments/environment';
import { WindowConfig } from '../environments/environment.dynamic';

import { POLLING_INTERVAL_MILLIS, RETRY_ON_HTTP_ERROR } from './_utils/app.utils';

import { AppComponent } from './app.component';
import { AppRouting } from './app.routing';
import { KtbRootComponent } from './ktb-root/ktb-root.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { DashboardLegacyComponent } from './dashboard-legacy/dashboard-legacy.component';
import { EvaluationBoardComponent } from './evaluation-board/evaluation-board.component';
import { KtbEnvironmentViewComponent } from './_views/ktb-environment-view/ktb-environment-view.component';
import { KtbLogoutViewComponent } from './_views/ktb-logout-view/ktb-logout-view.component';
import { KtbSequenceViewComponent } from './_views/ktb-sequence-view/ktb-sequence-view.component';
import { KtbServiceViewComponent } from './_views/ktb-service-view/ktb-service-view.component';
import { KtbSettingsViewComponent } from './_views/ktb-settings-view/ktb-settings-view.component';

import { AppInitService } from './_services/app.init';
import { EventService } from './_services/event.service';

import { KtbPipeModule } from './_pipes/ktb-pipe.module';
import { PendingChangesGuard } from './_guards/pending-changes.guard';

import { HttpDefaultInterceptor } from './_interceptors/http-default-interceptor';
import { HttpErrorInterceptor } from './_interceptors/http-error-interceptor';
import { HttpLoadingInterceptor } from './_interceptors/http-loading-interceptor';

import { DtAlertModule } from '@dynatrace/barista-components/alert';
import { DtButtonModule } from '@dynatrace/barista-components/button';
import { DtEmptyStateModule } from '@dynatrace/barista-components/empty-state';
import { DtExpandablePanelModule } from '@dynatrace/barista-components/expandable-panel';
import { DtIconModule } from '@dynatrace/barista-components/icon';
import { DtInfoGroupModule } from '@dynatrace/barista-components/info-group';
import { DtMenuModule } from '@dynatrace/barista-components/menu';
import { DtOverlayModule } from '@dynatrace/barista-components/overlay';
import { DtQuickFilterModule } from '@dynatrace/barista-components/quick-filter';
import { DtShowMoreModule } from '@dynatrace/barista-components/show-more';
import { DtTagModule } from '@dynatrace/barista-components/tag';

import { KtbAppHeaderModule } from './_components/ktb-app-header/ktb-app-header.module';
import { KtbCreateSecretFormModule } from './_components/ktb-create-secret-form/ktb-create-secret-form.module';
import { KtbCreateServiceModule } from './_components/ktb-create-service/ktb-create-service.module';
import { KtbDeploymentListModule } from './_components/ktb-deployment-list/ktb-deployment-list.module';
import { KtbEditServiceModule } from './_components/ktb-edit-service/ktb-edit-service.module';
import { KtbErrorViewModule } from './_views/ktb-error-view/ktb-error-view.module';
import { KtbEvaluationDetailsModule } from './_components/ktb-evaluation-details/ktb-evaluation-details.module';
import { KtbEventItemModule } from './_components/ktb-event-item/ktb-event-item.module';
import { KtbExpandableTileModule } from './_components/ktb-expandable-tile/ktb-expandable-tile.module';
import { KtbLoadingModule } from './_components/ktb-loading/ktb-loading.module';
import { KtbNoServiceInfoModule } from './_components/ktb-no-service-info/ktb-no-service-info.module';
import { KtbNotificationModule } from './_components/ktb-notification/ktb-notification.module';
import { KtbProjectListModule } from './_components/ktb-project-list/ktb-project-list.module';
import { KtbProjectSettingsModule } from './_components/ktb-project-settings/ktb-project-settings.module';
import { KtbRootEventsListModule } from './_components/ktb-root-events-list/ktb-root-events-list.module';
import { KtbSecretsListModule } from './_components/ktb-secrets-list/ktb-secrets-list.module';
import { KtbSequenceControlsModule } from './_components/ktb-sequence-controls/ktb-sequence-controls.module';
import { KtbSequenceTasksListModule } from './_components/ktb-sequence-tasks-list/ktb-sequence-tasks-list.module';
import { KtbSequenceTimelineModule } from './_components/ktb-sequence-timeline/ktb-sequence-timeline.module';
import { KtbServiceDetailsModule } from './_components/ktb-service-details/ktb-service-details.module';
import { KtbServiceSettingsModule } from './_components/ktb-service-settings/ktb-service-settings.module';
import { KtbStageDetailsModule } from './_views/ktb-environment-view/ktb-stage-details/ktb-stage-details.module';
import { KtbStageOverviewModule } from './_views/ktb-environment-view/ktb-stage-overview/ktb-stage-overview.module';

import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';

registerLocaleData(localeEn, 'en');

export function init_app(appLoadService: AppInitService): () => Promise<unknown> {
  return (): Promise<WindowConfig | null> => appLoadService.init();
}

const angularModules = [BrowserModule, BrowserAnimationsModule, CommonModule];

const dtModules = [
  DtAlertModule,
  DtButtonModule,
  DtEmptyStateModule,
  DtExpandablePanelModule,
  DtIconModule.forRoot({
    svgIconLocation: `assets/icons/{{name}}.svg`,
  }),
  DtInfoGroupModule,
  DtMenuModule,
  DtOverlayModule,
  DtQuickFilterModule,
  DtShowMoreModule,
  DtTagModule,
];

const ktbModules = [
  KtbAppHeaderModule,
  KtbCreateSecretFormModule,
  KtbCreateServiceModule,
  KtbDeploymentListModule,
  KtbEditServiceModule,
  KtbErrorViewModule,
  KtbEvaluationDetailsModule,
  KtbEventItemModule,
  KtbExpandableTileModule,
  KtbLoadingModule,
  KtbNoServiceInfoModule,
  KtbNotificationModule,
  KtbPipeModule,
  KtbProjectListModule,
  KtbProjectSettingsModule,
  KtbRootEventsListModule,
  KtbSecretsListModule,
  KtbSequenceControlsModule,
  KtbSequenceTasksListModule,
  KtbSequenceTimelineModule,
  KtbServiceDetailsModule,
  KtbServiceSettingsModule,
  KtbStageDetailsModule,
  KtbStageOverviewModule,
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
    KtbEnvironmentViewComponent,
    KtbSettingsViewComponent,
    KtbRootComponent,
    KtbLogoutViewComponent,
  ],
  imports: [...angularModules, ...dtModules, ...ktbModules, AppRouting, FlexModule, MomentModule],
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
