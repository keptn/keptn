import { APP_INITIALIZER, NgModule } from '@angular/core';
import { APP_BASE_HREF, CommonModule, registerLocaleData } from '@angular/common';
import { HTTP_INTERCEPTORS, HttpClientModule } from '@angular/common/http';
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
import { DashboardLegacyComponent } from './dashboard-legacy/dashboard-legacy.component';
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
import { KtbErrorViewModule } from './_views/ktb-error-view/ktb-error-view.module';
import { KtbLoadingModule } from './_components/ktb-loading/ktb-loading.module';
import { KtbNotificationModule } from './_components/ktb-notification/ktb-notification.module';
import { KtbProjectListModule } from './_components/ktb-project-list/ktb-project-list.module';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { KtbEvaluationDetailsModule } from './_components/ktb-evaluation-details/ktb-evaluation-details.module';

// Import BrowserModule, BrowserAnimationsModule, HttpModule or HttpClientModule only once!

registerLocaleData(localeEn, 'en');

export function init_app(appLoadService: AppInitService): () => Promise<unknown> {
  return (): Promise<WindowConfig | null> => appLoadService.init();
}

const angularModules = [BrowserModule, BrowserAnimationsModule, HttpClientModule, CommonModule];

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
  KtbErrorViewModule,
  KtbLoadingModule,
  KtbNotificationModule,
  KtbPipeModule,
  KtbProjectListModule,
];

@NgModule({
  declarations: [
    AppComponent,
    DashboardLegacyComponent,
    NotFoundComponent,
    KtbSettingsViewComponent,
    KtbRootComponent,
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
