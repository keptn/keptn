import {NgModule} from '@angular/core';
import {Routes, RouterModule, ExtraOptions} from '@angular/router';
import {DashboardComponent} from './dashboard/dashboard.component';
import {ProjectBoardComponent} from './project-board/project-board.component';
import {EvaluationBoardComponent} from './evaluation-board/evaluation-board.component';
import {ForwarderGuard} from './_forwarder/forwarder_guard';
import {KtbUniformViewComponent} from './_views/ktb-uniform-view/ktb-uniform-view.component';
import {KtbIntegrationViewComponent} from './_views/ktb-integration-view/ktb-integration-view.component';
import {KtbSettingsViewComponent} from './_views/ktb-settings-view/ktb-settings-view.component';
import {KtbServiceViewComponent} from './_views/ktb-service-view/ktb-service-view.component';
import {KtbSequenceViewComponent} from './_views/ktb-sequence-view/ktb-sequence-view.component';
import {KtbEnvironmentViewComponent} from './_views/ktb-environment-view/ktb-environment-view.component';

const routingConfiguration: ExtraOptions = {
  paramsInheritanceStrategy: 'always'
};

const routes: Routes = [
  {path: '', pathMatch: 'full', redirectTo: 'dashboard'},
  {path: 'dashboard', component: DashboardComponent},
  {path: 'project/:projectName', component: ProjectBoardComponent, children: [
      {path: '', pathMatch: 'full', component: KtbEnvironmentViewComponent},
      {path: 'uniform', component: KtbUniformViewComponent},
      {path: 'integration', component: KtbIntegrationViewComponent},
      {path: 'settings', component: KtbSettingsViewComponent},
      {path: 'service', component: KtbServiceViewComponent},
      {path: 'sequence', component: KtbSequenceViewComponent},
      {path: 'service/:serviceName', component: KtbServiceViewComponent},
      {path: 'service/:serviceName/context/:shkeptncontext', component: KtbServiceViewComponent},
      {path: 'service/:serviceName/context/:shkeptncontext/stage/:stage', component: KtbServiceViewComponent},
      {path: 'sequence/:shkeptncontext', component: KtbSequenceViewComponent},
      {path: 'sequence/:shkeptncontext/event/:eventId', component: KtbSequenceViewComponent},
      {path: 'sequence/:shkeptncontext/stage/:stage', component: KtbSequenceViewComponent},
  ]},
  {path: 'trace/:shkeptncontext', component: ProjectBoardComponent},
  {path: 'trace/:shkeptncontext/:eventselector', component: ProjectBoardComponent},
  {path: 'evaluation/:shkeptncontext', component: EvaluationBoardComponent},
  {path: 'evaluation/:shkeptncontext/:eventselector', component: EvaluationBoardComponent},
  {path: 'project/:projectName/:serviceName', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: 'project/:projectName/:serviceName/:contextId/:eventId', component: ProjectBoardComponent, canActivate: [ForwarderGuard]}, // deprecated
  {path: '**', redirectTo: ''}
];

@NgModule({
  imports: [RouterModule.forRoot(routes, routingConfiguration)],
  exports: [RouterModule]
})
class AppRouting {
}

export {AppRouting, routes};
