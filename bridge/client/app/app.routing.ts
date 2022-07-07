import { NgModule } from '@angular/core';
import { ExtraOptions, RouterModule, Routes } from '@angular/router';
import { DashboardLegacyComponent } from './dashboard-legacy/dashboard-legacy.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { EvaluationBoardComponent } from './evaluation-board/evaluation-board.component';
import { KtbSettingsViewComponent } from './_views/ktb-settings-view/ktb-settings-view.component';
import { KtbSecretsListComponent } from './_components/ktb-secrets-list/ktb-secrets-list.component';
import { KtbCreateSecretFormComponent } from './_components/ktb-create-secret-form/ktb-create-secret-form.component';
import { KtbProjectSettingsComponent } from './_components/ktb-project-settings/ktb-project-settings.component';
import { KtbCreateServiceComponent } from './_components/ktb-create-service/ktb-create-service.component';
import { KtbServiceSettingsOverviewComponent } from './_components/ktb-service-settings/ktb-service-settings-overview/ktb-service-settings-overview.component';
import { KtbServiceSettingsComponent } from './_components/ktb-service-settings/ktb-service-settings.component';
import { KtbEditServiceComponent } from './_components/ktb-edit-service/ktb-edit-service.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { PendingChangesGuard } from './_guards/pending-changes.guard';
import { KtbErrorViewComponent } from './_views/ktb-error-view/ktb-error-view.component';
import { AppComponent } from './app.component';

const routingConfiguration: ExtraOptions = {
  paramsInheritanceStrategy: 'always',
};

const lazyLoadEnvironmentView = (): Promise<unknown> =>
  import('./_views/ktb-environment-view/ktb-environment-view.module').then((m) => m.KtbEnvironmentViewModule);

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
      { path: 'dashboard', component: DashboardLegacyComponent },
      {
        path: 'create',
        component: ProjectBoardComponent,
        children: [{ path: 'project', component: KtbProjectSettingsComponent, canDeactivate: [PendingChangesGuard] }],
      },
      {
        path: 'project/:projectName',
        component: ProjectBoardComponent,
        children: [
          { path: '', pathMatch: 'full', loadChildren: lazyLoadEnvironmentView },
          { path: 'environment', pathMatch: 'full', loadChildren: lazyLoadEnvironmentView },
          { path: 'environment/stage/:stageName', loadChildren: lazyLoadEnvironmentView },
          {
            path: 'settings',
            component: KtbSettingsViewComponent,
            children: [
              { path: 'project', component: KtbProjectSettingsComponent, canDeactivate: [PendingChangesGuard] },
              {
                path: 'services',
                component: KtbServiceSettingsComponent,
                children: [
                  { path: 'create', component: KtbCreateServiceComponent },
                  { path: 'edit/:serviceName', component: KtbEditServiceComponent },
                  { path: '', pathMatch: 'full', component: KtbServiceSettingsOverviewComponent },
                ],
              },
              {
                path: 'uniform',
                children: [
                  {
                    path: 'integrations',
                    loadChildren: () =>
                      import('./_views/ktb-integration-view/ktb-integration-view.module').then(
                        (m) => m.KtbIntegrationViewModule
                      ),
                  },
                  {
                    path: 'secrets',
                    component: KtbSecretsListComponent,
                  },
                  { path: 'secrets/add', component: KtbCreateSecretFormComponent },
                  { path: '', pathMatch: 'full', redirectTo: 'integrations' },
                ],
              },
              {
                path: 'support',
                children: [
                  {
                    path: 'common-use-cases',
                    loadChildren: () =>
                      import('./_views/ktb-common-use-cases-view/ktb-common-use-cases-view.module').then(
                        (m) => m.KtbCommonUseCasesViewModule
                      ),
                  },
                  { path: '', pathMatch: 'full', redirectTo: 'common-use-cases' },
                ],
              },
              { path: '', pathMatch: 'full', redirectTo: 'project' },
            ],
          },
          {
            path: 'service',
            loadChildren: () =>
              import('./_views/ktb-service-view/ktb-service-view.module').then((m) => m.KtbServiceViewModule),
          },
          {
            path: 'sequence',
            loadChildren: () =>
              import('./_views/ktb-sequence-view/ktb-sequence-view.module').then((m) => m.KtbSequenceViewModule),
          },
        ],
      },
      { path: 'trace/:shkeptncontext', component: ProjectBoardComponent },
      { path: 'trace/:shkeptncontext/:eventselector', component: ProjectBoardComponent },
      { path: 'evaluation/:shkeptncontext', component: EvaluationBoardComponent },
      { path: 'evaluation/:shkeptncontext/:eventselector', component: EvaluationBoardComponent },
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
