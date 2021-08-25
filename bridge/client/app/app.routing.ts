import { NgModule } from '@angular/core';
import { ExtraOptions, RouterModule, Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/dashboard.component';
import { ProjectBoardComponent } from './project-board/project-board.component';
import { EvaluationBoardComponent } from './evaluation-board/evaluation-board.component';
import { ForwarderGuard } from './_forwarder/forwarder_guard';
import { KtbUniformViewComponent } from './_views/ktb-uniform-view/ktb-uniform-view.component';
import { KtbIntegrationViewComponent } from './_views/ktb-integration-view/ktb-integration-view.component';
import { KtbSettingsViewComponent } from './_views/ktb-settings-view/ktb-settings-view.component';
import { KtbServiceViewComponent } from './_views/ktb-service-view/ktb-service-view.component';
import { KtbSequenceViewComponent } from './_views/ktb-sequence-view/ktb-sequence-view.component';
import { KtbEnvironmentViewComponent } from './_views/ktb-environment-view/ktb-environment-view.component';
import { KtbKeptnServicesListComponent } from './_components/ktb-keptn-services-list/ktb-keptn-services-list.component';
import { KtbSecretsListComponent } from './_components/ktb-secrets-list/ktb-secrets-list.component';
import { KtbCreateSecretFormComponent } from './_components/ktb-create-secret-form/ktb-create-secret-form.component';
import { KtbProjectSettingsComponent } from './_components/ktb-project-settings/ktb-project-settings.component';

const routingConfiguration: ExtraOptions = {
  paramsInheritanceStrategy: 'always',
};

const routes: Routes = [
  {path: '', pathMatch: 'full', redirectTo: 'dashboard'},
  {path: 'dashboard', component: DashboardComponent},
  {
    path: 'create', component: ProjectBoardComponent, children: [
      {path: 'project', component: KtbProjectSettingsComponent, data: {isCreateMode: true}},
    ],
  },
  {
    path: 'project/:projectName', component: ProjectBoardComponent, children: [
      {path: '', pathMatch: 'full', component: KtbEnvironmentViewComponent},
      {
        path: 'uniform', component: KtbUniformViewComponent, children: [
          {path: 'services', component: KtbKeptnServicesListComponent},
          {path: 'secrets', component: KtbSecretsListComponent},
          {path: 'secrets/add', component: KtbCreateSecretFormComponent},
          {path: '', pathMatch: 'full', redirectTo: 'services'},
        ],
      },
      {path: 'integration', component: KtbIntegrationViewComponent},
      {
        path: 'settings', component: KtbSettingsViewComponent, children: [
          {path: 'project', component: KtbProjectSettingsComponent, data: {isCreateMode: false}},
          {path: '', pathMatch: 'full', redirectTo: 'project'},
        ],
      },
      {path: 'service', component: KtbServiceViewComponent},
      {path: 'sequence', component: KtbSequenceViewComponent},
      {path: 'service/:serviceName', component: KtbServiceViewComponent},
      {path: 'service/:serviceName/context/:shkeptncontext', component: KtbServiceViewComponent},
      {path: 'service/:serviceName/context/:shkeptncontext/stage/:stage', component: KtbServiceViewComponent},
      {path: 'sequence/:shkeptncontext', component: KtbSequenceViewComponent},
      {path: 'sequence/:shkeptncontext/event/:eventId', component: KtbSequenceViewComponent},
      {path: 'sequence/:shkeptncontext/stage/:stage', component: KtbSequenceViewComponent},
    ],
  },
  {path: 'trace/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext/:eventselector', component: ProjectBoardComponent},
  {path: 'evaluation/:shkeptncontext', component: EvaluationBoardComponent},
  {path: 'evaluation/:shkeptncontext/:eventselector', component: EvaluationBoardComponent},
  {path: 'project/:projectName/:serviceName', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId/:eventId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: '**', redirectTo: ''},
];

@NgModule({
  imports: [RouterModule.forRoot(routes, routingConfiguration)],
  exports: [RouterModule],
})
class AppRouting {
}

export { AppRouting, routes };
