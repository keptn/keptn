import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { KtbProjectViewComponent } from './ktb-project-view.component';
import { KtbSettingsViewComponent } from '../ktb-settings-view/ktb-settings-view.component';
import { KtbProjectSettingsComponent } from '../../_components/ktb-project-settings/ktb-project-settings.component';
import { PendingChangesGuard } from '../../_guards/pending-changes.guard';
import { KtbServiceSettingsComponent } from '../../_components/ktb-service-settings/ktb-service-settings.component';
import { KtbCreateServiceComponent } from '../../_components/ktb-create-service/ktb-create-service.component';
import { KtbEditServiceComponent } from '../../_components/ktb-edit-service/ktb-edit-service.component';
import { KtbServiceSettingsOverviewComponent } from '../../_components/ktb-service-settings/ktb-service-settings-overview/ktb-service-settings-overview.component';
import { KtbSecretsListComponent } from '../../_components/ktb-secrets-list/ktb-secrets-list.component';
import { KtbCreateSecretFormComponent } from '../../_components/ktb-create-secret-form/ktb-create-secret-form.component';

const lazyLoadEnvironmentView = (): Promise<unknown> =>
  import('../ktb-environment-view/ktb-environment-view.module').then((m) => m.KtbEnvironmentViewModule);

const routes: Routes = [
  {
    path: '',
    component: KtbProjectViewComponent,
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
                  import('../ktb-integration-view/ktb-integration-view.module').then((m) => m.KtbIntegrationViewModule),
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
                  import('../ktb-common-use-cases-view/ktb-common-use-cases-view.module').then(
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
        loadChildren: () => import('../ktb-service-view/ktb-service-view.module').then((m) => m.KtbServiceViewModule),
      },
      {
        path: 'sequence',
        loadChildren: () =>
          import('../ktb-sequence-view/ktb-sequence-view.module').then((m) => m.KtbSequenceViewModule),
      },
    ],
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class KtbProjectViewRoutingModule {}
