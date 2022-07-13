import { NgModule } from '@angular/core';
import { ExtraOptions, RouterModule, Routes } from '@angular/router';
import { NotFoundComponent } from './not-found/not-found.component';
import { KtbErrorViewComponent } from './_views/ktb-error-view/ktb-error-view.component';
import { AppComponent } from './app.component';
import { TraceDeepLinkGuard } from './_guards/trace-deep-link.guard';
import { KtbTestChartComponent } from './_components/ktb-chart/testing/ktb-test-chart.component';

const routingConfiguration: ExtraOptions = {
  paramsInheritanceStrategy: 'always',
};

const routes: Routes = [
  { path: 'error', component: KtbErrorViewComponent },
  {
    path: 'logoutsession',
    loadChildren: () => import('./_views/ktb-logout-view/ktb-logout-view.module').then((m) => m.KtbLogoutViewModule),
  },
  {
    path: '',
    component: AppComponent,
    children: [
      { path: '', pathMatch: 'full', redirectTo: 'dashboard' },
      {
        path: 'dashboard',
        loadChildren: () =>
          import('./_views/ktb-dashboard-view/ktb-dashboard-view.module').then((m) => m.KtbDashboardViewModule),
      },
      { path: 'chart', component: KtbTestChartComponent },
      {
        path: 'create',
        loadChildren: () =>
          import('./_views/ktb-project-view/ktb-project-view-create.module').then((m) => m.KtbProjectViewCreateModule),
      },
      {
        path: 'project/:projectName',
        loadChildren: () =>
          import('./_views/ktb-project-view/ktb-project-view.module').then((m) => m.KtbProjectViewModule),
      },
      { path: 'trace/:keptnContext', canActivate: [TraceDeepLinkGuard], children: [] },
      { path: 'trace/:keptnContext/:eventSelector', canActivate: [TraceDeepLinkGuard], children: [] },
      {
        path: 'evaluation/:shkeptncontext',
        loadChildren: () =>
          import('./_views/ktb-evaluation-view/ktb-evaluation-view.module').then((m) => m.KtbEvaluationViewModule),
      },
      { path: '**', component: NotFoundComponent },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes, routingConfiguration)],
  exports: [RouterModule],
})
class AppRouting {}

export { AppRouting, routes };
