import {NgModule} from '@angular/core';
import {HttpClientModule} from '@angular/common/http';
import {HTTP_INTERCEPTORS} from '@angular/common/http';
import {BrowserModule} from '@angular/platform-browser';
import {FlexLayoutModule} from "@angular/flex-layout";
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MomentModule} from "ngx-moment";

import {AppRouting} from './app.routing';
import {AppComponent} from './app.component';

import {DashboardComponent} from './dashboard/dashboard.component';
import {AppHeaderComponent} from './app-header/app-header.component';
import {ProjectBoardComponent} from './project-board/project-board.component';

import {KtbProjectTileComponent} from './_components/ktb-project-tile/ktb-project-tile.component';
import {KtbSelectableTileComponent} from "./_components/ktb-selectable-tile/ktb-selectable-tile.component";
import {KtbHttpLoadingBarComponent} from "./_components/ktb-http-loading-bar/ktb-http-loading-bar.component";
import {KtbRootEventsListComponent} from "./_components/ktb-root-events-list/ktb-root-events-list.component";
import {KtbHttpLoadingSpinnerComponent} from './_components/ktb-http-loading-spinner/ktb-http-loading-spinner.component';
import {KtbExpandableTileComponent, KtbExpandableTileHeader} from './_components/ktb-expandable-tile/ktb-expandable-tile.component';
import {KtbHorizontalSeparatorComponent, KtbHorizontalSeparatorTitle} from "./_components/ktb-horizontal-separator/ktb-horizontal-separator.component";

import {KtbShowHttpLoadingDirective} from './_directives/ktb-show-http-loading/ktb-show-http-loading.directive';
import {KtbHideHttpLoadingDirective} from "./_directives/ktb-hide-http-loading/ktb-hide-http-loading.directive";

import {HttpErrorInterceptor} from "./_interceptors/http-error-interceptor";
import {HttpLoadingInterceptor} from "./_interceptors/http-loading-interceptor";

import {DtThemingModule} from "@dynatrace/barista-components/theming";
import {DtButtonModule} from "@dynatrace/barista-components/button";
import {DtContextDialogModule} from "@dynatrace/barista-components/context-dialog";
import {DtEmptyStateModule} from "@dynatrace/barista-components/empty-state";
import {DtSelectModule} from "@dynatrace/barista-components/select";
import {DtMenuModule} from "@dynatrace/barista-components/menu";
import {DtDrawerModule} from "@dynatrace/barista-components/drawer";
import {DtInputModule} from "@dynatrace/barista-components/input";
import {DtCardModule} from "@dynatrace/barista-components/card";
import {DtTileModule} from "@dynatrace/barista-components/tile";
import {DtInfoGroupModule} from "@dynatrace/barista-components/info-group";
import {DtProgressBarModule} from "@dynatrace/barista-components/progress-bar";
import {DtLoadingDistractorModule} from "@dynatrace/barista-components/loading-distractor";
import {DtTagModule} from "@dynatrace/barista-components/tag";
import {DtExpandableTextModule} from "@dynatrace/barista-components/expandable-text";
import {DtExpandablePanelModule} from "@dynatrace/barista-components/expandable-panel";
import {DtShowMoreModule} from "@dynatrace/barista-components/show-more";
import {DtIconModule} from "@dynatrace/barista-components/icon";

import {registerLocaleData} from "@angular/common";
import localeEn from '@angular/common/locales/en';
import { KtbProjectListComponent } from './_components/ktb-project-list/ktb-project-list.component';
import {DtIndicatorModule} from "@dynatrace/barista-components/core";
import {DtProgressCircleModule} from "@dynatrace/barista-components/progress-circle";

registerLocaleData(localeEn, 'en');

@NgModule({
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
  ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    HttpClientModule,
    AppRouting,
    FlexLayoutModule,
    MomentModule,
    DtThemingModule,
    DtButtonModule,
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
    DtIconModule.forRoot({
      svgIconLocation: `/assets/icons/{{name}}.svg`,
    }),
    BrowserAnimationsModule
  ],
  providers: [
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpErrorInterceptor,
      multi: true,
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpLoadingInterceptor,
      multi: true
    },
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
